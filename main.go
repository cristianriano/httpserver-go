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
	time.Sleep(10 * time.Second)

	stream := make([]byte, 10)
	n, err := conn.Read(stream)

	if err != nil {
		fmt.Printf("Error reading %s\n", conn.RemoteAddr())
		conn.Close()
	}

	if n > 0 {
		fmt.Printf("[%s]: %s\n", conn.RemoteAddr(), stream)
	}
	fmt.Printf("Closing %s\n\n", conn.RemoteAddr())
	conn.Close()
}
