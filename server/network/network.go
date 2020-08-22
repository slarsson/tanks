package network

import (
	"fmt"
	"sync"

	"github.com/gorilla/websocket"
)

type Network struct {
	Clients    map[*Client]bool
	Register   chan *Client
	Unregister chan *Client
	Broadcast  chan []byte
	mutex      *sync.RWMutex
}

type Client struct {
	Conn          *websocket.Conn
	NetworkInput  chan []byte
	NetworkOutput chan []byte
}

type Action struct {
	MessageType int8
	Client      *Client
	Payload     []byte
}

func NewNetwork() *Network {
	return &Network{
		Clients:    make(map[*Client]bool),
		Register:   make(chan *Client),
		Unregister: make(chan *Client),
		Broadcast:  make(chan []byte),
		mutex:      &sync.RWMutex{},
	}
}

func (n *Network) Start(ca chan *Action) {
	for {
		select {
		case client := <-n.Register:
			n.mutex.Lock()
			n.Clients[client] = true
			n.mutex.Unlock()
			go n.reader(client, ca)
			go n.writer(client)

		case client := <-n.Unregister:
			err := client.Conn.Close()
			if err != nil {
				fmt.Println("NETWORK: failed to close client")
				continue
			}
			n.mutex.Lock()
			delete(n.Clients, client)
			n.mutex.Unlock()
			client.NetworkOutput <- []byte{} // send stop signal
			ca <- &Action{MessageType: 1, Client: client}

		case message := <-n.Broadcast:
			for client := range n.Clients {
				client.NetworkOutput <- message
			}
		}
	}
}

func (n *Network) reader(client *Client, ca chan *Action) {
	defer func() {
		// when sending fails, trigger unregister player event
		n.Unregister <- client
	}()

	for {
		_, message, err := client.Conn.ReadMessage()
		if err != nil {
			break // kill the goroutine
		}

		if len(message) == 0 {
			continue
		}

		// if MessageType = 1 send to NetworkInput (game state)
		// else send to action channel
		// TODO: this should not be done here..
		if message[0] == 1 {
			client.NetworkInput <- message
		} else {
			ca <- &Action{MessageType: int8(message[0]), Payload: message[1:], Client: client}
		}
	}
}

func (n *Network) writer(client *Client) {
	for {
		select {
		case message, ok := <-client.NetworkOutput:
			if ok {
				if len(message) == 0 {
					break
				}
				client.Conn.WriteMessage(websocket.BinaryMessage, message)
			} else {
				fmt.Println("NETWORK: could not write message, disconnect")
				n.Unregister <- client
				break
			}
		}
	}
}
