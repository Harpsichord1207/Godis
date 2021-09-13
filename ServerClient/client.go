package ServerClient

import "net"

type Client struct {
	Conn net.Conn
}