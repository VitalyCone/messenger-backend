package service

import (
	"fmt"

	"github.com/VitalyCone/websocket-messenger/internal/app/model"
	"github.com/redis/go-redis/v9"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)


type ChatService struct {
	db *gorm.DB
	rdb *redis.Client
}

func NewChatService(db *gorm.DB, rdb *redis.Client) *ChatService {
	return &ChatService{
		db: db,
		rdb: rdb,
	}
}

func (s *ChatService) CreateChat(chat *model.Chat) error {
	// var userIDs []uint

	if len(chat.Users) <= 1{
		return fmt.Errorf("chat users is empty")
	}

	for i, user := range chat.Users {
        var existingUser  model.User
        result := s.db.First(&existingUser , "username = ?", user.Username)
        if result.Error != nil {
            // If the user does not exist, return an error
            return fmt.Errorf("user %s does not exist : ", user.Username)
        }
		chat.Users[i] = existingUser
		// userIDs= append(userIDs, existingUser.ID)
    }
	if !chat.IsGroup{
		var existingChat model.Chat
		err := s.db.
			Model(&model.Chat{}).
			Joins("JOIN user_chats ON user_chats.chat_id = chats.id").
			Where("user_chats.user_id IN (?, ?)", chat.Users[0].ID, chat.Users[1].ID).
			Group("chats.id").
			Having("COUNT(DISTINCT user_chats.user_id) = 2").
			First(&existingChat).Error

		if err == nil {
			return fmt.Errorf("private chat between these users already exists")
		}
	}
	// chat.ChatKey = generateChatKey(userIDs)

	resoult := s.db.Create(&chat)
	if resoult.Error != nil {
		return resoult.Error
	}
	return nil
}

func (s *ChatService) GetChat(id uint) (model.Chat, error) {
	var chat model.Chat
	resoult := s.db.Preload("Users").First(&chat, id)
	if resoult.Error != nil{
		return model.Chat{}, resoult.Error
	}
	return chat, nil
}

func (s *ChatService) GetChat_ToResponse(id uint) (model.ChatResponse, error) {
	var chat model.Chat
	resoult := s.db.
		Preload("Users").
		Preload("Messages", func(db *gorm.DB) *gorm.DB {
			return db.Preload("Sender").
				Order("messages.created_at DESC").
				Limit(1)
		  }).
		First(&chat, id)
		
	if resoult.Error != nil{
		return model.ChatResponse{}, resoult.Error
	}
	chatResp := chat.ToResponse()
	return chatResp, nil
}

func (s *ChatService) GetChats(username string) ([]model.Chat, error) {
	var chats []model.Chat

    resoult := s.db.Model(&model.Chat{}).
        Preload("Users").
        Joins("JOIN user_chats ON user_chats.chat_id = chats.id").
        Joins("JOIN users ON users.id = user_chats.user_id").
        Where("users.username = ?", username).
        Find(&chats)

	if resoult.Error != nil{
		return chats, resoult.Error
	}
	return chats, nil
}

func (s *ChatService) GetChats_ToResponse(username string, offset, limit int) ([]model.ChatResponse, error) {
	var chats []model.Chat
	
    resoult := s.db.Model(&model.Chat{}).
		Preload("Users").
		Joins("JOIN user_chats ON user_chats.chat_id = chats.id").
		Joins("JOIN users ON users.id = user_chats.user_id").
		// Preload("Messages", func(db *gorm.DB) *gorm.DB {
		// 	return db.Preload("Sender").
		// 		Order("messages.created_at DESC").
		// 		Limit(1)
		// }).
		Where("users.username = ?", username).
		Offset(offset).
		Limit(limit).
		Find(&chats)
	
	if resoult.Error != nil{
		return nil, resoult.Error
	}

	chatResponses := make([]model.ChatResponse, len(chats))
	for i := range chats {
		err := s.db.Model(&chats[i]).
			Preload("Sender").
			Order("created_at DESC").
			Limit(1).
			Association("Messages").
			Find(&chats[i].Messages)
		
		if err != nil {
			return nil, err
		}
		chatResponses[i] = chats[i].ToResponse()
		logrus.Println(chatResponses[i].LastMessage)
	}
	return chatResponses, nil
}

func (s *ChatService) IsUserInChat(username string, chatID uint) bool{
	var count int64
    
    s.db.Model(&model.User{}).
        Joins("JOIN user_chats ON user_chats.user_id = users.id").
        Where("user_chats.chat_id = ? AND users.username = ?", chatID, username).
        Count(&count)
    
    return count > 0
}

//ЛИШНИЕ INSERTЫ USERS
func (s *ChatService) ModifyChatName(id uint, name string) error {
	err:= s.db.Model(model.Chat{}).
		Where(id).
		Update("name",name).
		Error
	return err
}

func (s *ChatService) ModifyChatUsers(id uint, users []model.User) error {
	var chat model.Chat
	if err := s.db.
		Preload("Users").
		First(&chat, id).
		Error; 
		err != nil{
		return err
	}

	if !chat.IsGroup{
		return fmt.Errorf("error: that's chat for 2 users only")
	}
	
	var fullUsers []model.User
	for _, user := range users {
		existingUser, err := getUserByUsername(user.Username, s.db, s.rdb)
		if err != nil {
			return err
		}
		fullUsers = append(fullUsers, existingUser)
	}
	
	err := s.db.
		Model(&chat).
		Association("Users").
		Replace(fullUsers)

	if err != nil{
		return err
	}
	
	return nil
}

// func generateChatKey(userIDs []uint) string {
//     ids := make([]uint, len(userIDs))
//     copy(ids, userIDs)
//     sort.Slice(ids, func(i, j int) bool { return ids[i] < ids[j] })
//     return fmt.Sprintf("%v", ids)
// }