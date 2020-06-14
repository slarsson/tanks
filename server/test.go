package main

// import (
// 	"log"
// 	"math/rand"
// 	"net/http"

// 	"github.com/gorilla/websocket"
// )

// var upgrader = websocket.Upgrader{}

// func server(s *Server, w http.ResponseWriter, r *http.Request) {
// 	c, err := upgrader.Upgrade(w, r, nil)
// 	if err != nil {
// 		log.Print("UPGRADE ERROR:", err)
// 		return
// 	}

// 	client := &Client{
// 		conn:          c,
// 		position:      &Point{X: 0, Y: 0, Z: 0},
// 		NetworkInput:  make(chan []byte, 100),
// 		NetworkOutput: make(chan []byte, 100),
// 		ID:            rand.Intn(10000),
// 	}
// 	s.register <- client // register new player
// }

// func main() {
// 	upgrader.CheckOrigin = func(r *http.Request) bool { return true }

// 	var x = newServer()
// 	go x.run()

// 	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
// 		server(x, w, r)
// 	})

// 	log.Fatal(http.ListenAndServe("localhost:1337", nil))
// }
