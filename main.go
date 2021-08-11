package main

import (
	"log"
	"os"
	"time"

	"github.com/EpsilonSolutions/seamstress/fabric"
	"github.com/nats-io/nats.go"
)

func requiredEnv(key string) string {
	v, ok := os.LookupEnv(key)
	if !ok {
		log.Fatalln(key, "must be defined")
	}
	return v
}

func defaultEnv(key, val string) string {
	v, ok := os.LookupEnv(key)
	if !ok {
		return val
	}
	return v
}

func main() {
	log.SetFlags(0)
	log.SetPrefix("seamstress: ")
	log.Println("Configuring worker...")

	var (
		natsURL           = requiredEnv("NATS_URL")
		fabricProfilePath = requiredEnv("FABRIC_CONNECTION_PROFILE_PATH")
		fabricChannel     = requiredEnv("FABRIC_CHANNEL")
		fabricOrg         = requiredEnv("FABRIC_ORG")
		fabricContract    = requiredEnv("FABRIC_CONTRACT")

		// These values should start with ./local/
		fabricIDPath  = requiredEnv("FABRIC_ID_PATH")
		fabricPEMPath = requiredEnv("FABRIC_PEM_PATH")
		fabricKeyPath = defaultEnv("FABRIC_KEY_PATH", "./local/keystore/privateKey")

		// These values are base64 encoded config files
		fabricConnectionProfileBase64 = requiredEnv("FABRIC_CONNECTION_PROFILE_B64")
		fabricClientKeyBase64         = requiredEnv("FABRIC_CLIENT_KEY_B64")
		fabricClientIDBase64          = requiredEnv("FABRIC_CLIENT_ID_B64")
		fabricClientPEMBase64         = requiredEnv("FABRIC_CLIENT_PEM_B64")
	)

	requireFile(fabricConnectionProfileBase64, fabricProfilePath)
	requireFile(fabricClientKeyBase64, fabricKeyPath)
	requireFile(fabricClientIDBase64, fabricIDPath)
	requireFile(fabricClientPEMBase64, fabricPEMPath)

	log.Println("Config confirmed: TODO show configuration")

	// TODO: configurable max workers
	maxWorkers := 5
	workch := make(chan *nats.Msg)
	for i := 1; i <= maxWorkers; i++ {
		w := &worker{
			fab: &fabric.Client{
				Profile:  fabricProfilePath,
				Channel:  fabricChannel,
				User:     fabricOrg,
				Org:      fabricOrg,
				Contract: fabricContract,
			},
			work: workch,
		}
		go w.run()
	}

	log.Println("Workers launched:", maxWorkers)

	l := listener{
		prefix:   fabricChannel + "." + fabricContract + ".",
		worker:   workch,
		shutdown: make(chan struct{}),
	}

	// Intercept system kill signal for graceful shutdown
	go cleanup(l.shutdown)

	log.Println("Listening for NATS messages...")

	l.listen(natsURL)

	log.Println("Exiting")
}

type worker struct {
	fab  *fabric.Client
	work chan *nats.Msg
}

func (w *worker) run() {
	for m := range w.work {
		w.invoke(m)
	}
}

func (w *worker) invoke(m *nats.Msg) {
	topic := w.fab.Channel + "." + w.fab.Contract + "."
	ret, err := w.fab.Invoke(m.Subject[len(topic):], m.Data)
	if err != nil {
		log.Println("error invoking smart contract:", err)
	}
	if err := m.Respond(ret); err == nats.ErrMsgNoReply {
		// If reply is not set we do nothing
	} else if err != nil {
		log.Println("error replying to NATS message:", err)
	}
}

type listener struct {
	prefix   string
	worker   chan *nats.Msg
	shutdown chan struct{}
}

func (l *listener) listen(natsURL string) {
	nc, err := nats.Connect(natsURL)
	if err != nil {
		// TODO: decide how to poll for NATS
		// (in particular what if NATS dies at some point /after/ startup)
		log.Fatalln("error connecting to NATS:", err)
	}

	sub, err := nc.ChanSubscribe(l.prefix+"*", l.worker)
	if err != nil {
		log.Fatalln("error subscribing to NATS topic:", err)
	}

	<-l.shutdown
	// TODO: test this to make sure that drain actually allows us to close the channel
	if err := sub.Drain(); err != nil {
		log.Fatalln("error draining worker on shutdown:", err)
	}
	close(l.worker)
}

func cleanup(ch chan struct{}) {
	// TODO: intercept kill signal from system (i.e. kubernetes)
	time.Sleep(20 * time.Second)
	close(ch)
}
