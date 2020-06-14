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

func newNetwork() *Network {
	return &Network{
		clients:    make(map[*Client]bool),
		register:   make(chan *Client),
		unregister: make(chan *Client),
		broadcast:  make(chan []byte),
	}
}

func (n *Network) start() {
	for {
		select {
		case client := <-n.register:
			n.clients[client] = true
			go n.reader(client)
			go n.writer(client)
		case client := <-n.unregister:
			delete(n.clients, client)
		case message := <-n.broadcast:
			for client := range n.clients {
				client.NetworkOutput <- message
			}
		}
	}
}

func (n *Network) reader(client *Client) {
	defer func() {
		n.unregister <- client // close connection?
	}()

	for {
		_, message, err := client.conn.ReadMessage()
		if err != nil {
			break
		}
		client.NetworkInput <- message
		//n.broadcast <- message
	}
}

func (n *Network) writer(client *Client) {
	// TODO: add defer func

	for {
		select {
		case message, ok := <-client.NetworkOutput:
			if ok {
				client.conn.WriteMessage(websocket.TextMessage, message)
			}
		}
	}
}
