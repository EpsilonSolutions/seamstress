package main

import (
	"io"
	"log"
	"os"
	"strings"
)

var help = map[string]string{}

func addHelp(cmd, msg string) {
	msg = strings.TrimSpace(msg) + "\n"
	help[cmd] = msg
}

func cmdhelp(cmd string) {
	msg, ok := help[cmd]
	if !ok {
		log.Fatalf("Unknown help topic %q.  Run 'seamstress help'.", cmd)
	}
	io.WriteString(os.Stderr, msg)
}
