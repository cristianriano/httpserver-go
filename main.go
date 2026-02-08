package main

import (
	"context"
	"fmt"
	"net"
	"os/signal"
	"syscall"
	"time"
)

func main() {
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGTERM, syscall.SIGINT)
	defer stop()

	listener, err := net.Listen("tcp", ":8085")
	if err != nil {
		panic(err)
	}

	fmt.Printf("Listening on %s\n", listener.Addr())

	ticker := time.NewTicker(10 * time.Second)
	for {
		select {
		case <-ctx.Done():
			listener.Close()
			fmt.Printf("Stopped: %v", ctx.Err())
		case <-ticker.C:
			fmt.Println("Listening...")
		}
	}
}
