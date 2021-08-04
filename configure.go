package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"

	"github.com/EpsilonSolutions/seamstress/fabric"
)

const localConfigFile = "~/.seamstress"

func init() {
	addHelp(
		"configure",
		`
TODO: add help for connect command
		`,
	)

	createFileIfNotExists(localConfigFile)
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
		if _, err := os.Create(filename); err != nil {
			log.Fatalln("unknown filesystem error creating local cache:", err)
		}
	} else if err != nil {
		log.Fatalln("unknown filesystem error:", err)
	}
}

func readConfig() *fabric.Client {
	f, err := os.Open(localConfigFile)
	if err != nil {
		log.Fatalln("error opening local config file:", err)
	}
	client := fabric.Client{}
	if err := json.NewDecoder(f).Decode(&client); err != nil {
		log.Fatalln("error writing config to local file:", err)
	}
	return &client
}

func writeConfig(client *fabric.Client) {
	f, err := os.Open(localConfigFile)
	if err != nil {
		log.Fatalln("error opening local config file:", err)
	}
	if err := json.NewEncoder(f).Encode(client); err != nil {
		log.Fatalln("error writing config to local file:", err)
	}
}
