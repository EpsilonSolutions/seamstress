package main

import (
	"encoding/base64"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
)

func envToFile(env, path string) error {
	v := os.Getenv(env)
	b, err := base64.StdEncoding.DecodeString(v)
	if err != nil {
		return err
	}
	if err := os.MkdirAll(filepath.Dir(path), 0777); err != nil {
		return err
	}
	return ioutil.WriteFile(path, b, 0666)
}

func requireFile(env, path string) {
	if err := envToFile(env, path); err != nil {
		log.Fatalln("could not write from", env, "to", path, ". error:", err)
	}
}
