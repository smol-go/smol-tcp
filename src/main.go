package main

import (
	"fmt"
	"log"
	"net"
)

type Server struct {
	listenAddress string
	ln            net.Listener
	quitch        chan struct{}
	msgChan       chan []byte
}

func NewServer(listenAddress string) *Server {
	return &Server{
		listenAddress: listenAddress,
		quitch:        make(chan struct{}),
		msgChan:       make(chan []byte, 10),
	}
}

func (server *Server) Start() error {
	ln, err := net.Listen("tcp", server.listenAddress)
	if err != nil {
		return err
	}
	defer ln.Close()
	server.ln = ln
	log.Println("Server started on port", server.listenAddress)

	go server.acceptLoop()

	<-server.quitch
	close(server.msgChan)

	return nil
}

func (server *Server) acceptLoop() {
	for {
		conn, err := server.ln.Accept()
		if err != nil {
			fmt.Println("accept error:", err)
			continue
		}

		fmt.Println("new connection to the server:", conn.RemoteAddr())

		go server.readLoop(conn)
	}
}

func (server *Server) readLoop(conn net.Conn) {
	defer conn.Close()
	buf := make([]byte, 2048)

	for {
		n, err := conn.Read(buf)
		if err != nil {
			log.Println("read error:", err)
			continue
		}

		server.msgChan <- buf[:n]
	}
}

func main() {
	server := NewServer(":3000")

	go func() {
		for msg := range server.msgChan {
			fmt.Println("received msaage for connection:", string(msg))
		}
	}()

	log.Fatal(server.Start())
}
