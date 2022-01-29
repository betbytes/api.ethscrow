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

type PoolComm struct {
	Register   chan *Client
	Unregister chan *Client
	Clients    map[*Client]*sync.Mutex
	Broadcast  chan *Message
	Pool       *models.Pool
}

func NewPool(id string) (*PoolComm, bool) {
	if _, ok := ActivePools[id]; ok {
		return ActivePools[id], true
	}

	ActivePools[id] = &PoolComm{
		Register:   make(chan *Client),
		Unregister: make(chan *Client),
		Clients:    make(map[*Client]*sync.Mutex),
		Broadcast:  make(chan *Message),
	}

	go ActivePools[id].Start()

	return ActivePools[id], false
}

func (p *PoolComm) Start() {
	for {
		var err error
		select {
		case client := <-p.Register:
			p.Clients[client] = &sync.Mutex{}
			for c, _ := range p.Clients {
				p.Clients[c].Lock()
				if c == client {
					poolJson, _ := json.Marshal(p.Pool)
					err = c.Conn.WriteJSON(&Message{
						Type: PoolDetails,
						Body: poolJson,
					})
				} else {
					err = c.Conn.WriteJSON(&ConnectBody{
						Type:     Connect,
						Username: client.Username,
					})
				}
				p.Clients[c].Unlock()
			}
			break
		case client := <-p.Unregister:
			delete(p.Clients, client)

			empty := true
			for client, _ = range p.Clients {
				p.Clients[client].Lock()
				err = client.Conn.WriteJSON(&DisconnectBody{
					Type:     Disconnect,
					Username: client.Username,
				})
				p.Clients[client].Unlock()
				empty = false
			}

			if empty {
				goto Close // Although not ideal, it is reasonable to use goto to exit loop
			}

			break
		case message := <-p.Broadcast:
			for client, _ := range p.Clients {
				if client != message.From {
					p.Clients[client].Lock()
					err = client.Conn.WriteJSON(message)
					p.Clients[client].Unlock()
				}
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
