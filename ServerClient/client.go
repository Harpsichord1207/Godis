package ServerClient

import (
	"Godis/lib/IOUtils"
	"Godis/lib/wait"
	"bufio"
	"context"
	"fmt"
	"io"
	"log"
	"net"
	"sync"
	"sync/atomic"
	"time"
)

type Client struct {
	Conn net.Conn
	Wait wait.Wait
}

type EchoHandler struct {
	activeConn sync.Map
	closing atomic.Value
}

func BuildEchoHandler() *EchoHandler {
	return &EchoHandler{}
}

func (e *EchoHandler) Handle (ctx context.Context, conn net.Conn) {
	if e.closing.Load()!= nil {
		IOUtils.Close(conn)
	}

	client := &Client{Conn: conn}
	e.activeConn.Store(client, struct{}{})  // Map当Set用

	reader := bufio.NewReader(conn)

	for {
		data, err := reader.ReadString('\n')
		if err == io.EOF {
			log.Println(fmt.Sprintf("Received EOF, close connection %#v.", conn))
			break
		} else if err != nil {
			log.Println(err.Error())
			break
		} else {
			IOUtils.Write(conn, []byte(data))
			client.Wait.Done()
		}
	}
}

func (c *Client) Close() error {
	c.Wait.WaitFinishOrTimeOut(10 * time.Second)
	return c.Conn.Close()
}

func (e *EchoHandler) Close() error {
	log.Println("Start to shutdown EchoHandler...")
	e.closing.Store(1)

	e.activeConn.Range(func(key, value interface{}) bool {
		c := key.(*Client)
		IOUtils.Close(c)
		return true
	})
	log.Println("EchoHandler shutdown.")
	return nil
}
