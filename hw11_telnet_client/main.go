package main

import (
	"context"
	"errors"
	"flag"
	"fmt"
	"io"
	"net"
	"os"
	"os/signal"
	"time"
)

func main() {
	var duration time.Duration
	flag.DurationVar(&duration, "timeout", time.Second*10, "the duration to wait before exit")
	flag.Parse()

	args := flag.Args()
	if len(args) < 2 {
		fmt.Fprintf(os.Stderr, "Need at least two arguments: host and port")
		os.Exit(1)
	}

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt)
	defer stop()

	client := NewTelnetClient(net.JoinHostPort(args[0], args[1]), duration, os.Stdin, os.Stdout)
	defer client.Close()

	if err := client.Connect(); err != nil {
		fmt.Fprintf(os.Stderr, "Startup failed: %v\n", err)
		return
	}
	fmt.Fprintf(os.Stderr, "Connected to %v\n", net.JoinHostPort(args[0], args[1]))

	errChan := make(chan error, 2)

	go func() {
		for {
			if err := client.Send(); err != nil {
				errChan <- err
				break
			}
		}
	}()

	go func() {
		for {
			if err := client.Receive(); err != nil {
				errChan <- err
				break
			}
		}
	}()

	select {
	case err := <-errChan:
		if err != nil && !errors.Is(err, io.EOF) {
			fmt.Fprintf(os.Stderr, "Communication error: %v\n", err)
		} else {
			fmt.Fprintf(os.Stderr, "Input closed\n")
		}
	case <-ctx.Done():
		fmt.Fprintf(os.Stderr, "Interrupt received\n")
	}
}
