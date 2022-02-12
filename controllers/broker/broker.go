package broker

import (
	"api.ethscrow/models"
	"api.ethscrow/utils"
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
	NeutralState  = 0
	ConflictState = -2
)

func CreatePool(w http.ResponseWriter, r *http.Request) {
	user := r.Context().Value("user").(models.User)
	pool := &models.Pool{}

	pool.Bettor = user.Username

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

	exists, _ := pool.Exists()
	poolComm, active := NewPool(pool.ID)
	if !active {
		poolComm.Pool = pool
	}

	if exists && (pool.Bettor != user.Username && pool.Caller != user.Username && pool.Mediator != user.Username) {
		utils.Error(w, http.StatusForbidden, "You are not part of the pool.")
		return
	} else if !exists {
		utils.Error(w, http.StatusForbidden, "Pool doesn't exist.")
		return
	}

	if !active {
		go poolComm.Start()
	}

	c, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		return
	}

	client := &Client{
		Username: user.Username,
		Conn:     c,
		PoolComm: poolComm,
	}

	poolComm.Register <- client

	client.Read()
}

func AcceptPool(w http.ResponseWriter, r *http.Request) {
	user := r.Context().Value("user").(models.User)
	pool := &models.Pool{
		ID: chi.URLParam(r, "PoolId"),
	}

	exists, _ := pool.Exists()
	if exists && pool.Caller != user.Username {
		utils.Error(w, http.StatusForbidden, "You are not part of the pool.")
		return
	} else if !exists {
		utils.Error(w, http.StatusForbidden, "Pool doesn't exist.")
		return
	}

	pool.Accepted = true

	if pool.Update() != nil {
		utils.JSON(w, http.StatusInternalServerError, "Failed to update pool state.")
		return
	}

	utils.JSON(w, http.StatusAccepted, pool)
}

func DeletePool(w http.ResponseWriter, r *http.Request) {
	user := r.Context().Value("user").(models.User)
	pool := &models.Pool{
		ID: chi.URLParam(r, "PoolId"),
	}

	exists, _ := pool.Exists()
	if exists && (pool.Bettor != user.Username || pool.Caller != user.Username) {
		utils.Error(w, http.StatusForbidden, "You are not part of the pool.")
		return
	} else if !exists {
		utils.Error(w, http.StatusForbidden, "Pool doesn't exist.")
		return
	}

	if pool.Accepted {
		utils.Error(w, http.StatusBadRequest, "Pool was already accepted.")
		return
	}

	if pool.Close() != nil {
		utils.Error(w, http.StatusInternalServerError, "Error closing pool.")
		return
	}

	utils.JSON(w, http.StatusAccepted, &utils.BasicData{Data: "Pool closed."})
}

func UpdatePoolState(w http.ResponseWriter, r *http.Request) {
	user := r.Context().Value("user").(models.User)

	pool := &models.Pool{
		ID: chi.URLParam(r, "PoolId"),
	}

	exists, _ := pool.Exists()
	if exists && (pool.Bettor != user.Username || pool.Caller != user.Username || pool.Mediator != user.Username) {
		utils.Error(w, http.StatusForbidden, "You are not part of the pool.")
		return
	} else if !exists {
		utils.Error(w, http.StatusForbidden, "Pool doesn't exist.")
		return
	}

	// If pool is in conflict
	if pool.BetterState != ConflictState || pool.CallerState != ConflictState {
		utils.JSON(w, http.StatusConflict, pool)
		return
	}

	// if pool's winner is already set
	if pool.ThresholdKey != nil {
		utils.JSON(w, http.StatusAlreadyReported, pool)
		return
	}

	state := &stateChangeRequest{}

	if err := utils.ParseRequestBody(r, &state); err != nil || state.NewState == 0 {
		utils.Error(w, http.StatusBadRequest, "Invalid request.")
		return
	}

	otherUserPoolState := &pool.BetterState
	currentUserPoolState := &pool.CallerState

	if user.Username == pool.Bettor {
		otherUserPoolState = &pool.CallerState
		currentUserPoolState = &pool.BetterState
	}

	if state.NewState == LostState && state.ThresholdKey != nil { // user lost
		pool.ThresholdKey = state.ThresholdKey
		*currentUserPoolState = LostState
		*otherUserPoolState = WonState
	} else if state.NewState == WonState && *otherUserPoolState == NeutralState { // user won and other user undecided
		*currentUserPoolState = WonState
	} else if state.NewState == WonState && *otherUserPoolState == WonState || state.NewState == ConflictState { // user won and other user also won
		pool.CallerState = ConflictState
		pool.BetterState = ConflictState
		pool.ThresholdKey = state.ThresholdKey
	} else {
		utils.Error(w, http.StatusBadRequest, "Invalid state change.")
		return
	}

	if pool.Update() != nil {
		utils.JSON(w, http.StatusInternalServerError, "Failed to update pool state.")
		return
	}

	utils.JSON(w, http.StatusAccepted, pool)
}

func ResolveConflict(w http.ResponseWriter, r *http.Request) {
	user := r.Context().Value("user").(models.User)
	pool := &models.Pool{
		ID: chi.URLParam(r, "PoolId"),
	}

	resolution := &resolveConflictRequest{}

	if err := utils.ParseRequestBody(r, &resolution); err != nil || resolution.WinnerUsername == "" || resolution.ThresholdKey != "" {
		utils.Error(w, http.StatusBadRequest, "Invalid request.")
		return
	}

	exists, _ := pool.Exists()
	if exists && pool.Mediator != user.Username {
		utils.Error(w, http.StatusForbidden, "You are not part of the pool.")
		return
	} else if !exists {
		utils.Error(w, http.StatusForbidden, "Pool doesn't exist.")
		return
	}

	// pool is not in conflict
	if pool.BetterState != ConflictState || pool.CallerState != ConflictState {
		utils.Error(w, http.StatusForbidden, "Not in conflict.")
		return
	}

	// invalid winner name set
	if resolution.WinnerUsername != pool.Bettor && resolution.WinnerUsername != pool.Caller {
		utils.Error(w, http.StatusForbidden, "Invalid winner.")
		return
	}

	pool.BetterState = WonState
	pool.CallerState = LostState

	if resolution.WinnerUsername == pool.Caller {
		pool.BetterState = LostState
		pool.CallerState = WonState
	}

	*pool.ThresholdKey = resolution.ThresholdKey

	if pool.Update() != nil {
		utils.JSON(w, http.StatusInternalServerError, "Failed to update pool state.")
		return
	}

	utils.JSON(w, http.StatusAccepted, pool)
}
