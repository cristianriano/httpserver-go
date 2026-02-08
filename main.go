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
	defaultBufferSize = 1024
)

func main() {
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGTERM, syscall.SIGINT)
	defer stop()

	listener, err := net.Listen("tcp", ":8085")
	if err != nil {
		panic(err)
	}

	fmt.Printf("Listening on %s\n\n", listener.Addr())
	go listen(listener)

	<-ctx.Done()
	listener.Close()
	fmt.Printf("Stopped: %v", ctx.Err())
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
		conn.SetDeadline(time.Now().Add(5 * time.Second))

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
