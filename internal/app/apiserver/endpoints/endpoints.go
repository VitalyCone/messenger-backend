package endpoints

import (
	"github.com/VitalyCone/websocket-messenger/docs"
	"github.com/VitalyCone/websocket-messenger/internal/app/apiserver/chat"
	"github.com/VitalyCone/websocket-messenger/internal/app/service"
	"github.com/gin-gonic/gin"
	swaggerfiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

var (
	mainPath string = "/api"
)


type Endpoints struct {
	services *service.Service
	router *gin.Engine
	hub *chat.Hub
}

func NewEndpoints(service *service.Service, router *gin.Engine, hub *chat.Hub) *Endpoints {
	return &Endpoints{
		services : service,
		router : router,
		hub : hub,
	}
}

func (e *Endpoints) InitRoutes() {
	docs.SwaggerInfo.BasePath = mainPath
	path := e.router.Group(mainPath)
	{
		v1 := path.Group("/v1")
		{
			
			v1.GET("/chatws", e.ConnectUserToChats)
			
			v1.GET("/accounts", e.GetUsersByUsername)
			v1.GET("/account", e.GetUserData)
			acc:= v1.Group("/account")
			{
				acc.POST("/register" , e.RegisterUser)
				acc.POST("/login" , e.LoginUser)
			}
			
			chat:= v1.Group("/chats")
			{
				chat.POST("/", e.CreateChat)
				chat.GET("/", e.GetChats)
				chat.GET("/:id", e.GetChat)
				chat.PATCH("/:id", e.ModifyChat)
				chat.GET("/:id/messages", e.GetMessages)
			}
			message := v1.Group("/messages")
			{
				message.POST("/", e.CreateMessage)
			}
			
		}
	}
	e.router.GET("/", e.Ping)
	e.router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerfiles.Handler))
}

func (ep *Endpoints) Ping(g *gin.Context){
	g.JSON(200, "PING")
}