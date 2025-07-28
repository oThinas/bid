package services

import (
	"context"
	"errors"
	"log/slog"
	"sync"
	"time"

	"github.com/google/uuid"
	"github.com/gorilla/websocket"
)

type MessageType int

const (
	// Requests
	PlaceBid MessageType = iota

	// Success
	SuccessfullyPlacedBid

	// Info
	NewBidPlaced
	AuctionEnded

	// Errors
	FailedToPlaceBid
	InvalidJSON
)

type Message struct {
	Message string      `json:"message,omitempty"`
	UserID  uuid.UUID   `json:"user_id,omitempty"`
	Amount  float64     `json:"amount,omitempty"`
	Type    MessageType `json:"type"`
}

type AuctionLobby struct {
	sync.Mutex
	Rooms map[uuid.UUID]*AuctionRoom
}

type AuctionRoom struct {
	ID          uuid.UUID
	Register    chan *Client
	Unregister  chan *Client
	Broadcast   chan Message
	Context     context.Context
	Clients     map[uuid.UUID]*Client
	BidsService BidsService
}

type Client struct {
	Conn   *websocket.Conn
	Send   chan Message
	Room   *AuctionRoom
	UserID uuid.UUID
}

func NewAuctionRoom(ctx context.Context, id uuid.UUID, bidsService BidsService) *AuctionRoom {
	return &AuctionRoom{
		ID:          id,
		Register:    make(chan *Client),
		Unregister:  make(chan *Client),
		Broadcast:   make(chan Message),
		Clients:     make(map[uuid.UUID]*Client),
		Context:     ctx,
		BidsService: bidsService,
	}
}

func NewClient(conn *websocket.Conn, room *AuctionRoom, userID uuid.UUID) *Client {
	return &Client{
		Conn:   conn,
		Send:   make(chan Message, 512),
		Room:   room,
		UserID: userID,
	}
}

func (r *AuctionRoom) Run() {
	slog.Info("Auction has begun", "AuctionID", r.ID)

	defer func() {
		close(r.Broadcast)
		close(r.Register)
		close(r.Unregister)
	}()

	for {
		select {
		case client := <-r.Register:
			r.registerClient(client)
		case client := <-r.Unregister:
			r.unregisterClient(client)
		case message := <-r.Broadcast:
			r.broadcastMessage(message)
		case <-r.Context.Done():
			slog.Info("Auction has ended", "AuctionID", r.ID)
			for _, client := range r.Clients {
				client.Send <- Message{
					Message: "Auction has ended",
					Type:    AuctionEnded,
				}
			}

			return
		}
	}
}

func (r *AuctionRoom) registerClient(client *Client) {
	slog.Info("New user connected", "Client:", client)
	r.Clients[client.UserID] = client
}

func (r *AuctionRoom) unregisterClient(client *Client) {
	slog.Info("User disconnected", "Client:", client)
	delete(r.Clients, client.UserID)
}

func (r *AuctionRoom) broadcastMessage(message Message) {
	slog.Info("New message received", "Room:", r.ID, "User:", message.UserID, "Message:", message.Message)
	switch message.Type {
	case PlaceBid:
		bid, err := r.BidsService.PlaceBid(r.Context, r.ID, message.UserID, message.Amount)
		if err != nil {
			if errors.Is(err, ErrBidAmountTooLow) {
				if client, ok := r.Clients[message.UserID]; ok {
					client.Send <- Message{
						Message: ErrBidAmountTooLow.Error(),
						Type:    FailedToPlaceBid,
						UserID:  message.UserID,
					}
				}

				return
			}
		}

		if client, ok := r.Clients[message.UserID]; ok {
			client.Send <- Message{
				Message: "Your bid was successfully placed",
				Type:    SuccessfullyPlacedBid,
				UserID:  message.UserID,
			}
		}

		for id, client := range r.Clients {
			newBidMessage := Message{
				Message: "A new bid was placed",
				Type:    NewBidPlaced,
				Amount:  bid.Amount,
				UserID:  message.UserID,
			}

			if id == message.UserID {
				continue
			}

			client.Send <- newBidMessage
		}

	case InvalidJSON:
		client, ok := r.Clients[message.UserID]
		if !ok {
			slog.Info("Client not found", "UserID", message.UserID)
			return
		}

		client.Send <- message
	}
}

func (c *Client) ReadEventLoop() {
	defer func() {
		c.Room.Unregister <- c
		c.Conn.Close()
	}()

	c.Conn.SetReadLimit(MaxMessageSize)
	c.Conn.SetReadDeadline(time.Now().Add(ReadDeadline))
	c.Conn.SetPongHandler(func(string) error {
		c.Conn.SetReadDeadline(time.Now().Add(ReadDeadline))
		return nil
	})

	for {
		var message Message
		message.UserID = c.UserID

		err := c.Conn.ReadJSON(&message)
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				slog.Error("Unexpected close error", "Error", err)
				return
			}

			c.Room.Broadcast <- Message{
				Message: "Invalid JSON",
				Type:    InvalidJSON,
				UserID:  message.UserID,
			}

			continue
		}

		c.Room.Broadcast <- message
	}
}

func (c *Client) WriteEventLoop() {
	ticker := time.NewTicker(PingInterval)
	defer func() {
		ticker.Stop()
		c.Conn.Close()
	}()

	for {
		select {
		case <-ticker.C:
			c.Conn.SetWriteDeadline(time.Now().Add(WriteDeadLine))
			if err := c.Conn.WriteMessage(websocket.PingMessage, nil); err != nil {
				slog.Error("Unexpected write error", "Error", err)
				return
			}

		case message, ok := <-c.Send:
			if !ok {
				c.Conn.WriteJSON(Message{
					Message: "Closing websocket connection",
					Type:    websocket.CloseMessage,
				})
				return
			}

			if message.Type == AuctionEnded {
				close(c.Send)
				return
			}

			c.Conn.SetWriteDeadline(time.Now().Add(WriteDeadLine))
			if err := c.Conn.WriteJSON(message); err != nil {
				c.Room.Unregister <- c
				return
			}
		}
	}
}
