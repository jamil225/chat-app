package mysocket

import "time"

type MessageStatus string

const (
	StatusSent      MessageStatus = "SENT"
	StatusDelivered MessageStatus = "DELIVERED"
	StatusRead      MessageStatus = "READ"
)

type Message struct {
	MessageID   string        `json:"message_id"`   // Unique ID of the message
	SenderID    int64         `json:"sender_id"`    // Sender's ID
	RecipientID int64         `json:"recipient_id"` // Recipient's ID
	Content     string        `json:"content"`      // Message content
	ChatRoomID  int64         `json:"chat_room_id"` // Chat room ID
	Timestamp   time.Time     `json:"timestamp"`    // Time of the message
	Status      MessageStatus `json:"status"`       // Status of the message
	ReadAt      *time.Time    `json:"read_at"`      // Time when the message was read (nullable)
}
