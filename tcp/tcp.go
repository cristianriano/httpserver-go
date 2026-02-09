package tcp

import (
	"context"
	"errors"
	"log/slog"
	"net"
	"time"
)

const (
	defaultBufferSize  = 1024
	connectionDeadline = 30 * time.Second
)

// Start creates a TCP listener in the provided address which echos back what it receives
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
			// Check if the error is because we intentionally closed the listener
			if errors.Is(err, net.ErrClosed) {
				slog.Debug("Listener closed, shutting down accept loop")
				return
			}

			slog.Error("Error accepting connections", "err", err)
			continue
		}

		slog.Debug("Accepted conn", "addr", conn.RemoteAddr())
		go processConn(conn)
	}
}

func processConn(conn net.Conn) {
	stream := make([]byte, defaultBufferSize)

	for {
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

				slog.Debug("Message received", "addr", conn.RemoteAddr(), "msg", stream[0:i])
				slog.Debug("Connection closed", "addr", conn.RemoteAddr())

				conn.Close()
				return
			}
		}

		conn.Write(stream[0:n])
		slog.Debug("Received", "addr", conn.RemoteAddr(), "msg", stream[0:n])
	}
}
