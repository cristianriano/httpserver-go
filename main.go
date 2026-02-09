package main

import (
	"context"
	"fmt"
	"net"
	"os/signal"
	"syscall"
	"time"
)

const (
	defaultBufferSize  = 1024
	connectionDeadline = 30 * time.Second
)

func main() {
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGTERM, syscall.SIGINT)
	defer stop()

	listener, err := Start(ctx, ":8085")
	if err != nil {
		panic(err)
	}

	<-ctx.Done()
	listener.Close()
	fmt.Printf("Stopped: %v", ctx.Err())
}

func Start(ctx context.Context, addr string) (net.Listener, error) {
	listener, err := net.Listen("tcp", addr)

	if err != nil {
		return nil, err
	}

	fmt.Printf("Listening on %s\n\n", listener.Addr())
	go listen(listener)

	return listener, nil
}

func listen(listener net.Listener) {
	for {
		conn, err := listener.Accept()
		if err != nil {
			fmt.Printf("Error accepting conn: %v", err)
		}

		fmt.Printf("Accepted conn %s\n", conn.RemoteAddr())
		go processConn(conn)
	}
}

func processConn(conn net.Conn) {
	for {
		stream := make([]byte, defaultBufferSize)
		conn.SetDeadline(time.Now().Add(connectionDeadline))

		_, err := conn.Read(stream)

		if err != nil {
			fmt.Printf("Error reading %s\n", conn.RemoteAddr())
			break
		}

		for i, c := range stream {
			if c == '#' {
				fmt.Printf("[%s]: %s", conn.RemoteAddr(), stream[0:i])
				fmt.Printf("'#' sent by %s. Connection closed\n\n", conn.LocalAddr())
				conn.Close()
				return
			}
		}

		fmt.Printf("[%s]: %s", conn.RemoteAddr(), stream)
	}

	fmt.Printf("Connection close by server %s\n\n", conn.RemoteAddr())
	conn.Close()
}
