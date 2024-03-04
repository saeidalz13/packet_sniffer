package main

import (
	"fmt"
	"net"
)

type TcpServer struct {
	ListenAddr string
	Listener   net.Listener
	Quitch     chan struct{}
	MsgChan    chan []byte
}

func NewTcpServer(listenAddr string) *TcpServer {
	return &TcpServer{
		ListenAddr: listenAddr,
		Quitch:     make(chan struct{}),
		// Needs to be buffered otherwise it will go on forever
		MsgChan:    make(chan []byte, 10),
	}
}

func (t *TcpServer) Start() error {
	listener, err := net.Listen("tcp", t.ListenAddr)
	if err != nil {
		return err
	}
	defer listener.Close()
	t.Listener = listener

	go t.acceptLoop()

	<-t.Quitch
	// to say party is over
	close(t.MsgChan)
	
	return nil
}

func (t *TcpServer) acceptLoop() {
	for {
		conn, err := t.Listener.Accept()
		if err != nil {
			fmt.Println("accept errorL:", err)
			continue
		}

		fmt.Println("new connection made...", conn.RemoteAddr())
		go t.readloop(conn)
	}
}

func (t *TcpServer) readloop(conn net.Conn) {
	defer conn.Close()
	buf := make([]byte, 2048)

	for {
		n, err := conn.Read(buf)
		if err != nil {
			fmt.Println("read error:", err)
			continue
		}
		t.MsgChan <- buf[:n]
	}
}

func main() {
	tcpServer := NewTcpServer(":3000")
	if err := tcpServer.Start(); err != nil {
		panic(err)
	}

}
