package main

import (
	"echoServer/server"
)

func main() {
	serv, err := server.NewEchoServer("127.0.0.1", ":3333")
	if err != nil {
		return
	}
	serv.Listen()
}
