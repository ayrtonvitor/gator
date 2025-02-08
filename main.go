package main

import (
	"fmt"
	"log"

	"github.com/ayrtonvitor/gator/internal/config"
)

func main() {
	conf, err := config.Read()
	if err != nil {
		log.Fatalf("[Error] Could not read %v", err)
	}
	fmt.Printf("Read config: %+v\n", conf)
}
