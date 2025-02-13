package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"net"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main() {
	timeout := flag.Duration("timeout", 10*time.Second, "timeout duration")
	flag.Parse()

	args := flag.Args()
	if len(args) < 2 {
		log.Fatalln("Both host and port are required")
	}
	host, port := args[0], args[1]
	if host == "" || port == "" {
		log.Fatalln("host and port cannot be empty")
	}

	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	address := net.JoinHostPort(host, port)

	client := NewTelnetClient(address, *timeout, os.Stdin, os.Stdout)

	if err := client.Connect(); err != nil {
		log.Printf("Failed to connect to %s: %v", address, err)
		stop()
	}
	defer client.Close()

	log.Printf("Connected to %s!\n", address)

	startRoutine(stop, "Send", client.Send)
	startRoutine(stop, "Receive", client.Receive)

	<-ctx.Done()
	log.Println("Connection closed.")
}

func startRoutine(stop func(), name string, task func() error) {
	go func() {
		defer func() {
			log.Printf("%s routine stopped.", name)
			stop()
		}()
		if err := task(); err != nil {
			fmt.Fprintf(os.Stderr, "%s error: %v\n", name, err)
		}
	}()
}
