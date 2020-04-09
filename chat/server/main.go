package main

import (
	"bufio"
	"errors"
	"fmt"
	"log"
	"net"
	"os"
	"strings"

	"github.com/01walid/goarabic"
)

const (
	defaultPort = "3000"

	cmdQuit = "QUIT"
	cmdEcho = "ECHO"
)

var (
	errDisconnected = errors.New("disconnected")

	makhadnahashFelKlass = arabic("ماخدناهاش في الكيلاس")
)

func init() { log.SetFlags(0) }

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = defaultPort
	}

	log.Printf("Listening on %s...", port)
	ln, err := net.Listen("tcp", fmt.Sprintf("127.0.0.1:%s", port))
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}
	defer ln.Close()

	for {
		conn, err := ln.Accept()
		if err != nil {
			log.Fatalf("Failed to accept: %v", err)
		}

		// TODO: Add error handling using channels?
		log.Printf("Accepted %v", conn.RemoteAddr())
		go func() {
			if err := handle(conn); err != nil && err != errDisconnected {
				log.Fatalf("Failed to handle: %v", err)
			}
			log.Printf("Disconnected from %v", conn.RemoteAddr())
		}()
	}
}

func handle(conn net.Conn) error {
	defer conn.Close()

	fmt.Fprint(conn, "Hello, how are you? Hope you're safe =)\n")

	sc := bufio.NewScanner(conn)
	for {
		if !sc.Scan() {
			return errDisconnected
		}
		args := strings.Fields(strings.TrimSpace(sc.Text()))
		if len(args) == 0 {
			fmt.Fprint(conn, "I'm not sure what was that. Try again?\n")
			continue
		}
		cmd, args := strings.ToUpper(args[0]), args[1:]
		switch cmd {
		case cmdQuit:
			fmt.Fprint(conn, "Ok, bye-bye! :)\n")
			return nil
		case cmdEcho:
			fmt.Fprintf(conn, "You said: %s\n", strings.Join(args, " "))
		default:
			fmt.Fprintf(conn, "I'm sorry, %s: %s %s\n", makhadnahashFelKlass, cmd, strings.Join(args, " "))
		}
	}
}

func arabic(s string) string {
	runes := []rune(goarabic.ToGlyph(s))
	for i, j := 0, len(runes)-1; i < j; i, j = i+1, j-1 {
		runes[i], runes[j] = runes[j], runes[i]
	}
	return string(runes)
}
