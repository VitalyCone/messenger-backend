package apiserver

import (
	"time"
	"github.com/VitalyCone/websocket-messenger/internal/app/apiserver/chat"
	"github.com/VitalyCone/websocket-messenger/internal/app/apiserver/endpoints"
	"github.com/VitalyCone/websocket-messenger/internal/app/service"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"net/http"
	"context"
	//_ "github.com/swaggo/swag/example/basic/docs"
)

type APIServer struct {
	config *Config
	router *gin.Engine
	services *service.Service
	srv *http.Server
}

func NewAPIServer(config *Config, service *service.Service) *APIServer {
	return &APIServer{
		config: config,
		router: gin.Default(),
		services: service,
	}
}

func (s *APIServer) Start(tokenSignedString string) error {

	s.router.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"*"},                                       // Разрешенные источники
		AllowMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE"},   // Разрешенные методы
		AllowHeaders:     []string{"Origin", "Authorization", "Content-Type", "token"}, // Разрешенные заголовки
		ExposeHeaders:    []string{"Content-Length"},                          // Заголовки, которые могут быть доступны клиенту
		AllowCredentials: true,                                                // Разрешить отправку учетных данных (например, куки)
		MaxAge:           12 * time.Hour,                                      // Время кэширования preflight-запросов
	}))

	if err := s.services.Start(); err != nil{
		return err
	}

	hub := chat.NewHub(s.services)
	go hub.Run()
	
	endpoint := endpoints.NewEndpoints(s.services, s.router, hub)

	endpoint.InitRoutes()
	
	logrus.Printf("SWAGGER : http://localhost%s/swagger/index.html\n", s.config.ApiAddr)

	s.srv = &http.Server{
		Addr:    ":8080",
		Handler: s.router,
	}

	return s.router.Run(s.config.ApiAddr)
}

func (s *APIServer) Close(ctx context.Context) error {
	err := s.srv.Shutdown(ctx)
	if err != nil{
		return err
	}
	return nil
}