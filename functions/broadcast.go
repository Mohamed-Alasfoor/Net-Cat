package functions

import (
    "fmt"
    "net"
    "time"
)

// getClientName retrieves the name of the client associated with the given connection
// from the server's Clients map. It locks the mutex to ensure thread safety.
func getClientName(server *Server, conn net.Conn) string {
    server.mutex.Lock()
    defer server.mutex.Unlock()
    return server.Clients[conn]
}

// formatMessageForBroadcast formats the message based on the specified mode
// and returns the formatted string to be broadcast to all clients.
func formatMessageForBroadcast(server *Server, conn net.Conn, message string, mode int) string {
    name := getClientName(server, conn) // Get the name of the client sending the message

    switch mode {
    case ModeSendMessage:
        if message == "\n" { // If the message is just a newline, return an empty string
            return ""
        }
        fallthrough // Fallthrough to the default case
    default:
        timestamp := time.Now().Format(TimeFormat) // Get the current timestamp formatted according to TimeFormat
        return fmt.Sprintf(PatternMessage, timestamp, name, message) // Format the message with timestamp and name
    case ModeJoinChat:
        return fmt.Sprintf(ColorYellow+PatternJoinChat+ColorReset, name) // Format the "joined chat" message with the client's name
    case ModeLeftChat:
        return fmt.Sprintf(ColorYellow+PatternLeftChat+ColorReset, name) // Format the "left chat" message with the client's name
    }
}
