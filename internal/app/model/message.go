package model

import (
	"time"

	"gorm.io/gorm"
)

type Message struct {
	gorm.Model
	Content  string
	Sender   User `gorm:"foreignKey:SenderID"`
	SenderID uint // Важно: это внешний ключ
	Chat     Chat `gorm:"foreignKey:ChatID"`
	ChatID   uint // Внешний ключ для чата
}

func (m *Message) ToResponse() MessageResponse {
	return MessageResponse{
		ID:        m.ID,
		Content:   m.Content,
		Sender:    m.Sender.ToResponse(),
		Chat:      m.Chat.ToResponse(),
		CreatedAt: m.CreatedAt,
		UpdatedAt: m.UpdatedAt,
	}
}

type CreateMessageDto struct {
	Content string `json:"content"`
	ChatID  uint   `json:"chatId"`
}

func (m *CreateMessageDto) ToModel(senderUsername string) Message {
	return Message{
		Content: m.Content,
		ChatID:  m.ChatID,
		Sender:  User{Username: senderUsername},
	}
}

type MessageResponse struct {
	ID        uint         `json:"id"`
	Content   string       `json:"content"`
	Sender    UserResponse `json:"sender"`
	Chat      ChatResponse `json:"chat"`
	CreatedAt time.Time    `json:"createdAt"`
	UpdatedAt time.Time    `json:"updatedAt"`
}

type MessageWS struct {
	Type       string `json:"type"`
	Sender     string
	Recipients []string
	Content    string `json:"content"`
	ChatID     uint   `json:"chat_id"`
	// {
	// 	"type":"message",
	// 	"content": "проверка",
	// 	"chat_id": 15
	// }
}

func (m *MessageWS) ToModel() Message {
	return Message{
		Sender:  User{Username: m.Sender},
		Content: m.Content,
		ChatID:    m.ChatID,
	}
}
