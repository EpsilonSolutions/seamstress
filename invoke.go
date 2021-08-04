package main

import "log"

func init() {
	addHelp(
		"invoke",
		`
TODO: add help for invoke command
		`,
	)
}

func cmdinvoke(function string, args [][]byte) {
	client := readConfig()
	if err := client.Invoke(function, args); err != nil {
		log.Fatalln("error invoking smart contract:", err)
	}
}
