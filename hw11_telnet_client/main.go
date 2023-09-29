package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"os"
	"os/signal"
	"time"
)

var (
	port    int
	host    string
	timeout time.Duration
)

func init() {
	flag.IntVar(&port, "port", 3000, "port")
	flag.StringVar(&host, "host", "localhost", "dns or ip address")
	flag.DurationVar(&timeout, "timeout", 10*time.Second, "timeout")
}

func main() {
	ctx, cancel := context.WithCancel(context.Background())

	defer cancel()

	go func() {
		s := make(chan os.Signal, 1)
		signal.Notify(s, os.Interrupt)

		<-s
		cancel()
	}()

	flag.Parse()
	address := fmt.Sprintf("%s:%d", host, port)
	client := NewTelnetClient(address, timeout, os.Stdin, os.Stdout)

	err := client.Connect()
	if err != nil {
		log.Printf("fail to connect: %v", err)
		return
	}

	defer func() {
		err := client.Close()
		if err != nil {
			log.Printf("fail to close: %v", err)
		}
	}()

	go func() {
		err := client.Receive()
		if err != nil {
			log.Printf("fail to receive: %v", err)
		}

		cancel()
	}()

	go func() {
		err := client.Send()
		if err != nil {
			log.Printf("fail to send: %v", err)
		}

		cancel()
	}()

	<-ctx.Done()
}
