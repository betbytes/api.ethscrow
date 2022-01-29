package models

import (
	"api.ethscrow/utils/database"
	"context"
	"errors"
	"github.com/google/uuid"
	"strings"
	"time"
)

type Pool struct {
	ID                 string     `json:"id,omitempty"`
	Address            *string    `json:"address,omitempty"`
	Bettor             string     `json:"bettor_username,omitempty"`
	BetterState        int16      `json:"betterState,omitempty"`
	Caller             string     `json:"caller_username,omitempty"`
	CallerState        int16      `json:"callerState,omitempty"`
	Mediator           *string    `json:"mediator_username,omitempty"`
	ThresholdKey       *string    `json:"threshold_key,omitempty"`
	Reason             string     `json:"reason,omitempty"`
	Chats              []string   `json:"chats,omitempty"`
	CreatedAt          time.Time  `json:"created_at,omitempty"`
	Balance            float64    `json:"balance,omitempty"`
	BalanceLastUpdated *time.Time `json:"balance_last_updated,omitempty"`
	Accepted           bool       `json:"accepted,omitempty"`
}

const poolExists = "SELECT * FROM pools WHERE id=$1"
const getChats = "SELECT message FROM chats WHERE pool_id=$1 ORDER BY timestamp"

func (p *Pool) Exists() (bool, error) {
	if p.ID == "" {
		return false, errors.New("missing pool id")
	}
	if err := database.DB.QueryRow(context.Background(), poolExists, p.ID).
		Scan(&p.ID, &p.Address, &p.Mediator, &p.Bettor, &p.Caller, &p.BetterState, &p.CallerState, &p.ThresholdKey, &p.CreatedAt, &p.Reason, &p.Balance, &p.BalanceLastUpdated, &p.Accepted); err != nil {
		return false, err
	}

	rows, err := database.DB.Query(context.Background(), getChats, p.ID)
	if err != nil {
		return true, err
	}
	defer rows.Close()

	for rows.Next() {
		var msg string
		if err = rows.Scan(&msg); err != nil {
			return true, err
		}
		p.Chats = append(p.Chats, msg)
	}

	return true, nil
}

const createPool = "INSERT INTO pools (id, bettor_username, caller_username, mediator_username, reason) VALUES ($1, $2, $3, $4, $5) RETURNING created_at"

func (p *Pool) Create() error {
	p.ID = strings.ReplaceAll(uuid.New().String(), "-", "")
	if p.Bettor == "" || p.Caller == "" {
		return errors.New("missing parameter")
	}

	err := database.DB.QueryRow(context.Background(), createPool, p.ID, p.Bettor, p.Caller, p.Mediator, p.Reason).Scan(&p.CreatedAt)
	return err
}

const closePool = "DELETE FROM pools WHERE id=$1"
const clearChats = "DELETE FROM chats WHERE pool_id=$1"

func (p *Pool) Close() error {
	if p.ID == "" {
		return errors.New("missing pool id")
	}
	_, err := database.DB.Exec(context.Background(), clearChats, p.ID)
	_, err = database.DB.Exec(context.Background(), closePool, p.ID)
	return err
}

const updatePool = "UPDATE pools SET bettor_state=$2, caller_state=$3, threshold_key=$4, accepted=$5 WHERE id=$1"

func (p *Pool) Update() error {
	_, err := database.DB.Exec(context.Background(), updatePool, p.ID, p.BetterState, p.CallerState, p.ThresholdKey, p.Accepted)
	return err
}
