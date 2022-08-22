package main

import (
	"context"
	"flag"
	"log"
	"net"
	"os"
	"sync"
	"time"
)

func main() {
	var timeout time.Duration
	flag.DurationVar(&timeout, "timeout", 10*time.Second, "timeout to connect")
	flag.Parse()

	telnet := NewTelnetClient(
		net.JoinHostPort(flag.Arg(0), flag.Arg(1)),
		timeout,
		os.Stdin,
		os.Stdout,
	)

	if err := telnet.Connect(); err != nil {
		log.Printf("connection error: %s", err)
		return
	}
	defer telnet.Close()

	var wg sync.WaitGroup
	ctx, cancelFn := context.WithCancel(context.Background())
	defer cancelFn()

	// receiver
	go func() {
		defer wg.Done()

	O:
		for {
			select {
			case <-ctx.Done():
				break O
			default:
				if err := telnet.Receive(); err != nil {
					log.Printf("telnet receive error: %s", err)
					cancelFn()
				}
			}
		}
	}()
	wg.Add(1)

	// sender
	go func() {
		defer wg.Done()

	O:
		for {
			select {
			case <-ctx.Done():
				break O
			default:
				if err := telnet.Send(); err != nil {
					log.Printf("telnet send error: %s", err)
					cancelFn()
				}
			}
		}
	}()
	wg.Add(1)

	wg.Wait()
}
