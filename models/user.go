package models

import (
	"api.ethscrow/utils/database"
	"context"
)

type User struct {
	Username  string `json:"username,omitempty"`
	PublicKey []byte `json:"publicKey,omitempty"`
}

type Friend struct {
	Username  string `json:"username,omitempty"`
	PublicKey string `json:"publicKey,omitempty"`
}

const existSQL = "SELECT username, public_key FROM users WHERE username=$1)"

func (u *User) Exists() (bool, error) {
	if err := database.DB.QueryRow(context.Background(), existSQL, u.Username).Scan(&u.Username, &u.PublicKey); err != nil {
		return false, err
	}
	return true, nil
}

const registerSQL = "INSERT INTO users VALUES ($1, $2, $3)"

func (u *User) Register() error {
	_, err := database.DB.Exec(context.Background(), registerSQL, u.Username, u.PublicKey)
	return err
}

const getFriends = "SELECT u.username, u.public_key FROM friendships f RIGHT JOIN users u on f.friend = u.username WHERE f.user=$1"

func (u *User) GetFriends() ([]Friend, error) {
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
