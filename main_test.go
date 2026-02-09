package main

import (
	"context"
	"net"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestStart(t *testing.T) {
	addr := ":0"

	t.Run("Start listens for tcp connections in the given port", func(t *testing.T) {
		// Use a random port assigned by the OS
		listener, err := Start(context.TODO(), addr)
		require.NoError(t, err)
		defer listener.Close()

		conn, err := net.Dial("tcp", listener.Addr().String())
		require.NoError(t, err)

		conn.SetDeadline(time.Now().Add(10 * time.Millisecond))
		n, err := conn.Write([]byte("Hello"))
		require.NoError(t, err)
		assert.Greater(t, n, 0)

		stream := make([]byte, 5)
		conn.Read(stream)
		assert.Equal(t, "Hello", stream)
	})
}
