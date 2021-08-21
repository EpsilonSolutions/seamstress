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

func cmdinvoke(function string, args ...[]byte) {
	client := readConfig()
	r, err := client.Invoke(function, args...)
	if err != nil {
		log.Fatalln("error invoking smart contract:", err)
	}
	if _, err := os.Stdout.Write(r.Payload); err != nil {
		log.Fatalln("error writing to stdout:", err)
	}
}
