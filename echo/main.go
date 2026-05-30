package main

import (
	"bufio"
	"fmt"
	"log"
	"net"
)

func handleConnection(conn net.Conn) {
	// Ensure the connection closes when this worker finishes
	defer conn.Close()

	fmt.Printf("[SERVER] New client connected from: %s\n", conn.RemoteAddr().String())
	conn.Write([]byte("Welcome to the Go Echo Server! Type anything...\n"))

	// Set up a scanner to read data coming over the network line-by-line
	scanner := bufio.NewScanner(conn)
	for scanner.Scan() {
		text := scanner.Text()

		// If the user types "exit", drop their connection gracefully
		if text == "exit" {
			conn.Write([]byte("Goodbye!\n"))
			break
		}

		// Echo the text back to the client
		reply := fmt.Sprintf("Echo: %s\n", text)
		conn.Write([]byte(reply))
	}

	fmt.Printf("[SERVER] Client %s disconnected.\n", conn.RemoteAddr().String())
}

func main() {
	// 1. Bind to port 8080 and listen for TCP traffic
	listener, err := net.Listen("tcp", ":8080")
	if err != nil {
		log.Fatal("Failed to start server:", err)
	}
	defer listener.Close()
	fmt.Println("[SERVER] Echo server running on port :8080...")

	// 2. Loop forever, waiting for clients to connect
	for {
		conn, err := listener.Accept() // Blocks gere until a client connects
		if err != nil {
			fmt.Println("Failed to accept connection:", err)
			continue
		}

		// 3. LAUNCH THE CONCURRENCY MAGIC!
		// We spawn a worker thread for this specific client, leaving the
		// loop instantly free to accept the NEXT client.
		go handleConnection(conn)
	}
}
