package main

import (
	"flag"
	"os"
)

func init() {
	addHelp(
		"",
		`
seamstress is a tool for manually connecting to and transacting with smart
contracts running on IBM Blockchain Platform.

Usage:

	seamstress command [arguments]

The commands are:

	connect        connects to the peer / channel / smart contract specified and persists connection information
	invoke         invokes the smart contract subaction specified with the data provided

Use "seamstress help [command]" for more information about a command.
		`,
	)
}

func main() {
	cmd := ""
	if len(os.Args) > 1 {
		cmd = os.Args[1]
	}
	switch cmd {
	case "connect":
		cmdconnect()
	case "invoke":
		cmdinvoke()
	case "help":
		fs := flag.NewFlagSet("help", 0)
		fs.Parse(os.Args[2:])
		cmdhelp(fs.Arg(0))
	default:
		cmdhelp("")
	}
}
