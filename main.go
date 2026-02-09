package main

import (
	"context"
	"flag"
	"log/slog"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/cristianriano/httpserver-go/tcp"
)

const (
	defaultBufferSize  = 1024
	connectionDeadline = 30 * time.Second
)

func main() {
	debugMode := flag.Bool("debug", false, "enable debug logging")
	flag.Parse()

	logOpts := &slog.HandlerOptions{Level: slog.LevelInfo}
	if *debugMode {
		logOpts.Level = slog.LevelDebug
	}

	logger := slog.New(slog.NewTextHandler(os.Stdout, logOpts))
	slog.SetDefault(logger)

	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGTERM, syscall.SIGINT)
	defer stop()

	listener, err := tcp.Start(ctx, ":8085")
	if err != nil {
		panic(err)
	}

	<-ctx.Done()
	listener.Close()
	slog.Debug("Stopped", "err", ctx.Err())
}
