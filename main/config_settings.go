package main

import (
	"os"

	"github.com/joho/godotenv"
)

type config_settings struct {
	boox_url string
	boox_ip  string
}

func create_config() config_settings {
	godotenv.Load(".env")
	ip := os.Getenv("BOOX_TABLET_IP")
	url := os.Getenv("BOOX_TABLET_ADDRESS")

	var config_settings config_settings

	config_settings.boox_ip = ip
	config_settings.boox_url = url

	return config_settings
}
