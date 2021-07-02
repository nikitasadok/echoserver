package server

import (
	"fmt"
	"log"
	"net"
	"os"
	"time"
)

type Server interface {
	Listen()
}

type EchoServer struct {
	listener     net.Listener
	conns        chan struct{}
	idleTimeout  time.Duration
	maxConns     int
	maxReadBytes int
	currentConns int
}

func NewEchoServer(host, port string) *EchoServer {
	listener, err := net.Listen("tcp", host+port)
	if err != nil {
		log.Panicln(err)
	}

	return &EchoServer{
		listener:    listener,
		conns:       make(chan struct{}),
		maxConns:    500000,
		idleTimeout: time.Second * 30,
	}
}

func (s *EchoServer) Listen() {
	for {
		if s.currentConns == s.maxConns {
			// TODO add finding space
		}
		s.currentConns++
		s.conns <- struct{}{}
		conn, err := s.listener.Accept()
		if err != nil {
			log.Panicln(err)
		}
		go s.handleRequest(conn)
	}
}

func (s *EchoServer) handleRequest(conn net.Conn) {
	fmt.Println("is called")
	log.Println("Accepted new connection")
	defer conn.Close()
	for {
		conn.SetReadDeadline(time.Now().Add(s.idleTimeout))
		buf := make([]byte, s.maxReadBytes)
		size, err := conn.Read(buf)
		if os.IsTimeout(err) {
			conn.Write([]byte("Exit due to idle.\n"))
			return
		}
		data := buf[:size]
		// todo check for quit signal
		log.Println("Read new data from connection", string(data))
		conn.Write(data)
	}
	<-s.conns
}
