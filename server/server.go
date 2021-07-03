package server

import (
	"echoServer/models"
	"echoServer/server/connectionQueue"
	"github.com/pkg/errors"

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
	connQueue    connectionQueue.ConnectionQueue
	idleTimeout  time.Duration
	maxConns     int
	maxReadBytes int
	currentConns int
}

func NewEchoServer(host, port string) (*EchoServer, error) {
	if host == "" || port == "" {
		return nil, errors.Wrap(models.ErrCreateServer, "host or port param is empty")
	}
	listener, err := net.Listen("tcp", host+port)
	if err != nil {
		log.Println("Fatal error creating server", err)
		return nil, errors.WithMessage(err, models.ErrCreateServer.Error())
	}

	return &EchoServer{
		listener:     listener,
		maxConns:     500000,
		maxReadBytes: 8192,
		idleTimeout:  time.Second * 30,
		connQueue:    connectionQueue.NewConnectionQueue(),
	}, nil
}

func (s *EchoServer) Listen() {
	for {
		if s.currentConns == s.maxConns {
			s.closeLeastUpdConn()
		}
		log.Println("currentConns:", s.currentConns)
		s.currentConns++
		conn, err := s.listener.Accept()
		if err != nil {
			log.Panicln(err)
		}
		c := models.Connection{
			Conn:       conn,
			LastUpdate: time.Now(),
		}
		s.connQueue.Push(&c)
		go s.handleRequest(&c)
	}
}

func (s *EchoServer) handleRequest(c *models.Connection) {
	log.Printf("Accepted new connection: %d\n", c.Index)
	defer func() {
		s.currentConns--
		if err := c.Conn.Close(); err != nil {
			log.Println("Error closing connection", err)
		}
	}()
	for {
		if err := c.Conn.SetReadDeadline(time.Now().Add(s.idleTimeout)); err != nil {

		}
		buf := make([]byte, s.maxReadBytes)
		size, err := c.Conn.Read(buf)
		if os.IsTimeout(err) {
			if _, err := c.Conn.Write([]byte("Exit due to idle.\n")); err != nil {

			}
			return
		}
		s.connQueue.Update(c, time.Now())
		data := buf[:size]
		if s.isQuit(data) {
			if _, err := c.Conn.Write([]byte("Got quit signal. Aborting.\n")); err != nil {

			}
			return
		}
		log.Printf("Read new data from connection: %s, date: %d\n", string(data), c.Index)
		if _, err := c.Conn.Write(data); err != nil {
			log.Println("Error writing response", err)
			return
		}
	}
}

func (s *EchoServer) isQuit(data []byte) bool {
	return string(data) == "quit"
}

func (s *EchoServer) closeLeastUpdConn() {
	c := s.connQueue.Pop()
	conn := c.(*models.Connection)
	if _, err := conn.Conn.Write([]byte("Exit due to idle\n")); err != nil {
		log.Println("Error writing end of stream message", err)
	}
	if err := conn.Conn.Close(); err != nil {
		log.Println("Error closing least active connection", err)
	}
	s.currentConns--
}
