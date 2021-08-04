package main

import (
	"log"
	"os"
)

func init() {
	addHelp(
		"invoke",
		`
TODO: add help for invoke command
		`,
	)
}

func cmdinvoke(function string, args []byte) {
	client := readConfig()
	b, err := client.Invoke(function, args)
	if err != nil {
		log.Fatalln("error invoking smart contract:", err)
	}
	if _, err := os.Stdout.Write(b); err != nil {
		log.Fatalln("error writing to stdout:", err)
	}
}
