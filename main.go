package main

import (
	"github.com/git-amw/devtalk/chatapp"
)

func main() {
	srv := chatapp.NewChatServer()
	srv.StartServer("7000")
}
