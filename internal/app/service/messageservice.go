package service

import (
	"errors"
	"fmt"

	"github.com/VitalyCone/websocket-messenger/internal/app/model"
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)



type MessageService struct {
	db *gorm.DB
	rdb *redis.Client
}

func NewMessageService(db *gorm.DB, rdb *redis.Client) *MessageService {
	return &MessageService{
		db: db,
		rdb: rdb,
	}
}

func (s *MessageService) CreateMessage(message *model.Message) error {
	if message.Content == ""{
		return errors.New("invalid message")
	}
	if message.Sender.Username != "" && message.SenderID == 0 {
        var user model.User
        if err := s.db.Where(model.User{Username: message.Sender.Username}).First(&user).Error; err != nil {
            return fmt.Errorf("user not found: %v", err)
        }
		message.SenderID = user.ID
		message.Sender = model.User{}
    }
	resoult := s.db.Create(&message)
	if resoult.Error != nil {
		return resoult.Error
	}
	if err := s.db.Preload("Sender").Preload("Chat").First(message, message.ID).Error; err != nil {
        return err
    }
	return nil
}

func (s *MessageService) GetMessages(chatID uint, limit, offset int) ([]model.Message, error) {
	messages := make([]model.Message, 0)
	
	resoult := s.db.Model(model.Message{}).
		Preload("Chat").
		Preload("Sender").
		Where("chat_id = ?", chatID).
		Offset(offset).
		Limit(limit).
		Find(&messages)
	if resoult.Error != nil {
		return messages, resoult.Error
	}
	return messages, nil
}

func (s *MessageService) GetMessages_ToResponse(chatID uint, limit, offset int) ([]model.MessageResponse, error) {
	messages := make([]model.Message, 0)
	
	resoult := s.db.Model(model.Message{}).
		Preload("Chat").
		Preload("Sender").
		Where("chat_id = ?", chatID).
		Order("created_at DESC").
		Offset(offset).
		Limit(limit).
		Find(&messages)
	
	if resoult.Error != nil {
		return nil, resoult.Error
	}
	respMessages := make([]model.MessageResponse, len(messages))
	for i, message := range messages {
		respMessages[i] = message.ToResponse()
	}

	return respMessages, nil
}