package broker

import (
	"api.ethscrow/models"
	"api.ethscrow/utils/database"
	"context"
	"errors"
	"github.com/google/uuid"
	"log"
)

var ActivePools = make(map[string]*PoolComm)

type PoolComm struct {
	Register   chan *Client
	Unregister chan *Client
	Clients    map[*Client]bool
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
		Clients:    make(map[*Client]bool),
		Broadcast:  make(chan *Message),
	}

	ActivePools[id].Start()

	return ActivePools[id], false
}

func (p *PoolComm) Start() {
	for {
		var err error
		select {
		case client := <-p.Register:
			p.Clients[client] = true
			for client, _ = range p.Clients {
				err = client.Conn.WriteJSON(&ConnectBody{
					Type: Connect,
					User: client.ID,
				})
			}
			break
		case client := <-p.Unregister:
			delete(p.Clients, client)

			empty := true
			for client, _ = range p.Clients {
				err = client.Conn.WriteJSON(&DisconnectBody{
					Type: Disconnect,
					User: client.ID,
				})
				empty = false
			}

			if empty {
				goto Close // Although not ideal, it is reasonable to use goto to exit loop
			}

			break
		case message := <-p.Broadcast:
			for client, _ := range p.Clients {
				if client != message.From {
					err = client.Conn.WriteJSON(message)
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

const logMessage = "INSERT INTO chats(id, pool_id, message) VALUES($1, $2, $3)"

func (p *PoolComm) Log(msg string) error {
	if msg == "" {
		return errors.New("missing message")
	}

	id := uuid.New().String()

	_, err := database.DB.Exec(context.Background(), logMessage, id, p.Pool.ID, msg)
	return err
}
