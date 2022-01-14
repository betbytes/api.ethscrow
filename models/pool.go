package models

import (
	"api.ethscrow/utils/database"
	"context"
)

type Pool struct {
	ID     string
	bettor string
	caller string
	Reason string
}

const poolExists = "SELECT bettor, caller, reason WHERE id=$1"

func (p *Pool) Exists() (bool, error) {
	if err := database.DB.QueryRow(context.Background(), poolExists, p.ID).Scan(&p.bettor, &p.caller, &p.Reason); err != nil {
		return false, err
	}
	return true, nil
}

const createPool = "INSERT INTO pools VALUES ($1, $2, $3, $4)"

func (p *Pool) Create() error {
	_, err := database.DB.Exec(context.Background(), createPool, p.ID, p.bettor, p.caller, p.Reason)
	return err
}

const closePool = "DELETE FROM pools WHERE id=$1"

func (p *Pool) Close() error {
	_, err := database.DB.Exec(context.Background(), closePool, p.ID)
	return err
}
