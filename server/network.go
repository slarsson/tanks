package main

import (
	"github.com/gorilla/websocket"
)

type Network struct {
	clients    map[*Client]bool
	register   chan *Client
	unregister chan *Client
	broadcast  chan []byte
}

type Client struct {
	conn          *websocket.Conn
	NetworkInput  chan []byte
	NetworkOutput chan []byte
}

func NewNetwork() *Network {
	return &Network{
		clients:    make(map[*Client]bool),
		register:   make(chan *Client),
		unregister: make(chan *Client),
		broadcast:  make(chan []byte),
	}
}

func (n *Network) Start(ca chan *Action) {
	for {
		select {
		case client := <-n.register:
			n.clients[client] = true
			go n.reader(client, ca)
			go n.writer(client)

		case client := <-n.unregister:
			ca <- &Action{MessageType: 1, Client: client}
			client.conn.Close()
			delete(n.clients, client)

		case message := <-n.broadcast:
			for client := range n.clients {
				client.NetworkOutput <- message
			}
		}
	}
}

func (n *Network) reader(client *Client, ca chan *Action) {
	defer func() {
		n.unregister <- client // close connection?
	}()

	for {
		_, message, err := client.conn.ReadMessage()
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
		n.unregister <- client
	}()

	for {
		select {
		case message, ok := <-client.NetworkOutput:
			if ok {
				client.conn.WriteMessage(websocket.BinaryMessage, message)
			} else {
				// error?
			}
		}
	}
}
