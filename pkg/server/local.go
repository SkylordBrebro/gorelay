package server

import (
	"fmt"
	"log"
	"net"
)

type LocalServer struct {
	port     int
	listener net.Listener
	clients  []net.Conn
}

func NewLocalServer(port int) *LocalServer {
	return &LocalServer{
		port:    port,
		clients: make([]net.Conn, 0),
	}
}

func (s *LocalServer) Start() error {
	listener, err := net.Listen("tcp", fmt.Sprintf(":%d", s.port))
	if err != nil {
		return err
	}
	s.listener = listener

	go s.acceptConnections()
	return nil
}

func (s *LocalServer) acceptConnections() {
	for {
		conn, err := s.listener.Accept()
		if err != nil {
			log.Printf("Error accepting connection: %v", err)
			continue
		}

		s.clients = append(s.clients, conn)
		go s.handleConnection(conn)
	}
}

func (s *LocalServer) handleConnection(conn net.Conn) {
	defer conn.Close()

	buffer := make([]byte, 4096)
	for {
		_, err := conn.Read(buffer)
		if err != nil {
			return
		}
		// Handle packet processing here
	}
}
