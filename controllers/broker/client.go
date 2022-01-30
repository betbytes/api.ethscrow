package broker

import (
	"api.ethscrow/models"
	"encoding/json"
	"fmt"
	"github.com/google/uuid"
	"github.com/gorilla/websocket"
	"log"
	"strings"
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
	Offer
	Answer
	OfferCandidate
	AnswerCandidate
)

type InitialBody struct {
	Type               MessageEnum  `json:"type"`
	Pool               *models.Pool `json:"pool"`
	OtherUserConnected bool         `json:"other_user_connected"`
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
				if chat.Important {
					_ = c.PoolComm.Log(chat)
				}
			}

			chatMarshal, _ := json.Marshal(chat)
			msg.Body = chatMarshal
			break
		case RefreshBalance:
			msg.Body = []byte{}
			break
		}

		c.PoolComm.Broadcast <- msg
		log.Printf("Forwarded: %+v\n", msg)
	}
}
