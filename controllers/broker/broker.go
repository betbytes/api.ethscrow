package broker

import (
	"github.com/go-chi/chi/v5"
	"github.com/gorilla/websocket"
	"net/http"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool { return true },
}

func Broker(w http.ResponseWriter, r *http.Request) {
	roomId := chi.URLParam(r, "roomId")
	c, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		return
	}

	pool := NewPool(roomId)
	client := &Client{
		Conn: c,
		Pool: pool,
	}

	pool.Register <- client

	client.Read()
}
