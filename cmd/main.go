package main

import (
	"help/cmd/api"
	"help/config"
	"log"
)

func main() {
	api := api.NewApi(config.Env.Port)
	if err := api.Run(); err != nil {
		log.Fatal(err)
	}
}
