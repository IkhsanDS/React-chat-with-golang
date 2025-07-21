package main

import (
	"log"
	"net/http"

	"github.com/IkhsanSeto/go-chat-app/hub"
)

func serveWs(h *hub.Hub, w http.ResponseWriter, r *http.Request) {
	conn, err := hub.Upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
		return
	}
	client := &hub.Client{Conn: conn, Send: make(chan []byte, 256)}
	h.Register <- client

	go client.WritePump()
	go client.ReadPump(h)
}

func main() {
	h := hub.NewHub()
	go h.Run()

	http.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
		serveWs(h, w, r)
	})

	log.Println("Server started on :8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
