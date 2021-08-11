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
	)

	// TODO: make these paths configurable
	requireFile("FABRIC_CONNECTION_PROFILE_B64", fabricProfilePath)
	requireFile("FABRIC_CLIENT_KEY_B64", fabricKeyPath)
	requireFile("FABRIC_CLIENT_ID_B64", fabricIDPath)
	requireFile("FABRIC_CLIENT_PEM_B64", fabricPEMPath)

	// TODO: consider spinning up a fixed number of greenthread workers
	w := &worker{
		fab: &fabric.Client{
			Profile:  fabricProfilePath,
			Channel:  fabricChannel,
			User:     fabricOrg,
			Org:      fabricOrg,
			Contract: fabricContract,
		},
		shutdown: make(chan struct{}),
	}

	log.Println("Worker configured: TODO show configuration")

	go cleanup(w.shutdown)
	w.listen(natsURL)

	log.Println("Exiting")
}

type worker struct {
	fab      *fabric.Client
	shutdown chan struct{}
}

func (w *worker) listen(natsURL string) {
	nc, err := nats.Connect(natsURL)
	if err != nil {
		// TODO: decide how to poll for NATS
		// (in particular what if NATS dies at some point /after/ startup)
		log.Fatalln("error connecting to NATS:", err)
	}

	topic := w.fab.Channel + "." + w.fab.Contract + "."
	sub, err := nc.Subscribe(topic+"*", func(m *nats.Msg) {
		ret, err := w.fab.Invoke(m.Subject[len(topic):], m.Data)
		if err != nil {
			log.Println("error invoking smart contract:", err)
		}
		if err := m.Respond(ret); err == nats.ErrMsgNoReply {
			// If reply is not set we do nothing
		} else if err != nil {
			log.Println("error replying to NATS message:", err)
		}
	})
	if err != nil {
		log.Fatalln("error subscribing to NATS topic:", err)
	}

	<-w.shutdown
	if err := sub.Drain(); err != nil {
		log.Println("error draining worker on shutdown:", err)
	}
}

// TODO: listen for shutdown instruction
func cleanup(ch chan struct{}) {
	// TODO: intercept kill signal from system (i.e. kubernetes)
	time.Sleep(20 * time.Second)
	close(ch)
}
