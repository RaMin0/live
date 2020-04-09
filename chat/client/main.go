package main

import (
	"bufio"
	"errors"
	"fmt"
	"log"
	"net"
	"os"
	"strings"
)

const (
	defaultServerPort = "3000"

	cmdQuit = "QUIT"
)

var (
	errClientDisconnected = errors.New("client disconnected")
	errServerDisconnected = errors.New("server disconnected")
)

func init() { log.SetFlags(0) }

func main() {
	serverPort := os.Getenv("SERVER_PORT")
	if serverPort == "" {
		serverPort = defaultServerPort
	}

	log.Printf("Connecting to %s...", serverPort)
	conn, err := net.Dial("tcp", fmt.Sprintf("127.0.0.1:%s", serverPort))
	if err != nil {
		log.Fatalf("Failed to connect: %v", err)
	}
	log.Printf("Connected to %v...", conn.RemoteAddr())

	if err := handle(conn); err != nil {
		log.Fatalf("Failed to handle: %v", err)
	}
}

func handle(conn net.Conn) error {
	var (
		connSc  = bufio.NewScanner(conn)
		stdinSc = bufio.NewScanner(os.Stdin)
	)

	for {
		if !connSc.Scan() {
			return errServerDisconnected
		}
		log.Printf("< %s", connSc.Text())
		fmt.Print("> ")
		if !stdinSc.Scan() {
			return errClientDisconnected
		}
		cmd := stdinSc.Text()
		fmt.Fprintln(conn, cmd)
		if strings.HasPrefix(strings.ToUpper(cmd), cmdQuit) {
			return nil
		}
	}
}
