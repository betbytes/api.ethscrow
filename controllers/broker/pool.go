package broker

import "log"

var ActivePools = make(map[string]*PoolComm)

type PoolComm struct {
	Register   chan *Client
	Unregister chan *Client
	Clients    map[*Client]bool
	Broadcast  chan *Message
}

func NewPool(id string) *PoolComm {
	if _, ok := ActivePools[id]; ok {
		return ActivePools[id]
	}

	ActivePools[id] = &PoolComm{
		Register:   make(chan *Client),
		Unregister: make(chan *Client),
		Clients:    make(map[*Client]bool),
		Broadcast:  make(chan *Message),
	}

	ActivePools[id].Start()

	return ActivePools[id]
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

			for client, _ = range p.Clients {
				err = client.Conn.WriteJSON(&DisconnectBody{
					Type: Disconnect,
					User: client.ID,
				})
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
}

func (p *PoolComm) Log() {

}
