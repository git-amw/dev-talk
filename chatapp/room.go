package chatapp

import (
	"bufio"
	"log/slog"
	"net"
	"strings"
)

type Room struct {
	name    string
	members map[net.Addr]*Client
}

func (r *Room) roomCommands(client *Client) {
	client.conn.Write([]byte("## Room Command ## \r\n"))
	client.conn.Write([]byte("1. /msg [your message] To send message \r\n"))
	client.conn.Write([]byte("2. /leave Leave current room \r\n"))
	for {
		msg, err := bufio.NewReader(client.conn).ReadString('\n')
		if err != nil {
			slog.Error("Error in reading room cmd", slog.Any("err ", err))
		}
		if r.handleRoomCommand(client, strings.TrimSpace(msg)) {
			return
		}
	}
}

func (r *Room) handleRoomCommand(client *Client, msg string) bool {
	parts := strings.SplitN(msg, " ", 2)
	command := parts[0]
	switch command {
	case "/msg":
		r.broadCastMessage(client, parts[1])
		return false
	case "/leave":
		r.removeMember(client)
		return true
	default:
		client.conn.Write([]byte("Enter a valid command\r\n"))
		return false
	}
}

func (r *Room) removeMember(client *Client) {
	delete(r.members, client.conn.RemoteAddr())
	client.room = nil
	client.server.mutx.Lock()
	if len(r.members) == 0 {
		delete(client.server.rooms, r.name)
	}
	client.server.mutx.Unlock()
	client.conn.Write([]byte("You've left the " + r.name + "\r\n"))
	r.broadCastMessage(client, client.name+" has left the room.\r\n")
}

func (r *Room) broadCastMessage(sender *Client, msg string) {
	for addr, client := range r.members {
		if sender.conn.RemoteAddr() != addr {
			client.messageCh <- msg
		}
	}
}
