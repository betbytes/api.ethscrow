package broker

import (
	"api.ethscrow/models"
	"encoding/json"
	"fmt"
	"github.com/google/uuid"
	"github.com/gorilla/websocket"
	"log"
	"strings"
	"time"
)

type Client struct {
	Username string
	Conn     *websocket.Conn
	PoolComm *PoolComm
}

type MessageEnum int

const (
	Connect = iota
	Disconnect
	Chat
	PoolDetails
	RefreshBalance
	GeneratingEscrow
	Offer
	Answer
	OfferCandidate
	AnswerCandidate
	InitializePool
	PoolStateChange
)

type InitialBody struct {
	Type               MessageEnum  `json:"type"`
	Pool               *models.Pool `json:"pool"`
	OtherUserConnected bool         `json:"other_user_connected"`
	MediatorPublicKey  string       `json:"mediator_public_key"`
}

// ConnectBody body of connect message type
type ConnectBody struct {
	Type     MessageEnum `json:"type"`
	Username string      `json:"username"`
}

// DisconnectBody body of connect message type
type DisconnectBody struct {
	Type     MessageEnum `json:"type"`
	Username string      `json:"username"`
}

// ChatBody body of connect message type
type ChatBody struct {
	Timestamp string `json:"timestamp"`
	Message   string `json:"message"`
	Important bool   `json:"important"`
}

// WebRTCInit body of connect message type
type WebRTCInit struct {
	Action        string `json:"action,omitempty"`
	SDP           string `json:"sdp,omitempty"`
	Candidate     string `json:"candidate,omitempty"`
	SDPMid        string `json:"sdpMid,omitempty"`
	SDPMLineIndex int    `json:"sdpMLineIndex,omitempty"`
}

type Message struct {
	Type MessageEnum     `json:"type"`
	Body json.RawMessage `json:"body"`
	From *string         `json:"-"`
}

type InitializedPool struct {
	Address string `json:"address"`
}

type BalanceResponse struct {
	Balance   int64     `json:"balance"`
	UpdatedAt time.Time `json:"updated_at"`
}

func (m *Message) String() string {
	return fmt.Sprintf("Type: %d Message: %s", m.Type, string(m.Body))
}

func (c *Client) Read() {
	defer func() {
		c.PoolComm.Unregister <- c
		c.Conn.Close()
	}()

	for {
		msg := &Message{}

		err := c.Conn.ReadJSON(msg)
		if err != nil {
			log.Println(err)
			return
		}

		switch msg.Type {
		case Chat:
			chat := &models.Chat{}
			chat.From = c.Username
			chat.ID = strings.ReplaceAll(uuid.New().String(), "-", "")
			if json.Unmarshal(msg.Body, &chat) == nil {
				c.PoolComm.Pool.Chats = append(c.PoolComm.Pool.Chats, *chat)
				if chat.Important {
					_ = c.PoolComm.Log(chat)
				}
			}

			chatMarshal, _ := json.Marshal(chat)
			msg.Body = chatMarshal
			break
		case RefreshBalance:
			if c.PoolComm.Pool.UpdateBalance() != nil {
				log.Println("Couldn't update wallet balance.")
			}

			balance := &BalanceResponse{
				Balance:   c.PoolComm.Pool.Balance,
				UpdatedAt: *c.PoolComm.Pool.BalanceLastUpdated,
			}
			msg.Body, err = json.Marshal(balance)
			if err != nil {
				log.Println("Failed to marshal balance.")
				return
			}
			break
		case Offer, Answer, OfferCandidate, AnswerCandidate, GeneratingEscrow:
			msg.From = &c.Username
			break
		case InitializePool:
			initiated := &InitializedPool{}
			if json.Unmarshal(msg.Body, &initiated) == nil {
				c.PoolComm.Pool.Address = &initiated.Address
				c.PoolComm.Pool.Initialized = true
				if c.PoolComm.Pool.Update() != nil {
					log.Println("Error updating pool")
					return
				}
			}
			break
		}

		c.PoolComm.Broadcast <- msg
		log.Printf("Forwarded: %+v\n", msg)
	}
}
