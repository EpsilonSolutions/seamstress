package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/user"

	"github.com/EpsilonSolutions/seamstress/fabric"
)

func init() {
	addHelp(
		"configure",
		`
TODO: add help for connect command
		`,
	)

}

func cmdconfigure() {
	// TODO: consider adding validation step
	writeConfig(&fabric.Client{
		Profile: prompt("enter filepath to profile:"),
		Channel: prompt("enter channel ID:"),
		//		User:     prompt("enter user:"),
		Org:      prompt("enter organization name:"),
		Contract: prompt("enter smart contract ID / chaincode ID:"),
	})
}

func prompt(msg string) string {
	log.Printf(msg)
	value := ""
	_, err := fmt.Scan(&value)
	if err != nil {
		log.Fatalln("failed to get value from prompt:", err)
	}
	return value
}

func createFileIfNotExists(filename string) {
	if _, err := os.Stat(filename); os.IsNotExist(err) {
		f, err := os.Create(filename)
		if err != nil {
			log.Fatalln("unknown filesystem error creating local cache:", err)
		}
		f.Close()
	} else if err != nil {
		log.Fatalln("unknown filesystem error:", err)
	}
}

func readConfig() *fabric.Client {
	u, err := user.Current()
	if err != nil {
		log.Fatalln("error getting user:", err)
	}

	contents, err := ioutil.ReadFile(u.HomeDir + "/.seamstress")
	if err != nil {
		log.Fatalln("error reading seamstress file:", err)
	}

	client := fabric.Client{}
	if err := json.Unmarshal(contents, &client); err != nil {
		log.Fatalln("error writing config to local file:", err)
	}

	return &client
}

func writeConfig(client *fabric.Client) {
	u, err := user.Current()
	if err != nil {
		log.Fatalln("error getting user:", err)
	}

	contents, err := json.Marshal(client)
	if err != nil {
		log.Fatalln("error encoding client:", err)
	}

	if err := ioutil.WriteFile(u.HomeDir+"/.seamstress", contents, 0644); err != nil {
		log.Fatalln("error writing to seamstress file:", err)
	}
}
