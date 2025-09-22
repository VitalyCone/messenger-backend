package service

import (
	"testing"

	"github.com/VitalyCone/websocket-messenger/internal/app/model"
	"github.com/go-redis/redismock/v9"
	"github.com/stretchr/testify/assert"
)

// Mock for Redis client is not used in current logic, but included for completeness.

func TestCreateMessage_InvalidContent(t *testing.T) {
	db := setupTestDB()
	rdb, _ := redismock.NewClientMock()
	service := NewMessageService(db, rdb)

	msg := &model.Message{Content: ""}
	err := service.CreateMessage(msg)
	assert.Error(t, err)
	assert.Equal(t, "invalid message", err.Error())
}

func TestCreateMessage_UserNotFound(t *testing.T) {
	db := setupTestDB()
	rdb, _ := redismock.NewClientMock()
	service := NewMessageService(db, rdb)

	msg := &model.Message{
		Content: "Hello",
		Sender:  model.User{Username: "nonexistent"},
	}
	err := service.CreateMessage(msg)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "user not found")
}

func TestCreateMessage_Success(t *testing.T) {
	db := setupTestDB()
	rdb, _ := redismock.NewClientMock()
	service := NewMessageService(db, rdb)

	user := model.User{Username: "alice"}
	db.Create(&user)
	chat := model.Chat{Name: "testchat"}
	db.Create(&chat)

	msg := &model.Message{
		Content:  "Hello",
		Sender:   model.User{Username: "alice"},
		ChatID:   chat.ID,
	}
	err := service.CreateMessage(msg)
	assert.NoError(t, err)
	assert.NotZero(t, msg.ID)
	assert.Equal(t, user.ID, msg.SenderID)
}

func TestGetMessages_Empty(t *testing.T) {
	db := setupTestDB()
	rdb, _ := redismock.NewClientMock()
	service := NewMessageService(db, rdb)

	messages, err := service.GetMessages(1, 10, 0)
	assert.NoError(t, err)
	assert.Len(t, messages, 0)
}

func TestGetMessages_Success(t *testing.T) {
	db := setupTestDB()
	rdb, _ := redismock.NewClientMock()
	service := NewMessageService(db, rdb)

	user := model.User{Username: "bob"}
	db.Create(&user)
	chat := model.Chat{Name: "chat1"}
	db.Create(&chat)
	msg := model.Message{
		Content:  "Hi",
		SenderID: user.ID,
		ChatID:   chat.ID,
	}
	db.Create(&msg)

	messages, err := service.GetMessages(chat.ID, 10, 0)
	assert.NoError(t, err)
	assert.Len(t, messages, 1)
	assert.Equal(t, "Hi", messages[0].Content)
}

func TestGetMessages_ToResponse(t *testing.T) {
	db := setupTestDB()
	rdb, _ := redismock.NewClientMock()
	service := NewMessageService(db, rdb)

	user := model.User{Username: "carol"}
	db.Create(&user)
	chat := model.Chat{Name: "chat2"}
	db.Create(&chat)
	msg := model.Message{
		Content:  "Hey",
		SenderID: user.ID,
		ChatID:   chat.ID,
	}
	db.Create(&msg)

	responses, err := service.GetMessages_ToResponse(chat.ID, 10, 0)
	assert.NoError(t, err)
	assert.Len(t, responses, 1)
	assert.Equal(t, "Hey", responses[0].Content)
}