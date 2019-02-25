package main

import (
	"bufio"
	"flag"
	"fmt"
	"github.com/containous/whoamitcp/tcp"
	"log"
	"net"
	"strings"
	"time"
)

var port string

func init() {
	flag.StringVar(&port, "port", "80", "give me a port number")
}

func main() {
	flag.Parse()

	fmt.Println("Starting up on port " + port)
	listener, err := buildListener(port)
	if err != nil {
		log.Fatal("could not start server: %v", err)
		return
	}

	handler := buildDefaultHandler()
	for {
		conn, err := listener.Accept()
		log.Printf("Connection arrived")
		if err != nil {
			panic(err)
		}
		log.Printf("Forwarding to handler")
		go handler.ServeTCP(conn)
	}
}

func buildListener (port string) (net.Listener, error) {
	listener, err := net.Listen("tcp", port)
	if err != nil {
		return nil, fmt.Errorf("error opening listener: %v", err)
	}

	return tcpKeepAliveListener{listener.(*net.TCPListener)}, nil
}

// tcpKeepAliveListener sets TCP keep-alive timeouts on accepted
// connections.
type tcpKeepAliveListener struct {
	*net.TCPListener
}

func (ln tcpKeepAliveListener) Accept() (net.Conn, error) {
	tc, err := ln.AcceptTCP()
	if err != nil {
		return nil, err
	}

	if err = tc.SetKeepAlive(true); err != nil {
		return nil, err
	}

	if err = tc.SetKeepAlivePeriod(3 * time.Minute); err != nil {
		return nil, err
	}

	return tc, nil
}

func buildDefaultHandler() tcp.Handler {
	handler := tcp.Handler{}
	handler.ServeTCP = func (conn net.Conn) {
		defer conn.Close()

		log.Printf("Serving %s\n", conn.RemoteAddr().String())

		for {
			netData, err := bufio.NewReader(conn).ReadString('\n')
			if err != nil {
				fmt.Println(err)
				return
			}

			temp := strings.TrimSpace(string(netData))
			if temp == "STOP" {
				break
			}

			result := "Me me me \n"
			conn.Write([]byte(result))
		}

	}
   return handler
}