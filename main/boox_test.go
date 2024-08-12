package main

import (
	"fmt"
	"log"
	"testing"

	"github.com/go-ping/ping"
)

func TestBooxConnnection(*testing.T) {
	var config = create_config()
	boox_ip := config.boox_ip
	pinger, err := ping.NewPinger(boox_ip)
	if err != nil {
		log.Fatal(err)
	}
	pinger.Count = 3
	pinger.Run()
	if err != nil {
		panic(err)
	}
	stats := pinger.Statistics()

	fmt.Println(stats)

}
