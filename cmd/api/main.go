package main

import (
	"log"

	server "github.com/alexPavlikov/time-tracker/cmd"
)

func main() {
	if err := server.Run(); err != nil {
		log.Fatal(err)
	}
}
