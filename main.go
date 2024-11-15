package main

import (
	"github.com/git-amw/devtalk/chatapp"
)

func main() {
	poolsize := 4 // set the maxm number of users can access at a given time
	srv := chatapp.NewChatServer(poolsize)
	srv.StartServer("7000", poolsize)
}
