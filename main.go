package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"net"
	"os"
	"strconv"
	"time"
)

var semaphore chan struct{}

func main() {
	semaphore = make(chan struct{}, 500000)
	port := flag.Int("port", 3333, "Port to accept connection on")
	host := flag.String("host", "127.0.0.1", "Host or IP to bind to")

	l, err := net.Listen("tcp", *host+":"+strconv.Itoa(*port))
	if err != nil {
		log.Panicln(err)
	}

	defer l.Close()

	for {
		semaphore <- struct{}{}
		conn, err := l.Accept()
		if err != nil {
			log.Panicln(err)
		}

		// 	conn.SetReadDeadline(time.Now().Add(time.Second * 10))
		go handleRequest(conn)
	}
}

func handleRequest(conn net.Conn) {
	fmt.Println("is called")
	log.Println("Accepted new connection")
	defer conn.Close()
	for {
		conn.SetReadDeadline(time.Now().Add(time.Second * 3))
		buf := make([]byte, 8192)
		size, err := conn.Read(buf)
		if os.IsTimeout(err) {
			conn.Write([]byte("Exit due to idle.\n"))
			return
		}
		data := buf[:size]
		log.Println("Read new data from connection", string(data))
		conn.Write(data)
	}
	<-semaphore
}

func readWithTimeout(ctx *context.Context, conn net.Conn) {
	buf := make([]byte, 8192)
	size, err := conn.Read(buf)
	if err != nil {
		return
	}
	data := buf[:size]
	log.Println("Read new data from connection", string(data))
	conn.Write(data)
}
