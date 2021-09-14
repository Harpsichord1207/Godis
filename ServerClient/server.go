package ServerClient

import (
	"Godis/lib/IOUtils"
	"context"
	"fmt"
	"io"
	"log"
	"net"
	"os"
	"os/signal"
	"sync"
	"syscall"
)

type Handler interface {
	Handle(ctx context.Context, conn net.Conn)
	io.Closer
}

type Config struct {
	Address string
}

func ListenAndServe(listener net.Listener, handler Handler, closeChan <-chan struct{}) {
	// 如果closeChan中读取到数据就关闭Server
	go func() {
		<-closeChan
		log.Println("Shutdown Server...")
		IOUtils.Close(listener)
		IOUtils.Close(handler)
	}()

	defer func() {
		IOUtils.Close(listener)
		IOUtils.Close(handler)
	}()

	ctx := context.Background()
	var wg sync.WaitGroup

	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Println(fmt.Sprintf("Failed to accept %#v, error is : %s", listener, err.Error()))
			break
		} else {
			log.Println(fmt.Sprintf("Accept %#v", listener))
		}
		// 每次Accept的conn处理完才会进行下一次循环
		wg.Add(1)

		go func() {
			defer wg.Done()
			handler.Handle(ctx, conn)
		}()
	}
	wg.Wait()
}

func ListenAndServeWithSignal(config *Config, handler Handler) {
	closeChan := make(chan struct{})
	signalChan := make(chan os.Signal)
	// 收到指定信号时发送给signalChan
	signal.Notify(signalChan, syscall.SIGHUP, syscall.SIGQUIT, syscall.SIGTERM, syscall.SIGINT)

	go func() {
		sig := <-signalChan
		log.Println(fmt.Sprintf("Received Signal %#v", sig))
		closeChan <- struct{}{}
	}()

	listener, err := net.Listen("tcp", config.Address)
	if err != nil {
		log.Fatal(fmt.Sprintf("Failed to listen %s, error is: %s", config.Address, err.Error()))
	} else {
		log.Println(fmt.Sprintf("bind to %s, start listening", config.Address))
	}
	ListenAndServe(listener, handler, closeChan)
}
