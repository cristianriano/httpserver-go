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

		conn.SetDeadline(time.Now().Add(100 * time.Millisecond))
		n, err := conn.Write([]byte("Hello"))
		require.NoError(t, err)
		assert.Greater(t, n, 0)

		stream := make([]byte, 10)
		n, err = conn.Read(stream)
		assert.Equal(t, "Hello", string(stream[0:n]))
		conn.Close()
	})

	t.Run("Start returns a listener that close connections with '#'", func(t *testing.T) {
		// Use a random port assigned by the OS
		listener, err := Start(context.TODO(), addr)
		require.NoError(t, err)
		require.NotNil(t, listener)
		defer listener.Close()

		conn, err := net.Dial("tcp", listener.Addr().String())
		require.NoError(t, err)

		conn.SetDeadline(time.Now().Add(100 * time.Millisecond))
		n, err := conn.Write([]byte("Hello#Invalid"))

		require.NoError(t, err)
		assert.Greater(t, n, 0)

		stream := make([]byte, 20)
		n, err = conn.Read(stream)

		require.NoError(t, err)
		assert.Equal(t, "Hello", string(stream[0:n]))
	})
}
