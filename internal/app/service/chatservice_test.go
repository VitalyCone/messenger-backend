package service

import (
	"encoding/json"
	"testing"

	"github.com/VitalyCone/websocket-messenger/internal/app/model"
	"github.com/go-redis/redismock/v9"
	"github.com/stretchr/testify/assert"
)

// func setupTestDB() *gorm.DB {
// 	db, _ := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
// 	db.AutoMigrate(&model.User{}, &model.Chat{}, &model.Message{})
// 	return db
// }

func TestCreateChat_InvalidUsers(t *testing.T) {
	db := setupTestDB()
	rdb, _ := redismock.NewClientMock()
	service := NewChatService(db, rdb)

	chat := &model.Chat{Users: []model.User{}}
	err := service.CreateChat(chat)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "chat users is empty")
}

func TestCreateChat_UserDoesNotExist(t *testing.T) {
	db := setupTestDB()
	rdb, _ := redismock.NewClientMock()
	service := NewChatService(db, rdb)

	chat := &model.Chat{Users: []model.User{{Username: "ghost"}, {Username: "phantom"}}}
	err := service.CreateChat(chat)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "does not exist")
}

func TestCreateChat_PrivateChatAlreadyExists(t *testing.T) {
	db := setupTestDB()
	rdb, _ := redismock.NewClientMock()
	service := NewChatService(db, rdb)

	user1 := model.User{Username: "alice"}
	user2 := model.User{Username: "bob"}
	db.Create(&user1)
	db.Create(&user2)
	chat := model.Chat{Users: []model.User{user1, user2}, IsGroup: false}
	db.Create(&chat)

	chat2 := &model.Chat{Users: []model.User{{Username: "alice"}, {Username: "bob"}}, IsGroup: false}
	err := service.CreateChat(chat2)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "already exists")
}

func TestCreateChat_Success(t *testing.T) {
	db := setupTestDB()
	rdb, _ := redismock.NewClientMock()
	service := NewChatService(db, rdb)

	user1 := model.User{Username: "alice"}
	user2 := model.User{Username: "bob"}
	db.Create(&user1)
	db.Create(&user2)
	chat := &model.Chat{Users: []model.User{{Username: "alice"}, {Username: "bob"}}, IsGroup: false}
	err := service.CreateChat(chat)
	assert.NoError(t, err)
	assert.NotZero(t, chat.ID)
}

func TestGetChat_NotFound(t *testing.T) {
	db := setupTestDB()
	rdb, _ := redismock.NewClientMock()
	service := NewChatService(db, rdb)

	_, err := service.GetChat(999)
	assert.Error(t, err)
}

func TestGetChat_Success(t *testing.T) {
	db := setupTestDB()
	rdb, _ := redismock.NewClientMock()
	service := NewChatService(db, rdb)

	user := model.User{Username: "alice"}
	db.Create(&user)
	chat := model.Chat{Name: "test", Users: []model.User{user}, IsGroup: true}
	db.Create(&chat)

	got, err := service.GetChat(chat.ID)
	assert.NoError(t, err)
	assert.Equal(t, chat.ID, got.ID)
}

func TestGetChat_ToResponse(t *testing.T) {
	db := setupTestDB()
	rdb, _ := redismock.NewClientMock()
	service := NewChatService(db, rdb)

	user := model.User{Username: "alice"}
	db.Create(&user)
	chat := model.Chat{Name: "test", Users: []model.User{user}, IsGroup: true}
	db.Create(&chat)

	resp, err := service.GetChat_ToResponse(chat.ID)
	assert.NoError(t, err)
	assert.Equal(t, chat.ID, resp.ID)
}

func TestGetChats_Empty(t *testing.T) {
	db := setupTestDB()
	rdb, _ := redismock.NewClientMock()
	service := NewChatService(db, rdb)

	chats, err := service.GetChats("nobody")
	assert.NoError(t, err)
	assert.Len(t, chats, 0)
}

func TestGetChats_Success(t *testing.T) {
	db := setupTestDB()
	rdb, _ := redismock.NewClientMock()
	service := NewChatService(db, rdb)

	user := model.User{Username: "alice"}
	db.Create(&user)
	chat := model.Chat{Name: "test", Users: []model.User{user}, IsGroup: true}
	db.Create(&chat)
	db.Model(&chat).Association("Users").Append(&user)

	chats, err := service.GetChats("alice")
	assert.NoError(t, err)
	assert.Len(t, chats, 1)
	assert.Equal(t, chat.ID, chats[0].ID)
}

func TestIsUserInChat(t *testing.T) {
	db := setupTestDB()
	rdb, _ := redismock.NewClientMock()
	service := NewChatService(db, rdb)

	user := model.User{Username: "alice"}
	db.Create(&user)
	chat := model.Chat{Name: "test", Users: []model.User{user}, IsGroup: true}
	db.Create(&chat)
	db.Model(&chat).Association("Users").Append(&user)

	in := service.IsUserInChat("alice", chat.ID)
	assert.True(t, in)
	out := service.IsUserInChat("bob", chat.ID)
	assert.False(t, out)
}

func TestModifyChatName(t *testing.T) {
	db := setupTestDB()
	rdb, _ := redismock.NewClientMock()
	service := NewChatService(db, rdb)

	chat := model.Chat{Name: "old", IsGroup: true}
	db.Create(&chat)

	err := service.ModifyChatName(chat.ID, "new")
	assert.NoError(t, err)
	var updated model.Chat
	db.First(&updated, chat.ID)
	assert.Equal(t, "new", updated.Name)
}

func TestModifyChatUsers_Group(t *testing.T) {
	db := setupTestDB()
	rdb, mock := redismock.NewClientMock()

	user1 := model.User{Username: "alice"}
	user2 := model.User{Username: "bob"}

	service := NewChatService(db, rdb)
	
	db.Create(&user1)
	db.Create(&user2)

	user2Byte, _ := json.Marshal(user2)
	mock.ExpectSet("user_bob", user2Byte, 0).SetVal("OK")
	mock.ExpectGet("user_bob").SetVal(string(user2Byte))

	chat := model.Chat{Name: "group", IsGroup: true}
	db.Create(&chat)
	db.Model(&chat).Association("Users").Append(&user1)

	err := service.ModifyChatUsers(chat.ID, []model.User{user2})
	assert.NoError(t, err)
	var updated model.Chat
	db.Preload("Users").First(&updated, chat.ID)
	assert.Len(t, updated.Users, 1)
	assert.Equal(t, "bob", updated.Users[0].Username)
}

func TestModifyChatUsers_Private(t *testing.T) {
	db := setupTestDB()
	rdb, _ := redismock.NewClientMock()
	service := NewChatService(db, rdb)

	user1 := model.User{Username: "alice"}
	db.Create(&user1)
	chat := model.Chat{Name: "private", IsGroup: false}
	db.Create(&chat)
	db.Model(&chat).Association("Users").Append(&user1)

	err := service.ModifyChatUsers(chat.ID, []model.User{{Username: "alice"}})
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "for 2 users only")
}
