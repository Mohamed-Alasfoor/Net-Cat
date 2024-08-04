package functions

import (
	"bufio"
	"errors"
	"fmt"
	"log"
	"net"
	"os"
	"strings"
	"time"
)

// GetPort retrieves the port number from command-line arguments or returns the default port.
func GetPort() (string, error) {
	args := os.Args
	if len(args) < 2 {
		return DefaultPort, nil
	} else if len(args) > 2 {
		return "", errors.New("too many arguments")
	}
	return ":" + args[1], nil
}

// GetLocalIP retrieves the local IP address of the machine.
func GetLocalIP() string {
	conn, err := net.Dial("udp", "8.8.8.8:80")
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()
	return conn.LocalAddr().(*net.UDPAddr).IP.String()
}

// Start initializes the server and starts listening for incoming connections.
func (s *Server) Start(port string, maxConn int) error {
	listener, err := net.Listen("tcp", port)
	if err != nil {
		return err
	}
	s.Listener = listener
	s.MaxConnections = maxConn
	s.Clients = make(map[net.Conn]string)
	s.clientNames = make(map[string]bool)
	return nil
}

// acceptConnection checks if a new connection can be accepted based on the maximum connection limit.
func (s *Server) acceptConnection() bool {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	return len(s.Clients) < s.MaxConnections
}

// HandleConnection handles a new incoming connection.
func (s *Server) HandleConnection(conn net.Conn) {
	if !s.acceptConnection() {
		fmt.Fprint(conn, "The room is full, please try again later.")
		conn.Close()
		return
	}
	fmt.Fprint(conn, WelcomeMessage)
	reader := bufio.NewReader(conn)
	name, err := reader.ReadString('\n')
	if err != nil {
		log.Println("Error reading name:", err)
		conn.Close()
		return
	}
	name = strings.TrimSpace(name)
	if err := s.addConnection(conn, name); err != nil {
		fmt.Fprint(conn, err.Error())
		conn.Close()
		return
	}
	s.startChat(conn)
	s.removeConnection(conn)
}

// startChat handles the chat session for a connected client.
func (s *Server) startChat(conn net.Conn) {
	s.loadMessages(conn)
	message := formatMessageForBroadcast(s, conn, "", ModeJoinChat)
	s.broadcastMessage(conn, message)

	reader := bufio.NewReader(conn)
	for {
		message, err := reader.ReadString('\n')
		if err != nil {
			break
		}
		formattedMessage := formatMessageForBroadcast(s, conn, message, ModeSendMessage)
		s.broadcastMessage(conn, formattedMessage)
		s.saveMessage(formattedMessage)
	}

	message = formatMessageForBroadcast(s, conn, "", ModeLeftChat)
	s.broadcastMessage(conn, message)
}

// broadcastMessage sends a message to all connected clients.
func (s *Server) broadcastMessage(conn net.Conn, message string) {
	if message == "" {
		fmt.Fprintf(conn, PatternSending, time.Now().Format(TimeFormat), s.Clients[conn])
		return
	}
	timestamp := time.Now().Format(TimeFormat)
	formattedMessage := fmt.Sprintf("%s%s\n%s", ColorYellow, message, ColorReset)
	s.mutex.Lock()
	defer s.mutex.Unlock()
	for otherConn := range s.Clients {
		if otherConn != conn {
			fmt.Fprint(otherConn, formattedMessage)
		}
		fmt.Fprintf(otherConn, PatternSending, timestamp, s.Clients[otherConn])
	}
}

// loadMessages sends all previously saved messages to a connected client.
func (s *Server) loadMessages(conn net.Conn) {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	for _, message := range s.AllMessages {
		fmt.Fprint(conn, message)
	}
}

// saveMessage saves a message to the server's message history.
func (s *Server) saveMessage(message string) {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	s.AllMessages = append(s.AllMessages, message)
}

// addConnection adds a new connection to the server's client list.
func (s *Server) addConnection(conn net.Conn, name string) error {
	if name == "" {
		return errors.New("name cannot be empty")
	}
	s.mutex.Lock()
	defer s.mutex.Unlock()
	if _, exists := s.clientNames[name]; exists {
		return errors.New("name already in use")
	}
	s.clientNames[name] = true
	s.Clients[conn] = name
	return nil
}

// removeConnection removes a connection from the server's client list.
func (s *Server) removeConnection(conn net.Conn) {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	name := s.Clients[conn]
	conn.Close()
	delete(s.Clients, conn)
	delete(s.clientNames, name)
}

// CloseServer closes the server and all connected clients.
func (s *Server) CloseServer() {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	for conn := range s.Clients {
		fmt.Fprintf(conn, "\n%sServer is shutting down%s", BgColorRed, ColorReset)
		conn.Close()
	}
	s.Listener.Close()
}
