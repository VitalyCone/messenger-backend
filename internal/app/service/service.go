package service

import (
	"context"
	"fmt"

	"github.com/VitalyCone/websocket-messenger/internal/app/model"
	"github.com/redis/go-redis/v9"
	"github.com/sirupsen/logrus"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

//go:generate mockgen -source=service.go -destination=mocks/mock.go

type Service struct {
	User
	Chat
	Message
	config *Config
	db *gorm.DB
	rdb *redis.Client
}

func NewService(config *Config) *Service {
	return &Service{
		config: config,
	}
}

type User interface {
	RegisterUser(m model.User) (interface{}, error)
	GetUsernameFromToken(tokenString string) (string, error)
	LoginUser(m model.User) (interface{}, error)
	GetUserData(tokenString string) (model.User, error)
	GetUsersWithQuery_ToResponse(username string, offset, limit int) ([]model.UserResponse, error)
}

type Chat interface {
	CreateChat(chat *model.Chat) error
	GetChat(id uint) (model.Chat, error)
	GetChat_ToResponse(id uint) (model.ChatResponse, error)
	GetChats_ToResponse(username string, offset, limit int) ([]model.ChatResponse, error)
	GetChats(username string) ([]model.Chat, error)
	ModifyChatName(id uint, name string) error
	ModifyChatUsers(id uint, users []model.User) error 
	IsUserInChat(username string, chatID uint) bool
}

type Message interface {
	CreateMessage(message *model.Message) error 
	GetMessages(chatID uint, limit, offset int) ([]model.Message, error) 
	GetMessages_ToResponse(chatID uint, limit, offset int) ([]model.MessageResponse, error)
}


func (s *Service) Start() error {
	db, err := gorm.Open(postgres.Open(s.config.DatabaseURL), &gorm.Config{})
	if err != nil {
		return err
	}

	sqlDB, err := db.DB()
	if err != nil {
		return err
	}

	if err := sqlDB.Ping(); err != nil {
		return err
	}

	if err := s.migrateDatabase(db); err != nil {
		return err
	}

	logrus.Println("Migrations updated")
	logrus.Println("Database is working!")

	rdb, err := s.configureRedis(s.config.RedisURL)
	if err != nil {
		return err
	}

	logrus.Println("Redis is working!")

	s.db = db
	s.rdb = rdb

	s.User = NewUserService(db, rdb, s.config.TokenKey)
	s.Chat = NewChatService(db, rdb)
	s.Message = NewMessageService(db,rdb)

	return nil
}

func (s *Service) migrateDatabase(db *gorm.DB) error {
	// Отключаем проверку внешних ключей на время миграции
	db.Config.DisableForeignKeyConstraintWhenMigrating = true

	// Порядок важен: сначала таблицы без зависимостей, затем зависимые
	err := db.AutoMigrate(
		&model.User{},
		&model.Chat{},
		&model.Message{},
	)
	if err != nil {
		return fmt.Errorf("failed to migrate models: %v", err)
	}

	// Включаем проверку внешних ключей обратно
	db.Config.DisableForeignKeyConstraintWhenMigrating = false

	// Создаем таблицу many2many связей (если не создалась автоматически)
	if !db.Migrator().HasTable("user_chats") {
		err = db.Table("user_chats").AutoMigrate(&struct {
			UserID uint `gorm:"primaryKey"`
			ChatID uint `gorm:"primaryKey"`
		}{})
		if err != nil {
			return fmt.Errorf("failed to migrate user_chats table: %v", err)
		}
	}

	if !db.Migrator().HasConstraint(&model.Message{}, "Sender") {
		err = db.Migrator().CreateConstraint(&model.Message{}, "Sender")
		logrus.Println("Added Sender fk")
		if err != nil {
			return fmt.Errorf("failed to create Sender foreign key: %v", err)
		}
	}

	if !db.Migrator().HasConstraint(&model.Message{}, "Chat") {
		err = db.Migrator().CreateConstraint(&model.Message{}, "Chat")
		logrus.Println("Added Chat fk")
		if err != nil {
			return fmt.Errorf("failed to create Chat foreign key: %v", err)
		}
	}

	// if !db.Migrator().HasConstraint(&model.Chat{}, "LastMessage") {
	// 	err = db.Migrator().CreateConstraint(&model.Chat{}, "LastMessage")
	// 	logrus.Println("Added LastMessage fk")
	// 	if err != nil {
	// 		return fmt.Errorf("failed to create LastMessage foreign key: %v", err)
	// 	}
	// }

	// if !db.Migrator().HasConstraint(&model.Chat{}, "LastSender") {
	// 	err = db.Migrator().CreateConstraint(&model.Chat{}, "LastSender")
	// 	logrus.Println("Added LastSender fk")
	// 	if err != nil {
	// 		return fmt.Errorf("failed to create LastSender foreign key: %v", err)
	// 	}
	// }

	// Для сложных связей можно добавить дополнительные индексы
	// if !s.db.Migrator().HasIndex(&model.Message{}, "SenderID") {
	//     err = s.db.Migrator().CreateIndex(&model.Message{}, "SenderID")
	//     if err != nil {
	//         return fmt.Errorf("failed to create SenderID index: %v", err)
	//     }
	// }

	// if !s.db.Migrator().HasIndex(&model.Message{}, "ChatID") {
	//     err = s.db.Migrator().CreateIndex(&model.Message{}, "ChatID")
	//     if err != nil {
	//         return fmt.Errorf("failed to create ChatID index: %v", err)
	//     }
	// }

	return nil
}

func (s *Service) configureRedis(addr string) (*redis.Client, error) {
	rdb := redis.NewClient(
		&redis.Options{
			Addr:     addr,
			Password: "", // no password set
			DB:       0,  // use default DB
		})
	ctx := context.Background()
	_, err := rdb.Ping(ctx).Result()
	return rdb, err
}
