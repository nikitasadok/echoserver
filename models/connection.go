package models

import (
	"net"
	"time"
)

type Connection struct {
	Conn       net.Conn
	LastUpdate time.Time
	Index      int
}
