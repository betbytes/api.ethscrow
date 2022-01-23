package broker

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/websocket"
	"log"
)

type Client struct {
	ID       string
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

// ConnectBody body of connect message type
type ConnectBody struct {
	Type MessageEnum `json:"type"`
	User string      `json:"user"`
}

// DisconnectBody body of connect message type
type DisconnectBody struct {
	Type MessageEnum `json:"type"`
	User string      `json:"user"`
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
	From *Client         `json:"-"`
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
		msg := &Message{
			From: c,
		}

		err := c.Conn.ReadJSON(msg)
		if err != nil {
			log.Println(err)
			return
		}

		switch msg.Type {
		case Chat:
			chat := ChatBody{}
			if json.Unmarshal(msg.Body, &chat) == nil {
				if chat.Important {
					_ = c.PoolComm.Log(chat.Message)
				}
			}
			break
		case RefreshBalance:
			msg.Body = []byte{}
			break
		}

		c.PoolComm.Broadcast <- msg
		log.Printf("Forwarded: %+v\n", msg)
	}
}
