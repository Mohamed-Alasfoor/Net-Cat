package main

import (
	"log"
	"net-cat/functions" // Import the functions package
	"os"
	"os/signal"
)

func main() {
	// Get the port number from the command-line arguments
	port, err := functions.GetPort()
	if err != nil {
		log.Fatalln("[USAGE]: ./TCPChat $port")
	}

	// Create a new Server instance
	server := &functions.Server{}

	// Start the server with the specified port and a maximum of 10 connections
	err = server.Start(port, 10)
	if err != nil {
		log.Fatalf("ERROR Starting the server %v\n", err) // Log and exit if server fails to start
	}

	// Log the server's IP address and port
	log.Printf("Listening on IP %s and port %s\n", functions.GetLocalIP(), port)

	// Create a channel to receive OS signals
	CHAN := make(chan os.Signal, 1)
	signal.Notify(CHAN, os.Interrupt) // Notify the channel for SIGINT (Ctrl+C) signal

	// Start a goroutine to handle incoming connections
	go func() {
		for {
			// Accept incoming connections
			conn, err := server.Listener.Accept()
			if err != nil {
				break // Exit the loop if there's an error accepting connections
			}
			// Handle the accepted connection in a new goroutine
			go server.HandleConnection(conn)
		}
	}()

	// Wait for a signal to close the server
	<-CHAN

	// Close the server
	server.CloseServer()
}
