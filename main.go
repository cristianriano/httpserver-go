package main

import (
	"context"
	"log/slog"
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
	slog.Debug("Stopped", "err", ctx.Err())
}

func Start(ctx context.Context, addr string) (net.Listener, error) {
	listener, err := net.Listen("tcp", addr)

	if err != nil {
		return nil, err
	}

	slog.Debug("Started", "addr", listener.Addr())
	go listen(listener)

	return listener, nil
}

func listen(listener net.Listener) {
	for {
		conn, err := listener.Accept()
		if err != nil {
			slog.Error("Error accepting connections", "err", err)
		}

		slog.Debug("Accepted conn", "addr", conn.RemoteAddr())
		go processConn(conn)
	}
}

func processConn(conn net.Conn) {
	for {
		stream := make([]byte, defaultBufferSize)
		conn.SetDeadline(time.Now().Add(connectionDeadline))

		n, err := conn.Read(stream)

		if err != nil {
			slog.Error("Error reading", "err", err, "addr", conn.RemoteAddr())
			conn.Close()
			return
		}

		for i, c := range stream {
			if c == '#' {
				conn.Write(stream[0:i])

				slog.Info("Message received", "addr", conn.RemoteAddr(), "msg", stream[0:i])
				slog.Debug("Connection closed", "addr", conn.RemoteAddr())

				conn.Close()
				return
			}
		}

		conn.Write(stream[0:n])
		slog.Info("Received", "addr", conn.RemoteAddr(), "msg", stream[0:n])
	}
}
