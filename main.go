package main

import (
	"log"

	"github.com/Ewan-Greer09/finance-app/api"
)

func main() {
	err := api.NewAPI().Run()
	if err != nil {
		log.Panic(err)
	}
}
