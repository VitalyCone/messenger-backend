package endpoints

import (
	"net/http"
	"strconv"

	"github.com/VitalyCone/websocket-messenger/internal/app/model"
	"github.com/gin-gonic/gin"
)

// @Summary Get user data
// @Schemes
// @Description Get user data by token
// @Security ApiKeyAuth
// @Tags Messages
// @Accept json
// @Produce json
// @Param createMessageDto body model.CreateMessageDto true "Create message dto"
// @Success 201 {object} model.MessageResponse "message response"
// @Failure 400,404,401 {object} errorResponse
// @Failure 500 {object} errorResponse
// @Failure default {object} errorResponse
// @Router /v1/messages [POST]
func (ep *Endpoints) CreateMessage(g *gin.Context){
	var createMessageDto model.CreateMessageDto

	tokenString := g.GetHeader("token")
    if tokenString == ""{
        newErrorResponse(g, http.StatusUnauthorized, "token nil")
        return
	}

	if err := g.BindJSON(&createMessageDto); err != nil {
        newErrorResponse(g, http.StatusBadRequest, err.Error())
        return
    }

	username, err := ep.services.User.GetUsernameFromToken(tokenString)
	if err != nil{
		newErrorResponse(g, http.StatusUnauthorized, err.Error())
		return
	}

	message := createMessageDto.ToModel(username)
	if err = ep.services.Message.CreateMessage(&message); err != nil{
		newErrorResponse(g, http.StatusInternalServerError, err.Error())
	}

	messageResp := message.ToResponse()
	g.JSON(http.StatusCreated, messageResp)
}

// @Summary Get user data
// @Schemes
// @Description Get user data by token
// @Security ApiKeyAuth
// @Tags Messages,Chats
// @Accept json
// @Produce json
// @Param id path int true "chat id"
// @Param limit query int false "limit"
// @Param offest query int false "offest"
// @Success 200 {object} []model.MessageResponse "messages responsee"
// @Failure 400,404,401,403 {object} errorResponse
// @Failure 500 {object} errorResponse
// @Failure default {object} errorResponse
// @Router /v1/chats/{id}/messages [GET]
func (ep *Endpoints) GetMessages(g *gin.Context){
	tokenString := g.GetHeader("token")
    if tokenString == ""{
        newErrorResponse(g, http.StatusUnauthorized, "token nil")
        return
	}
	username, err := ep.services.User.GetUsernameFromToken(tokenString)
	if err != nil{
		newErrorResponse(g, http.StatusUnauthorized, err.Error())
		return
	}

	id, err := strconv.Atoi(g.Param("id"))
	if err != nil{
		newErrorResponse(g, http.StatusBadRequest, err.Error())
	}
	limit, err := strconv.Atoi(g.Query("limit"))
	if err != nil || limit < 0{
		limit = 15
	}
	offest, err := strconv.Atoi(g.Query("offest"))
	if err != nil || offest < 0{
		offest = 0
	}

	ok := ep.services.Chat.IsUserInChat(username, uint(id))
	if !ok{
		newErrorResponse(g, http.StatusForbidden, "user can't read this chat")
		return
	}

	modelResps, err := ep.services.Message.GetMessages_ToResponse(uint(id),limit,offest)
	if err != nil{
		newErrorResponse(g, http.StatusInternalServerError, err.Error())	
	}
	// if err = ep.services.Message.CreateMessage(&message); err != nil{
	// 	newErrorResponse(g, http.StatusInternalServerError, err.Error())
	// }

	g.JSON(http.StatusOK, modelResps)
}