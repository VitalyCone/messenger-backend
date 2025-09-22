package model

import (
	"time"

	"gorm.io/gorm"
)

type Chat struct {
	gorm.Model
	Name     string
	// ChatKey  string    `gorm:"not null;unique"`
	IsGroup  bool      `gorm:"not null;default:false"`
	Users    []User    `gorm:"many2many:user_chats;constraint:OnUpdate:CASCADE,OnDelete:SET NULL;"`
	Messages []Message `gorm:"foreignKey:ChatID"`
}

func (c *Chat) ToResponse() ChatResponse {
	userResponse := make([]UserResponse, 0)
	for _, user := range c.Users {
		userResponse = append(userResponse, user.ToResponse())
	}

	var lastMessage MessageResponse
	if len(c.Messages) > 0 {
		// lastMes := &c.Messages
		// lastMessage = lastMes.ToResponse()
		lastMessage = c.Messages[0].ToResponse()
	}
	// if c.LastSender != nil {
	// 	lastSenderResponse := c.LastSender.ToResponse()
	// 	lastSender = &lastSenderResponse
	// }

	// if c.LastMessage != nil {
	// 	lastMessageResponse := c.LastMessage.ToResponse()
	// 	lastMessage = &lastMessageResponse
	// }

	return ChatResponse{
		ID:          c.ID,
		Name:        c.Name,
		IsGroup:     c.IsGroup,
		Users:       userResponse,
		LastMessage: &lastMessage,
		CreatedAt:   c.CreatedAt,
		UpdatedAt:   c.UpdatedAt,
	}
}

type CreateChatDto struct {
	Name                string   `json:"name"`
	CompanionsUsernames []string `json:"companions_usernames"`
	IsGroup             bool     `json:"isGroup"`
}

func (c *CreateChatDto) ToModel(ownUsername string) Chat {
	companions := make([]User, 0)
	companions = append(companions, User{Username: ownUsername})

	for _, username := range c.CompanionsUsernames {
		companions = append(companions, User{Username: username})
	}
	return Chat{
		Name:    c.Name,
		Users:   companions,
		IsGroup: c.IsGroup,
	}
}

type ChatResponse struct {
	ID          uint             `json:"id"`
	Name        string           `json:"name"`
	IsGroup     bool             `json:"isGroup"`
	Users       []UserResponse   `json:"users"`
	LastMessage *MessageResponse `json:"lastMessage"`
	CreatedAt   time.Time        `json:"createdAt"`
	UpdatedAt   time.Time        `json:"updatedAt"`
}

type ModifyChatDto struct {
	ID            uint     `json:"id"`
	Name          *string  `json:"name"`
	UserUsernames []string `json:"userUsernames"`
}

func (c *ModifyChatDto) ToModel() Chat {
	users := make([]User, len(c.UserUsernames))
	for i, username := range c.UserUsernames {
		users[i] = User{Username: username}
	}
	var name string
	if c.Name != nil {
		name = *c.Name
	}
	return Chat{
		Model: gorm.Model{
			ID: c.ID,
		},
		Name:  name,
		Users: users,
	}
}
