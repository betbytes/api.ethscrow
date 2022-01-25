package broker

import (
	"api.ethscrow/models"
	"api.ethscrow/utils"
	"encoding/json"
	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"github.com/gorilla/websocket"
	"net/http"
	"strings"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool { return true },
}

const (
	LostState     = -1
	WonState      = 1
	ConflictState = 0
)

func CreatePool(w http.ResponseWriter, r *http.Request) {
	pool := &models.Pool{}
	pool.Bettor = "ahmad"

	if err := utils.ParseRequestBody(r, &pool); err != nil {
		utils.Error(w, http.StatusInternalServerError, "Invalid")
		return
	}

	pool.ID = strings.ReplaceAll(uuid.New().String(), "-", "")

	if pool.Create() != nil {
		utils.Error(w, http.StatusInternalServerError, "Error creating pool.")
		return
	}

	utils.JSON(w, http.StatusCreated, pool)
}

// ConnectToPool /connect/{roomId} - Authenticated
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

	if exists && (pool.Bettor != user.Username || pool.Caller != user.Username || *pool.Mediator != user.Username) {
		utils.Error(w, http.StatusForbidden, "You are not part of the pool.")
		return
	} else if !exists {
		utils.Error(w, http.StatusForbidden, "Pool doesn't exist.")
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

	poolJson, err := json.Marshal(pool)

	poolComm.Broadcast <- &Message{
		From: client,
		Type: PoolDetails,
		Body: poolJson,
	}

	client.Read()
}

func DeletePool(w http.ResponseWriter, r *http.Request) {
	user := r.Context().Value("claims").(models.User)
	pool := &models.Pool{
		ID: chi.URLParam(r, "PoolId"),
	}

	exists, _ := pool.Exists()
	if exists && (pool.Bettor != user.Username || pool.Caller != user.Username || *pool.Mediator != user.Username) {
		utils.Error(w, http.StatusForbidden, "You are not part of the pool.")
		return
	} else if !exists {
		utils.Error(w, http.StatusForbidden, "Pool doesn't exist.")
		return
	}

	if pool.Close() != nil {
		utils.Error(w, http.StatusInternalServerError, "Error closing pool.")
		return
	}

	utils.JSON(w, http.StatusAccepted, &utils.BasicData{Data: "Pool closed."})
}

func UpdatePoolState(w http.ResponseWriter, r *http.Request) {
	user := r.Context().Value("claims").(models.User)

	pool := &models.Pool{
		ID: chi.URLParam(r, "PoolId"),
	}

	exists, _ := pool.Exists()
	if exists && (pool.Bettor != user.Username || pool.Caller != user.Username || *pool.Mediator != user.Username) {
		utils.Error(w, http.StatusForbidden, "You are not part of the pool.")
		return
	} else if !exists {
		utils.Error(w, http.StatusForbidden, "Pool doesn't exist.")
		return
	}

	state := &stateChangeRequest{}

	if err := utils.ParseRequestBody(r, &state); err != nil || state.NewState == 0 {
		utils.Error(w, http.StatusBadRequest, "Invalid request.")
		return
	}

	if state.NewState == LostState && state.ThresholdKey != nil {
		pool.ThresholdKey = state.ThresholdKey
		// TODO: Continue from here
	}
}
