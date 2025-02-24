package util

import (
	"encoding/binary"
	"fmt"
	"io"
	"net"
)

// ReadUTF reads a UTF string from a reader
func ReadUTF(r io.Reader) (string, error) {
	var length uint16
	if err := binary.Read(r, binary.BigEndian, &length); err != nil {
		return "", fmt.Errorf("failed to read string length: %v", err)
	}

	buf := make([]byte, length)
	if _, err := io.ReadFull(r, buf); err != nil {
		return "", fmt.Errorf("failed to read string data: %v", err)
	}

	return string(buf), nil
}

// WriteUTF writes a UTF string to a writer
func WriteUTF(w io.Writer, s string) error {
	length := uint16(len(s))
	if err := binary.Write(w, binary.BigEndian, length); err != nil {
		return fmt.Errorf("failed to write string length: %v", err)
	}

	if _, err := w.Write([]byte(s)); err != nil {
		return fmt.Errorf("failed to write string data: %v", err)
	}

	return nil
}

// ReadNullTerminatedString reads a null-terminated string from a reader
func ReadNullTerminatedString(r io.Reader) (string, error) {
	var bytes []byte
	for {
		var b byte
		if err := binary.Read(r, binary.BigEndian, &b); err != nil {
			return "", fmt.Errorf("failed to read byte: %v", err)
		}
		if b == 0 {
			break
		}
		bytes = append(bytes, b)
	}
	return string(bytes), nil
}

// WriteNullTerminatedString writes a null-terminated string to a writer
func WriteNullTerminatedString(w io.Writer, s string) error {
	if _, err := w.Write([]byte(s)); err != nil {
		return fmt.Errorf("failed to write string: %v", err)
	}
	if err := binary.Write(w, binary.BigEndian, byte(0)); err != nil {
		return fmt.Errorf("failed to write null terminator: %v", err)
	}
	return nil
}

// GetLocalIP returns the local IP address
func GetLocalIP() (string, error) {
	addrs, err := net.InterfaceAddrs()
	if err != nil {
		return "", fmt.Errorf("failed to get interface addresses: %v", err)
	}

	for _, addr := range addrs {
		if ipnet, ok := addr.(*net.IPNet); ok && !ipnet.IP.IsLoopback() {
			if ipnet.IP.To4() != nil {
				return ipnet.IP.String(), nil
			}
		}
	}

	return "", fmt.Errorf("no suitable IP address found")
}
