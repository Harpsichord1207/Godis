package Echo

import (
	"Godis/lib/IOUtils"
	"bufio"
	"fmt"
	"io"
	"log"
	"net"
)


func ListenAndServe(addr string) {
	listener, err := net.Listen("tcp", addr)
	if err != nil {
		log.Fatal(fmt.Sprintf("Failed to listen %s, error is: %s", addr, err.Error()))
	} else {
		log.Println(fmt.Sprintf("bind to %s, start listening", addr))
	}

	defer IOUtils.Close(listener)

	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Fatal(fmt.Sprintf("Failed to accept %#v, error is : %s", listener, err.Error()))
		}
		go Handle(conn)
	}
}

func Handle(conn net.Conn) {
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
		}
	}
}
