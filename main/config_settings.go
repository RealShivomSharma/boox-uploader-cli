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
	url := os.Getenv("BOOX_TABLET_URL")
	return config_settings{boox_url: url, boox_ip: ip}
}
