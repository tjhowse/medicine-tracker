package main

import (
	"log"
	"os"

	"github.com/kelseyhightower/envconfig"
	"github.com/pelletier/go-toml/v2"
)

func main() {

	s := &Server{}
	// Check if settings.toml is present, if so: load it
	var settings ServerSettings
	// Try to load the settings from the environment
	err := envconfig.Process("FBF", &settings)
	if err != nil {
		log.Println("Failed to load settings from environment, using defaults")
	}
	// Load the file
	if f, err := os.Open("settings.toml"); err == nil {
		// Decode the file
		if err := toml.NewDecoder(f).Decode(&settings); err != nil {
			log.Println("Failed to decode settings.toml.")
		}
	}

	log.Printf("Starting with settings: %+v", settings)
	s.Init(settings)
	defer s.Close()

	s.Run()
}
