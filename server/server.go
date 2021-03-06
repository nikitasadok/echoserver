package server

import (
	"container/heap"
	"echoServer/models"
	"echoServer/server/connectionQueue"
	"github.com/pkg/errors"
	"io"
	"sync"

	"log"
	"net"
	"os"
	"time"
)

type Server interface {
	Listen()
}

type echoServer struct {
	listener        net.Listener
	connQueue       connectionQueue.ConnectionQueue
	idleTimeout     time.Duration
	maxConns        int
	maxReadBytes    int
	currentConns    int
	currentConnsMux sync.Mutex
}

func NewEchoServer(host, port string) (*echoServer, error) {
	if host == "" || port == "" {
		return nil, errors.Wrap(models.ErrCreateServer, "host or port param is empty")
	}
	listener, err := net.Listen("tcp", host+port)
	if err != nil {
		log.Println("Fatal error creating server", err)
		return nil, errors.WithMessage(err, models.ErrCreateServer.Error())
	}

	return &echoServer{
		listener:     listener,
		maxConns:     models.MaxConns,
		maxReadBytes: models.MaxReadBytes,
		idleTimeout:  time.Minute * 3,
		connQueue:    connectionQueue.NewConnectionQueue(),
	}, nil
}

func (s *echoServer) Listen() {
	for {
		if s.connQueue.Len() == s.maxConns {
			s.closeLeastUpdConn()
		}
		conn, err := s.listener.Accept()
		if err != nil {
			log.Println("Error accepting connection", err)
			continue
		}
		c := models.Connection{
			Conn:       conn,
			LastUpdate: time.Now(),
		}
		s.connQueue.Push(&c)
		go s.handleRequest(&c)
	}
}

func (s *echoServer) handleRequest(c *models.Connection) {
	s.currentConnsMux.Lock()
	s.currentConns++
	s.currentConnsMux.Unlock()
	defer func() {
		s.currentConnsMux.Lock()
		s.currentConns--
		s.currentConnsMux.Unlock()
		if err := c.Conn.Close(); err != nil {
			log.Println("Error closing connection", err)
		}
	}()
	for {
		if err := c.Conn.SetReadDeadline(time.Now().Add(s.idleTimeout)); err != nil {
			log.Println("Error setting read deadline", err)
			return
		}
		buf := make([]byte, s.maxReadBytes)
		size, err := c.Conn.Read(buf)
		if err != nil {
			s.handleReadError(err, c.Conn)
			return
		}
		s.connQueue.Update(c, time.Now())
		data := buf[:size]
		if s.isQuit(data) {
			if _, err := c.Conn.Write([]byte(models.MsgQuit)); err != nil {
				log.Println("Error writing quit response", err)
			}
			return
		}
		if _, err := c.Conn.Write(data); err != nil {
			log.Println("Error writing response", err)
			return
		}
	}
}

func (s *echoServer) isQuit(data []byte) bool {
	return string(data) == "quit"
}

func (s *echoServer) closeLeastUpdConn() {
	heap.Init(&s.connQueue)
	c := s.connQueue.Pop()
	conn := c.(*models.Connection)
	if _, err := conn.Conn.Write([]byte(models.MsgTimeout)); err != nil {
		log.Println("Error writing end of stream message", err)
	}
	if err := conn.Conn.Close(); err != nil {
		log.Println("Error closing least active connection", err)
	}
	s.currentConnsMux.Lock()
	s.currentConns--
	s.currentConnsMux.Unlock()
}

func (s *echoServer) handleReadError(err error, conn net.Conn) {
	if err == io.EOF {
		if _, err := conn.Write([]byte(models.MsgEOF)); err != nil {
			log.Println("Error writing EOF response", err)
		}
		return
	}
	if os.IsTimeout(err) {
		if _, err := conn.Write([]byte(models.MsgTimeout)); err != nil {
			log.Println("Error writing timeout response", err)
		}
		return
	}
	log.Println("Not a normal error", err)
}
