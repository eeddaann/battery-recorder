package main

import (
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"time"

	socketio "github.com/googollee/go-socket.io"
)

func sendRandom() float64 {
	return (rand.Float64() * 5) + 5
}

func main() {
	server := socketio.NewServer(nil)

	server.OnConnect("/", func(s socketio.Conn) error {
		s.SetContext("")
		fmt.Println("connected:", s.ID())
		return nil
	})

	server.OnError("/", func(s socketio.Conn, e error) {
		fmt.Println("meet error:", e)
	})

	server.OnDisconnect("/", func(s socketio.Conn, reason string) {
		fmt.Println("closed", reason)
	})

	go server.Serve()

	go func() {
		for {
			time.Sleep(time.Second / 5)
			server.BroadcastToNamespace("/", "temp", fmt.Sprintf("%v", sendRandom()))
		}
	}()

	defer server.Close()

	http.Handle("/socket.io/", server)
	http.Handle("/", http.FileServer(http.Dir("./asset")))
	log.Println("Serving at localhost:8000...")
	log.Fatal(http.ListenAndServe(":8000", nil))
}
