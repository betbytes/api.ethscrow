package models

import (
	"api.ethscrow/utils/database"
	"context"
)

const existSQL = "SELECT exists(SELECT 1 FROM users where username=$1)"
const registerSQL = "INSERT INTO users VALUES ($1, $2, $3)"
const loginSQL = "SELECT email, password, public_key FROM users WHERE username=$1"

type User struct {
	Username  string `json:"username,omitempty"`
	Email     string `json:"email,omitempty"`
	PublicKey string `json:"publicKey,omitempty"`
	Token     string `json:"token,omitempty"`
}

type Friend struct {
	Username  string `json:"username,omitempty"`
	PublicKey string `json:"publicKey,omitempty"`
}

func (u *User) Exists() (bool, error) {
	var exists bool
	if err := database.DB.QueryRow(context.Background(), existSQL, u.Username).Scan(&exists); err != nil {
		return true, err
	}
	return exists, nil
}

func (u *User) Register() error {
	_, err := database.DB.Exec(context.Background(), registerSQL, u.Username, u.PublicKey)
	return err
}

func (u *User) Login() error {
	return nil
}

func (u *User) GetFriends() ([]Friend, error) {
	return nil, nil
}
