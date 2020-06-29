package network

import (
	"github.com/gorilla/websocket"
)

type Network struct {
	Clients    map[*Client]bool
	Register   chan *Client
	Unregister chan *Client
	Broadcast  chan []byte
}

type Client struct {
	Conn          *websocket.Conn
	NetworkInput  chan []byte
	NetworkOutput chan []byte
}

type Action struct {
	MessageType int8
	Client      *Client
}

func NewNetwork() *Network {
	return &Network{
		Clients:    make(map[*Client]bool),
		Register:   make(chan *Client),
		Unregister: make(chan *Client),
		Broadcast:  make(chan []byte),
	}
}

func (n *Network) Start(ca chan *Action) {
	for {
		select {
		case client := <-n.Register:
			n.Clients[client] = true
			go n.reader(client, ca)
			go n.writer(client)

		case client := <-n.Unregister:
			ca <- &Action{MessageType: 1, Client: client}
			client.Conn.Close()
			delete(n.Clients, client)

		case message := <-n.Broadcast:
			for client := range n.Clients {
				client.NetworkOutput <- message
			}
		}
	}
}

func (n *Network) reader(client *Client, ca chan *Action) {
	defer func() {
		n.Unregister <- client // close connection?
	}()

	for {
		_, message, err := client.Conn.ReadMessage()
		if err != nil {
			break
		}

		//fmt.Println("NETWORK:", message)
		if message[0] == 0 {
			ca <- &Action{MessageType: 0, Client: client}
		} else {
			client.NetworkInput <- message
		}
	}
}

func (n *Network) writer(client *Client) {
	defer func() {
		n.Unregister <- client
	}()

	for {
		select {
		case message, ok := <-client.NetworkOutput:
			if ok {
				client.Conn.WriteMessage(websocket.BinaryMessage, message)
			} else {
				// error?
			}
		}
	}
}
