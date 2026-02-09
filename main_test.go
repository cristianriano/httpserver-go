package main

import (
	"context"
	"fmt"
	"net"
	"sync"
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
		defer conn.Close()

		conn.SetDeadline(time.Now().Add(100 * time.Millisecond))
		n, err := conn.Write([]byte("Hello"))
		require.NoError(t, err)
		assert.Greater(t, n, 0)

		stream := make([]byte, 10)
		n, err = conn.Read(stream)
		assert.Equal(t, "Hello", string(stream[:n]))
	})

	t.Run("Start returns a listener that close connections with '#'", func(t *testing.T) {
		// Use a random port assigned by the OS
		listener, err := Start(context.TODO(), addr)
		require.NoError(t, err)
		require.NotNil(t, listener)
		defer listener.Close()

		conn, err := net.Dial("tcp", listener.Addr().String())
		require.NoError(t, err)
		defer conn.Close()

		conn.SetDeadline(time.Now().Add(100 * time.Millisecond))
		n, err := conn.Write([]byte("Hello#Invalid"))

		require.NoError(t, err)
		assert.Greater(t, n, 0)

		stream := make([]byte, 20)
		n, err = conn.Read(stream)

		require.NoError(t, err)
		assert.Equal(t, "Hello", string(stream[:n]))
	})

	t.Run("Start when message sent is bigger than buffer size", func(t *testing.T) {
		// Use a random port assigned by the OS
		listener, err := Start(context.TODO(), addr)
		require.NoError(t, err)
		defer listener.Close()

		conn, err := net.Dial("tcp", listener.Addr().String())
		require.NoError(t, err)
		defer conn.Close()

		conn.SetDeadline(time.Now().Add(100 * time.Millisecond))

		packet := make([]byte, defaultBufferSize+1)
		for i := range packet {
			packet[i] = '1'
		}

		n, err := conn.Write(packet)
		require.NoError(t, err)
		assert.Greater(t, n, 0)

		stream := make([]byte, defaultBufferSize*2)
		n, err = conn.Read(stream)
		assert.Equal(t, n, defaultBufferSize+1)
	})

	t.Run("It keeps memory isolated for concurrent clients", func(t *testing.T) {
		// Use a random port assigned by the OS
		listener, err := Start(context.TODO(), addr)
		require.NoError(t, err)
		defer listener.Close()

		var wg sync.WaitGroup

		for i := range 50 {
			wg.Add(1)

			go func(j int) {
				defer wg.Done()

				conn, err := net.Dial("tcp", listener.Addr().String())
				assert.NoError(t, err)

				conn.SetDeadline(time.Now().Add(100 * time.Millisecond))

				msg := fmt.Sprintf("Hello %d", j)
				_, err = conn.Write([]byte(msg))
				assert.NoError(t, err)

				stream := make([]byte, 10)
				n, err := conn.Read(stream)
				assert.Equal(t, msg, string(stream[:n]))
			}(i)
		}

		wg.Wait()
	})
}
