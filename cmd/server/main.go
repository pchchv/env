package main

import (
	"log"
	"os"

	"github.com/labstack/echo"
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
	e := echo.New()
	log.Fatal(e.Start(":" + getEnvValue("PORT")).Error())
}
