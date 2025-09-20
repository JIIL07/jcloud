package socket

import (
	"sync"

	"github.com/gorilla/websocket"
)

var clients = make(map[*websocket.Conn]bool)
var mutex = &sync.Mutex{}

func Broadcast(message string) error {
	mutex.Lock()
	defer mutex.Unlock()

	for client := range clients {
		err := client.WriteMessage(websocket.TextMessage, []byte(message))
		if err != nil {
			client.Close()
			delete(clients, client)
			return err
		}
	}
	return nil
}
