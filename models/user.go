package models

import (
	"api.ethscrow/utils/database"
	"context"
	"errors"
	"time"
)

type User struct {
	Username     string    `json:"username,omitempty"`
	PublicKey    string    `json:"public_key,omitempty"`
	EncPublicKey string    `json:"enc_public_key,omitempty"`
	Email        *string   `json:"email,omitempty"`
	CreatedAt    time.Time `json:"created_at,omitempty"`
}

type Friend struct {
	Username  string `json:"username,omitempty"`
	PublicKey string `json:"public_key,omitempty"`
}

const createSQL = "INSERT INTO users (username, public_key, enc_public_key, email) VALUES($1, $2, $3, $4)"

func (u *User) Create() error {
	if u.Username == "" || u.PublicKey == "" {
		return errors.New("missing username")
	}

	_, err := database.DB.Exec(context.Background(), createSQL, u.Username, u.PublicKey, u.EncPublicKey, u.Email)
	return err
}

const existSQL = "SELECT username, public_key, enc_public_key, email FROM users WHERE username=$1"

func (u *User) Exists() (bool, error) {
	if u.Username == "" {
		return false, errors.New("missing username")
	}
	if err := database.DB.QueryRow(context.Background(), existSQL, u.Username).Scan(&u.Username, &u.PublicKey, &u.EncPublicKey, &u.Email); err != nil {
		return false, err
	}
	return true, nil
}

const registerSQL = "INSERT INTO users VALUES ($1, $2, $3)"

func (u *User) Register() error {
	if u.Username == "" || u.PublicKey == "" {
		return errors.New("missing username and/or public key")
	}
	_, err := database.DB.Exec(context.Background(), registerSQL, u.Username, u.PublicKey, u.Email)
	return err
}

const getFriends = "SELECT u.username, u.public_key FROM friendships f RIGHT JOIN users u on f.friend = u.username WHERE f.user=$1"

func (u *User) GetFriends() ([]Friend, error) {
	if u.Username == "" {
		return nil, errors.New("missing username")
	}
	rows, err := database.DB.Query(context.Background(), getFriends, u.Username)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var friends []Friend

	for rows.Next() {
		var friend Friend
		if err = rows.Scan(&friend.Username, &friend.PublicKey); err != nil {
			return nil, err
		}
		friends = append(friends, friend)
	}

	return friends, nil
}

const getPublicKey = "SELECT enc_public_key FROM users WHERE username=$1"

func (u *User) GetPublicKey() error {
	if err := database.DB.QueryRow(context.Background(), getPublicKey, u.Username).Scan(&u.EncPublicKey); err != nil {
		return err
	}
	return nil
}

const getAllPools = "SELECT * FROM pools WHERE bettor_username=$1 OR caller_username=$1 OR mediator_username=$1"

func (u *User) GetAllPools() ([]Pool, error) {

	rows, err := database.DB.Query(context.Background(), getAllPools, u.Username)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var pools []Pool

	for rows.Next() {
		var p Pool
		if err = rows.Scan(&p.ID, &p.Address, &p.Mediator, &p.Bettor, &p.Caller, &p.BetterState, &p.CallerState, &p.ThresholdKey, &p.CreatedAt, &p.Reason, &p.Balance, &p.BalanceLastUpdated, &p.Accepted, &p.Initialized, &p.ConflictTempData); err != nil {
			return nil, err
		}
		pools = append(pools, p)
	}

	return pools, nil
}
