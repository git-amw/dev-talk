package chatapp

import (
	"bufio"
	"log/slog"
	"net"
	"strings"
	"sync"
)

type ChatServer struct {
	rooms              map[string]*Room
	registerClientCh   chan *Client
	unregisterClientCh chan *Client
	connectionQueue    chan net.Conn
	mutx               sync.Mutex
	wg                 sync.WaitGroup
}

func NewChatServer(poolSize int) *ChatServer {
	return &ChatServer{
		rooms:              make(map[string]*Room),
		registerClientCh:   make(chan *Client),
		unregisterClientCh: make(chan *Client),
		connectionQueue:    make(chan net.Conn, poolSize),
	}
}

func (srv *ChatServer) StartServer(port string, poolSize int) {
	listener, err := net.Listen("tcp", ":7000")
	if err != nil {
		slog.Error("Error in listener", slog.Any("err", err))
	}
	defer listener.Close()
	slog.Info("Server is up and running at :7000")

	srv.startWorkerPool(poolSize)
	go srv.hangleClientRegistration()

	for {
		conn, err := listener.Accept()
		if err != nil {
			slog.Error("Error in Accepting new client", slog.Any("err", err))
		}
		slog.Info("Client", slog.Any("addr ", conn.RemoteAddr()))
		srv.addToQueue(conn)
	}
}

func (srv *ChatServer) startWorkerPool(poolSize int) {
	for i := 0; i < poolSize; i++ {
		go func() {
			for conn := range srv.connectionQueue {
				srv.HandleNewClient(conn)
			}
		}()
	}
}

func (srv *ChatServer) addToQueue(conn net.Conn) {
	srv.wg.Add(1)
	srv.connectionQueue <- conn
}

func (srv *ChatServer) hangleClientRegistration() {
	for {
		select {
		case client := <-srv.registerClientCh:
			client.conn.Write([]byte("Welcome, " + client.name + "! Type /rooms to see active rooms.\r\n"))
			client.conn.Write([]byte("To interact with app use these commands:\r\n"))
			client.conn.Write([]byte("1. /rooms - List available rooms\r\n"))
			client.conn.Write([]byte("2. /create [room name] - Create a room\r\n"))
			client.conn.Write([]byte("3. /join [room name] - Join a room\r\n"))
			client.conn.Write([]byte("4. /quit - Quit app\r\n"))
		case client := <-srv.unregisterClientCh:
			client.conn.Write([]byte("Disconnecting " + client.name + "\r\n"))
			client.conn.Close()
		}
	}
}

func (srv *ChatServer) HandleNewClient(conn net.Conn) {
	conn.Write([]byte("Enter your name: \r\n"))
	name, _ := bufio.NewReader(conn).ReadString('\n')
	name = strings.TrimSpace(name)
	client := &Client{
		name:      name,
		conn:      conn,
		messageCh: make(chan string),
		server:    srv,
	}
	srv.registerClientCh <- client
	go client.listenForMessages()
	client.readInput()
	defer close(client.messageCh)
	defer srv.wg.Done()

}
