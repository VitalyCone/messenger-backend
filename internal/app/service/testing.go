package service

import (
	"fmt"
	"log"
	"strings"
	"testing"

	"github.com/VitalyCone/websocket-messenger/internal/app/model"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func TestService(t *testing.T) (*Service, func(db *gorm.DB) ,error){

	type ConfigData struct {
		ApiAddr         string `yaml:"api_addr"`
		DbUrl           string `yaml:"database_url"`
		DbDockerUrl     string `yaml:"database_docker_url"`
		TestDbUrl       string `yaml:"test_database_url"`
		TestDbDockerUrl string `yaml:"test_database_docker_url"`
	}

	type TokenData struct {
		Token string `yaml:"token"`
	}

	cfg := ConfigData{}
	token := TokenData{}

	cfg.TestDbDockerUrl = "host=localhost port=5432 user=postgres dbname=websocket_messenger_test password=postgres sslmode=disable"
	token.Token = "x"
	configService := NewConfig(cfg.TestDbDockerUrl, "127.0.0.1:6379", token.Token)
	service := NewService(configService)

	if err := service.Start(); err != nil {
		return nil ,nil,err
	}
	trunicator := func(db *gorm.DB){
		tabl , err:= db.Migrator().GetTables()
		if err != nil {
			log.Println(err)
		}
		res := db.Exec(fmt.Sprintf("TRUNCATE %s CASCADE", strings.Join(tabl, ", ")))
		if res.Error != nil{
			log.Printf("%v", res.Error)
		}
	}
	return service,trunicator ,nil
}

func setupTestDB() *gorm.DB {
	db, _ := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	db.AutoMigrate(&model.User{}, &model.Chat{}, &model.Message{})
	return db
}