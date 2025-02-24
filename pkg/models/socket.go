package models

import (
	"net"
	"sync"
)

// SocketWrapper wraps a net.Conn with additional metadata
type SocketWrapper struct {
	ID     int32
	Socket net.Conn
	mu     sync.RWMutex

	// Connection state
	Connected  bool
	LastActive int64
	BytesSent  int64
	BytesRecv  int64

	// Buffer management
	ReadBuffer  []byte
	WriteBuffer []byte
}

// NewSocketWrapper creates a new socket wrapper
func NewSocketWrapper(id int32, socket net.Conn) *SocketWrapper {
	return &SocketWrapper{
		ID:          id,
		Socket:      socket,
		Connected:   true,
		ReadBuffer:  make([]byte, 8192),
		WriteBuffer: make([]byte, 8192),
	}
}

// Close closes the socket and cleans up resources
func (sw *SocketWrapper) Close() error {
	sw.mu.Lock()
	defer sw.mu.Unlock()

	if !sw.Connected {
		return nil
	}

	sw.Connected = false
	if sw.Socket != nil {
		err := sw.Socket.Close()
		sw.Socket = nil
		return err
	}
	return nil
}

// Write writes data to the socket
func (sw *SocketWrapper) Write(data []byte) (int, error) {
	sw.mu.Lock()
	defer sw.mu.Unlock()

	if !sw.Connected {
		return 0, net.ErrClosed
	}

	n, err := sw.Socket.Write(data)
	if err == nil {
		sw.BytesSent += int64(n)
	}
	return n, err
}

// Read reads data from the socket
func (sw *SocketWrapper) Read(data []byte) (int, error) {
	sw.mu.Lock()
	defer sw.mu.Unlock()

	if !sw.Connected {
		return 0, net.ErrClosed
	}

	n, err := sw.Socket.Read(data)
	if err == nil {
		sw.BytesRecv += int64(n)
	}
	return n, err
}

// IsConnected returns whether the socket is connected
func (sw *SocketWrapper) IsConnected() bool {
	sw.mu.RLock()
	defer sw.mu.RUnlock()
	return sw.Connected
}

// GetStats returns socket statistics
func (sw *SocketWrapper) GetStats() (bytesSent, bytesRecv int64) {
	sw.mu.RLock()
	defer sw.mu.RUnlock()
	return sw.BytesSent, sw.BytesRecv
}
