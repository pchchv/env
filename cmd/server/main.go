package main

import (
	"log"
	"net/http"
	"os"

	"github.com/pchchv/env"
)

func init() {
	// Load values from .env into the system
	if err := env.Load(); err != nil {
		log.Panic("No .env file found")
	}
}

func getEnvValue(v string) string {
	value, exist := os.LookupEnv(v)
	if !exist {
		log.Fatalf("Value %v does not exist", v)
	}
	return value
}

func main() {
	err := http.ListenAndServe(":"+getEnvValue("PORT"), nil)
	if err != nil {
		log.Panic(err)
	}
}
