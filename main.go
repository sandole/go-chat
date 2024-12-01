package main

import (
	"bufio"
	"fmt"
	"log"
	"net"
)

type Message struct {
	sender  string
	content string
}

type Server struct {
	clients   map[net.Conn]string
	broadcast chan Message
}

func NewServer() *Server {
	return &Server{
		clients:   make(map[net.Conn]string),
		broadcast: make(chan Message),
	}
}

func (s *Server) Start() {
	listener, err := net.Listen("tcp", ":8080")
	if err != nil {
		log.Fatalf("Unable to start server: %v", err)
	}
	defer listener.Close()

	// Handle broadcasts
	go func() {
		for msg := range s.broadcast {
			// Send message to all clients
			for conn := range s.clients {
				fmt.Fprintf(conn, "%s: %s\n", msg.sender, msg.content)
			}
		}
	}()

	log.Println("Chat server running on :8080")

	// Accept connections
	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Printf("Failed to accept connection: %v", err)
			continue
		}
		go s.handleConnection(conn)
	}
}

func (s *Server) handleConnection(conn net.Conn) {
	defer conn.Close()

	// Get username
	fmt.Fprintf(conn, "Enter your name: ")
	scanner := bufio.NewScanner(conn)
	scanner.Scan()
	username := scanner.Text()

	s.clients[conn] = username
	s.broadcast <- Message{sender: "Server", content: fmt.Sprintf("%s joined the chat", username)}

	// Handle messages
	for scanner.Scan() {
		msg := scanner.Text()
		s.broadcast <- Message{sender: username, content: msg}
	}

	// Client disconnected
	delete(s.clients, conn)
	s.broadcast <- Message{sender: "Server", content: fmt.Sprintf("%s left the chat", username)}
}

func main() {
	server := NewServer()
	server.Start()
}
