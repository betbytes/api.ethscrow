package broker

import (
	"github.com/go-chi/chi/v5"
	"github.com/gorilla/websocket"
	"net/http"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool { return true },
}

func ConnectToPool(w http.ResponseWriter, r *http.Request) {
	poolId := chi.URLParam(r, "PoolId")

	// TODO: Check if room is in database
	// TODO: if it exists, check if user is part of the pool
	// TODO: if new, create a new pool, with current user as one of participants

	c, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		return
	}

	pool := NewPool(poolId)
	client := &Client{
		Conn: c,
		Pool: pool,
	}

	pool.Register <- client

	client.Read()
}
