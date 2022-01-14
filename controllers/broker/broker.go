package broker

import (
	"api.ethscrow/models"
	"api.ethscrow/utils"
	"github.com/go-chi/chi/v5"
	"github.com/gorilla/websocket"
	"net/http"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool { return true },
}

func ConnectToPool(w http.ResponseWriter, r *http.Request) {
	user := r.Context().Value("user").(models.User)
	pool := &models.Pool{
		ID: chi.URLParam(r, "PoolId"),
	}

	var body map[string]string
	if err := utils.ParseRequestBody(r, &body); err != nil {
		utils.Error(w, http.StatusInternalServerError, "Invalid")
		return
	}

	var exists bool
	poolComm, active := NewPool(pool.ID)
	if active {
		pool = poolComm.Pool
		exists = true
	} else {
		exists, _ = pool.Exists()
	}

	if exists && (pool.Bettor != user.Username || pool.Caller != user.Username || pool.Mediator != user.Username) {
		utils.Error(w, http.StatusForbidden, "You are not part of the pool.")
		return
	} else if !exists && body["caller"] != "" && body["mediator"] != "" {
		pool.Bettor = user.Username
		pool.Caller = body["caller"]
		pool.Mediator = body["mediator"]

		if pool.Create() != nil {
			utils.Error(w, http.StatusInternalServerError, "Error creating pool.")
			return
		}
	} else if !exists && (body["caller"] == "" || body["mediator"] == "") {
		utils.Error(w, http.StatusBadRequest, "Missing caller and mediator details.")
		return
	}

	c, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		return
	}

	client := &Client{
		Conn:     c,
		PoolComm: poolComm,
	}

	poolComm.Register <- client

	client.Read()
}
