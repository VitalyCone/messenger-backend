package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/VitalyCone/websocket-messenger/internal/app/apiserver"
	"github.com/VitalyCone/websocket-messenger/internal/app/service"
	"github.com/sirupsen/logrus"
	"gopkg.in/yaml.v3"
)

var (
	configPath  string
	dockerCheck string
	tokenPath   string
)

func init() {
	configPath = "config/apiserver.yaml"
	tokenPath = "config/token.yaml"
	dockerCheck = "DOCKER_ENV"
}

type ConfigData struct {
	ApiAddr            string `yaml:"api_addr"`
	DbUrl              string `yaml:"database_url"`
	DbDockerUrl        string `yaml:"database_docker_url"`
	RedisHostUrl       string `yaml:"redis_host"`
	RedisDockerHostUrl string `yaml:"redis_docker_host"`
	TestDbUrl          string `yaml:"test_database_url"`
	TestDbDockerUrl    string `yaml:"test_database_docker_url"`
}

type TokenData struct {
	Token string `yaml:"token"`
}

// @title Account API
// @version 1.0
// @description API for managing users
// @securityDefinitions.apiKey ApiKeyAuth
// @in header
// @name token
func main() {
	var configServer *apiserver.Config
	var configService *service.Config
	var redisUrl string
	var dbUrl string
	var testDbUrl string

	cfg := ConfigData{}
	token := TokenData{}

	isDocker := os.Getenv(dockerCheck) == "true"

	data, err := os.ReadFile(configPath)
	if err != nil {
		logrus.Fatal(err)
	}

	tokenData, err := os.ReadFile(tokenPath)
	if err != nil {
		logrus.Fatal(err)
	}

	err = yaml.Unmarshal(data, &cfg)
	if err != nil {
		logrus.Fatal(err)
	}

	err = yaml.Unmarshal(tokenData, &token)
	if err != nil {
		logrus.Fatal(err)
	}

	if isDocker {
		logrus.Println("App running in Docker. Using Docker database url")
		dbUrl = cfg.DbDockerUrl
		redisUrl = cfg.RedisDockerHostUrl
		testDbUrl = cfg.TestDbDockerUrl
	} else {
		//TODO: FIX IT
		logrus.Println("App running without Docker. Using Local database url")
		dbUrl = cfg.DbDockerUrl
		redisUrl = cfg.RedisDockerHostUrl
		testDbUrl = cfg.TestDbUrl
	}

	configServer = apiserver.NewConfig(cfg.ApiAddr, dbUrl, testDbUrl)
	configService = service.NewConfig(dbUrl, redisUrl, token.Token)
	service := service.NewService(configService)
	server := apiserver.NewAPIServer(configServer, service)

	go func() {
		if err := server.Start(token.Token); err != nil {
			logrus.Fatal(err)
		}
	}()

	quit := make(chan os.Signal, 3)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM, syscall.SIGTSTP)

	<-quit
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := server.Close(ctx); err != nil {
		log.Fatalf("Forced shutdown: %v", err)
	}
	log.Println("Server stopped gracefully")
}
