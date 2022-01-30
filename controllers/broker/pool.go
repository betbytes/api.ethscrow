package broker

import (
	"api.ethscrow/models"
	"api.ethscrow/utils/database"
	"context"
	"encoding/json"
	"errors"
	"log"
	"sync"
)

var ActivePools = make(map[string]*PoolComm)

type UserComm struct {
	Client *Client
	Mutex  *sync.Mutex
}

type PoolComm struct {
	Register    chan *Client
	Unregister  chan *Client
	ActiveUsers map[string]UserComm
	Broadcast   chan *Message
	Pool        *models.Pool
}

func NewPool(id string) (*PoolComm, bool) {
	if _, ok := ActivePools[id]; ok {
		return ActivePools[id], true
	}

	ActivePools[id] = &PoolComm{
		Register:    make(chan *Client),
		Unregister:  make(chan *Client),
		ActiveUsers: make(map[string]UserComm),
		Broadcast:   make(chan *Message),
	}

	go ActivePools[id].Start()

	return ActivePools[id], false
}

func (p *PoolComm) Start() {
	for {
		var err error
		select {
		case client := <-p.Register:
			p.ActiveUsers[client.Username] = UserComm{client, &sync.Mutex{}}

			for username, userComm := range p.ActiveUsers {
				userComm.Mutex.Lock()
				if username == client.Username {
					otherUsername := p.Pool.Bettor
					if p.Pool.Bettor == username {
						otherUsername = p.Pool.Caller
					}

					_, otherActive := p.ActiveUsers[otherUsername]

					init := &InitialBody{
						Type:               PoolDetails,
						Pool:               p.Pool,
						OtherUserConnected: otherActive,
					}

					initJson, _ := json.Marshal(init)
					err = userComm.Client.Conn.WriteJSON(&Message{
						Type: PoolDetails,
						Body: initJson,
					})
				} else {
					err = userComm.Client.Conn.WriteJSON(&ConnectBody{
						Type:     Connect,
						Username: client.Username,
					})
				}
				userComm.Mutex.Unlock()
			}
			break
		case client := <-p.Unregister:
			delete(p.ActiveUsers, client.Username)

			empty := true
			for _, userComm := range p.ActiveUsers {
				userComm.Mutex.Lock()
				err = userComm.Client.Conn.WriteJSON(&DisconnectBody{
					Type:     Disconnect,
					Username: client.Username,
				})
				userComm.Mutex.Unlock()
				empty = false
			}

			if empty {
				goto Close // Although not ideal, it is reasonable to use goto to exit loop
			}

			break
		case message := <-p.Broadcast:
			for username, userComm := range p.ActiveUsers {
				if message.From != nil && username == *message.From {
					continue
				}
				userComm.Mutex.Lock()
				err = userComm.Client.Conn.WriteJSON(message)
				userComm.Mutex.Unlock()
			}
		}

		if err != nil {
			log.Println(err)
		}
	}
Close:
	defer delete(ActivePools, p.Pool.ID)
}

const logMessage = "INSERT INTO chats(id, pool_id, message, from_username) VALUES($1, $2, $3, $4)"

func (p *PoolComm) Log(chat *models.Chat) error {
	if chat.Message == "" {
		return errors.New("missing message")
	}

	_, err := database.DB.Exec(context.Background(), logMessage, chat.ID, p.Pool.ID, chat.Message, chat.From)
	return err
}
