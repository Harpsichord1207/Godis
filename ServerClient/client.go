package ServerClient

import (
	"Godis/lib/wait"
	"net"
	"sync"
	"sync/atomic"
)

type Client struct {
	Conn net.Conn
	Wait wait.Wait
}

type EchoHandler struct {
	activeConn sync.Map
	closing atomic.Value
}