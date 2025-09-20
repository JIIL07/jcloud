package socket

import (
	"github.com/gorilla/websocket"
	"sync"
)

var clients = make(map[*websocket.Conn]bool)
var mutex = &sync.Mutex{}

func Broadcast(message string) error {
	mutex.Lock()
	defer mutex.Unlock()

	for client := range clients {
		err := client.WriteMessage(websocket.TextMessage, []byte(message))
		if err != nil {
			client.Close() // nolint:errcheck
			delete(clients, client)
			return err
		}
	}
	return nil
}
