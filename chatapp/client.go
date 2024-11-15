package chatapp

import (
	"bufio"
	"log/slog"
	"net"
	"strconv"
	"strings"
)

type Client struct {
	name      string
	conn      net.Conn
	room      *Room
	messageCh chan string
	server    *ChatServer
}

func (client *Client) readInput() {
	for {
		msg, err := bufio.NewReader(client.conn).ReadString('\n')
		if err != nil { // work later on how to improve
			slog.Error("Error in reading input", slog.Any("err ", err))
			client.server.unregisterClientCh <- client
			return
		}
		client.handleCommand(strings.TrimSpace(msg))
	}
}

func (client *Client) handleCommand(input string) {
	parts := strings.SplitN(input, " ", 2)
	command := parts[0]
	switch command {
	case "/rooms":
		client.listRooms()
	case "/create":
		client.createRoom(parts[1])
	case "/join":
		if len(parts) < 2 {
			client.conn.Write([]byte("Please enter a room name.\r\n"))
			return
		}
		client.joinRoom(parts[1])
	case "/quit":
		client.server.unregisterClientCh <- client
	default:
		client.conn.Write([]byte("Enter a valid command\r\n"))

	}
}

func (client *Client) listRooms() {
	client.server.mutx.Lock()
	defer client.server.mutx.Unlock()
	if len(client.server.rooms) == 0 {
		client.conn.Write([]byte("No active rooms.\r\n"))
		return
	}
	client.conn.Write([]byte("Active rooms:\r\n"))
	for name, r := range client.server.rooms {
		client.conn.Write([]byte(" - " + name + " online :" + strconv.Itoa(len(r.members)) + "\r\n"))
	}
}

func (client *Client) createRoom(roomName string) {
	client.server.mutx.Lock()
	_, exists := client.server.rooms[roomName]
	if !exists {
		room := &Room{
			name:    roomName,
			members: make(map[net.Addr]*Client),
		}
		client.server.rooms[roomName] = room
		client.conn.Write([]byte("Created room - " + roomName + " !!\r\n"))
	} else {
		client.conn.Write([]byte("Room - " + roomName + "aleady exists\r\n"))
	}
	client.server.mutx.Unlock()
}

func (client *Client) joinRoom(roomName string) {
	client.server.mutx.Lock()
	room, exists := client.server.rooms[roomName]
	client.server.mutx.Unlock()
	if !exists {
		client.conn.Write([]byte("Room - " + roomName + " does not exists\r\n"))
		return
	}

	room.members[client.conn.RemoteAddr()] = client
	client.room = room
	room.broadCastMessage(client, client.name+" has joined the room.\r\n")
	client.conn.Write([]byte("------  Joined - " + roomName + "------\r\n"))
	room.roomCommands(client)
}

func (client *Client) listenForMessages() {
	for msg := range client.messageCh {
		client.conn.Write([]byte(msg + "\r\n"))
	}
}
