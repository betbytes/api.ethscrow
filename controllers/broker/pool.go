package broker

import "log"

var Pools = make(map[string]*Pool)

type Pool struct {
	Register   chan *Client
	Unregister chan *Client
	Clients    map[*Client]bool
	Broadcast  chan *Message
}

func NewPool(id string) *Pool {
	if _, ok := Pools[id]; ok {
		return Pools[id]
	}

	Pools[id] = &Pool{
		Register:   make(chan *Client),
		Unregister: make(chan *Client),
		Clients:    make(map[*Client]bool),
		Broadcast:  make(chan *Message),
	}

	Pools[id].Start()

	return Pools[id]
}

func (p *Pool) Start() {
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

func (p *Pool) Log() {

}
