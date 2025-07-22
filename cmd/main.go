package main

import (
	"log"

	"github.com/Babushkin05/subscription-organizer/internal/config"
)

func main() {
	// Config
	cfg := config.MustLoad()
	log.Println("Config was loaded correctful!")

}
