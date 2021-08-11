package main

import (
	"encoding/base64"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
)

func envToFile(b64, path string) error {
	b, err := base64.StdEncoding.DecodeString(b64)
	if err != nil {
		return err
	}
	if err := os.MkdirAll(filepath.Dir(path), 0777); err != nil {
		return err
	}
	return ioutil.WriteFile(path, b, 0666)
}

func requireFile(b64, path string) {
	if err := envToFile(b64, path); err != nil {
		log.Fatalln("could not write to", path, ". error:", err)
	}
}
