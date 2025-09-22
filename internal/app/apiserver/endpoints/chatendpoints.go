package endpoints

import (
	"net/http"
	"slices"
	"strconv"

	"github.com/VitalyCone/websocket-messenger/internal/app/model"
	"github.com/gin-gonic/gin"
)

// @Summary Get user data
// @Schemes
// @Description Get user data by token
// @Security ApiKeyAuth
// @Tags Chats
// @Accept json
// @Produce json
// @Param createChatDto body model.CreateChatDto true "Create chat dto for register chat"
// @Success 201 {object} model.ChatResponse "chat responce"
// @Failure 400,404,401 {object} errorResponse
// @Failure 500 {object} errorResponse
// @Failure default {object} errorResponse
// @Router /v1/chats [POST]
func (ep *Endpoints) CreateChat(g *gin.Context){
	var createChatDto model.CreateChatDto

	tokenString := g.GetHeader("token")
    if tokenString == ""{
        newErrorResponse(g, http.StatusUnauthorized, "token nil")
        return
	}

	if err := g.BindJSON(&createChatDto); err != nil {
        newErrorResponse(g, http.StatusBadRequest, err.Error())
        return
    }

	username, err := ep.services.User.GetUsernameFromToken(tokenString)
	if err != nil{
		newErrorResponse(g, http.StatusUnauthorized, err.Error())
		return
	}

	chat := createChatDto.ToModel(username)
	if err = ep.services.Chat.CreateChat(&chat); err != nil{
		newErrorResponse(g, http.StatusInternalServerError, err.Error())
		return
	}

	chatResponse := chat.ToResponse()
	g.JSON(http.StatusCreated, chatResponse)
}

// @Summary Get user data
// @Schemes
// @Description Get user data by token
// @Security ApiKeyAuth
// @Tags Chats
// @Accept json
// @Produce json
// @Param offset query int false "offset of chats"
// @Param limit query int false "limit of chats"
// @Success 200 {object} []model.ChatResponse "chat response"
// @Failure 400,404,401 {object} errorResponse
// @Failure 500 {object} errorResponse
// @Failure default {object} errorResponse
// @Router /v1/chats [GET]
func (ep *Endpoints) GetChats(g *gin.Context){
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

	limit, err := strconv.Atoi(g.Query("limit"))
	if err != nil || limit < 0{
		limit = 15
	}
	offset, err := strconv.Atoi(g.Query("offset"))
	if err != nil || offset < 0{
		offset = 0
	}

	chatsResponces, err := ep.services.Chat.GetChats_ToResponse(username, offset, limit)
	if err != nil{
		newErrorResponse(g, http.StatusInternalServerError, err.Error())
		return
	}
	g.JSON(http.StatusOK, chatsResponces)
}

// @Summary Get user data
// @Schemes
// @Description Get user data by token
// @Security ApiKeyAuth
// @Tags Chats
// @Accept json
// @Produce json
// @Param id path int true "chat id"
// @Success 200 {object} model.ChatResponse "chat response"
// @Failure 400,404,401 {object} errorResponse
// @Failure 500 {object} errorResponse
// @Failure default {object} errorResponse
// @Router /v1/chats/{id} [GET]
func (ep *Endpoints) GetChat(g *gin.Context){
	tokenString := g.GetHeader("token")
    if tokenString == ""{
        newErrorResponse(g, http.StatusUnauthorized, "token nil")
        return
	}

	id, err := strconv.Atoi(g.Param("id"))
	if err != nil {
		newErrorResponse(g, http.StatusBadRequest, "Invalid type of id")
		return
	}

	username, err := ep.services.User.GetUsernameFromToken(tokenString)
	if err != nil{
		newErrorResponse(g, http.StatusUnauthorized, err.Error())
		return
	}

	chatResp, err := ep.services.Chat.GetChat_ToResponse(uint(id))
	if err != nil{
		newErrorResponse(g, http.StatusInternalServerError, err.Error())
		return
	}

	usernames := make([]string, len(chatResp.Users))
	for i, user := range chatResp.Users {
		usernames[i] = user.Username
	}

	if !slices.Contains(usernames, username){
		newErrorResponse(g, http.StatusUnauthorized, "you are not owner of this chat")
		return
	}

	g.JSON(http.StatusOK, chatResp)
}

