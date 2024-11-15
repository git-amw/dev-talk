# dev-talk a Real time TCP Chat Application

This is a simple concurrent chat application built with Go, allowing multiple clients to join chat rooms, set their usernames, and communicate with others within specific rooms. The server supports basic functionality like listing active rooms, number of members in a room, creating or joining a room, and broadcasting messages within a room.

## Features

> **Concurrent Client Handling**: Uses goroutines and channels to handle multiple clients simultaneously.

> **Thread Pool Support**: The server uses a thread pool to handle client connections and requests, ensuring a limited number of concurrent goroutines.

> **Room Management**: Allows clients to list available rooms, create new rooms, or join existing rooms.

> **Broadcasting**: Messages are broadcasted to all clients within a specific room.

## Commands

| Command               | Description                                                      |
| --------------------- | ---------------------------------------------------------------- |
| `/rooms`              | List all _rooms_ availabe on server and number of online members |
| `/create [room name]` | Create a room                                                    |
| `/join [room name]`   | Join a room with give _name_                                     |
| `/quit`               | Quit app                                                         |
| `/msg`                | To send message in the room                                      |
| `/leave`              | Leave a room                                                     |

## Build and Run application

```
$ go build -o chat-app main.go

$ ./chat-app

```
