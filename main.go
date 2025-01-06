package main

import (
	"fmt"
	"log"

	"github.com/juaniten/gator/internal/config"
)

func main() {
	configuration, err := config.Read()
	if err != nil {
		log.Fatalf("error reading gator configuration: %v", err)
	}
	fmt.Printf("Read configuration file: %+v\n", configuration)

	err = configuration.SetUser("juan")
	if err != nil {
		log.Fatalf("error setting user name: %v", err)
	}

	configuration, err = config.Read()
	if err != nil {
		log.Fatalf("error reading gator configuration: %v", err)
	}
	fmt.Printf("Read configuration file again: %+v\n", configuration)

}