// // @Summary Modify chat data
// // @Schemes
// // @Description Modify chat data
// // @Security ApiKeyAuth
// // @Tags Chats
// // @Accept json
// // @Produce json
// // @Param modifyChatDto body model.ModifyChatDto true "modify chat dto"
// // @Success 200 {object} model.ChatResponse "chat response"
// // @Failure 400,404,401 {object} errorResponse
// // @Failure 500 {object} errorResponse
// // @Failure default {object} errorResponse
// // @Router /chats/{id} [PUT]
// func (ep *Endpoints) ModifyChat(g *gin.Context){
// 	var modifyChatDto model.ModifyChatDto
// 	tokenString := g.GetHeader("token")
//     if tokenString == ""{
//         newErrorResponse(g, http.StatusUnauthorized, "token nil")
//         return
// 	}

// 	if err := g.BindJSON(&modifyChatDto); err != nil {
//         newErrorResponse(g, http.StatusBadRequest, err.Error())
//         return
//     }

// 	username, err := ep.services.User.GetUsernameFromToken(tokenString)
// 	if err != nil{
// 		newErrorResponse(g, http.StatusUnauthorized, err.Error())
// 		return
// 	}

// 	ok := ep.services.Chat.IsUserInChat(username, modifyChatDto.ID)
// 	if !ok{
// 		newErrorResponse(g, http.StatusForbidden, "user can't read this chat")
// 		return
// 	}

// 	chatModel := modifyChatDto.ToModel()
// 	err = ep.services.Chat.ModifyChat(&chatModel)
// 	if err != nil{
// 		newErrorResponse(g, http.StatusInternalServerError, err.Error())
// 		return
// 	}

// 	g.JSON(http.StatusOK, chatModel.ToResponse())
// }

// @Summary Modify chat data
// @Schemes
// @Description Modify chat data
// @Security ApiKeyAuth
// @Tags Chats
// @Accept json
// @Produce json
// @Param modifyChatDto body model.ModifyChatDto true "modify chat dto"
// @Success 200 {object} model.ChatResponse "chat response"
// @Failure 400,404,401 {object} errorResponse
// @Failure 500 {object} errorResponse
// @Failure default {object} errorResponse
// @Router /v1/chats/{id} [PATCH]
func (ep *Endpoints) ModifyChat(g *gin.Context){
	var modifyChatDto model.ModifyChatDto
	tokenString := g.GetHeader("token")
    if tokenString == ""{
        newErrorResponse(g, http.StatusUnauthorized, "token nil")
        return
	}

	if err := g.BindJSON(&modifyChatDto); err != nil {
        newErrorResponse(g, http.StatusBadRequest, err.Error())
        return
    }

	username, err := ep.services.User.GetUsernameFromToken(tokenString)
	if err != nil{
		newErrorResponse(g, http.StatusUnauthorized, err.Error())
		return
	}

	ok := ep.services.Chat.IsUserInChat(username, modifyChatDto.ID)
	if !ok{
		newErrorResponse(g, http.StatusForbidden, "user can't read this chat")
		return
	}

	chat := modifyChatDto.ToModel()
	if modifyChatDto.Name != nil{
		err := ep.services.Chat.ModifyChatName(chat.ID, chat.Name)
		if err != nil{
			newErrorResponse(g, http.StatusInternalServerError, err.Error())
			return
		}
	}
	if modifyChatDto.UserUsernames != nil{
		err := ep.services.Chat.ModifyChatUsers(chat.ID, chat.Users)
		if err != nil{
			newErrorResponse(g, http.StatusInternalServerError, err.Error())
			return
		}
	}

	chat, err = ep.services.Chat.GetChat(chat.ID)
	if err != nil{
		newErrorResponse(g, http.StatusInternalServerError, err.Error())
		return
	}
	
	g.JSON(http.StatusOK, chat.ToResponse())

	// chatModel := modifyChatDto.ToModel()
	// err = ep.services.Chat.ModifyChat(&chatModel)
	// if err != nil{
	// 	newErrorResponse(g, http.StatusInternalServerError, err.Error())
	// 	return
	// }

	// g.JSON(http.StatusOK, chatModel.ToResponse())
}