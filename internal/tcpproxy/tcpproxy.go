package tcpproxy

import (
	"fmt"
	"io"
	"log"
	"net"
	"sync"

	"github.com/google/uuid"
	"github.com/pikami/servicebus-spy/internal/amqphandler"
)

type TcpProxy struct {
	sourcePort         int
	destinationAddress string
	amqpHandler        *amqphandler.AMQPHandler
}

func NewTcpProxy(sourcePort int, destinationAddress string, amqpHandler *amqphandler.AMQPHandler) *TcpProxy {
	return &TcpProxy{
		sourcePort:         sourcePort,
		destinationAddress: destinationAddress,
		amqpHandler:        amqpHandler,
	}
}

func (p *TcpProxy) Start() {
	listener, err := net.Listen("tcp", fmt.Sprintf(":%d", p.sourcePort))
	if err != nil {
		log.Fatalf("Failed to listen on %d: %v", p.sourcePort, err)
	}
	defer listener.Close()

	log.Printf("Listening on %d, forwarding to %s", p.sourcePort, p.destinationAddress)

	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Fatalf("Failed to accept connection: %v", err)
		}

		go p.handleConnection(conn)
	}
}

func (p *TcpProxy) handleConnection(conn net.Conn) {
	connID := uuid.New().String()
	p.amqpHandler.OnNewConnection(connID)
	defer p.amqpHandler.OnConnectionClosed(connID)
	defer conn.Close()

	log.Printf("New connection from %s", conn.RemoteAddr().String())

	serverConn, err := net.Dial("tcp", p.destinationAddress)
	if err != nil {
		log.Fatalf("Failed to connect to %s: %v", p.destinationAddress, err)
	}
	defer serverConn.Close()

	serverToClientWriter := NewMultiWriter(serverConn, p.amqpHandler.GetConnectionWriter(connID))
	clientToServerWriter := NewMultiWriter(conn, p.amqpHandler.GetConnectionWriter(connID))

	var wg sync.WaitGroup
	wg.Add(2)

	go func() {
		defer wg.Done()
		io.Copy(serverToClientWriter, conn)
	}()

	go func() {
		defer wg.Done()
		io.Copy(clientToServerWriter, serverConn)
	}()

	wg.Wait()
	log.Printf("Connection from %s closed", conn.RemoteAddr().String())
}
