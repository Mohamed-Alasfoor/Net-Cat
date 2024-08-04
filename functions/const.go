package functions

import (
	"net"
	"sync"
)

type Server struct {
	Listener       net.Listener
	Clients        map[net.Conn]string
	clientNames    map[string]bool
	MaxConnections int
	AllMessages    []string
	mutex          sync.Mutex
}

const (
	DefaultPort    = ":8989"
	MaxConnections = 10

	WelcomeMessage  = "Welcome to TCP-Chat!\n         _nnnn_\n        dGGGGMMb\n       @p~qp~~qMb\n       M|@||@) M|\n       @,----.JM|\n      JS^\\__/  qKL\n     dZP        qKRb\n    dZP          qKKb\n   fZP            SMMb\n   HZM            MMMM\n   FqM            MMMM\n __| \".        |\\dS\"qML\n |    `.       | `' \\Zq\n_)      \\.___.,|     .'\n\\____   )MMMMMP|   .'\n     `-'       `--'\n[ENTER YOUR NAME]: "
	PatternSending  = "[%v][%v]:"
	PatternMessage  = "[%v][%s]: %s"
	PatternJoinChat = "%s has joined the chat...\n"
	PatternLeftChat = "%s has left the chat...\n"

	TimeFormat = "2006-01-02 15:04:05"

	ColorReset  = "\u001b[0m"
	ColorYellow = "\u001b[33m"
	BgColorRed  = "\u001b[41m"
	BgColorGray = "\u001b[47;1m"

	ModeSendMessage = iota
	ModeJoinChat
	ModeLeftChat
)
