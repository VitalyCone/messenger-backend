package endpoints

import (
	"net/http"
	"strconv"
	"strings"

	"github.com/VitalyCone/websocket-messenger/internal/app/model"
	"github.com/gin-gonic/gin"
)

// @Summary Register for user
// @Schemes
// @Description Register in api
// @Tags User
// @Accept json
// @Produce json
// @Param createUserDto body model.CreateUserDto true "Create user dto for register in"
// @Success 201 {string} string "token"
// @Failure 400,404 {object} errorResponse
// @Failure 500 {object} errorResponse
// @Failure default {object} errorResponse
// @Router /v1/account/register [POST]
func (ep *Endpoints) RegisterUser(g *gin.Context){
    var CreateUserDto model.CreateUserDto

    if err := g.BindJSON(&CreateUserDto); err != nil {
        newErrorResponse(g, http.StatusBadRequest, err.Error())
        return
    }

    user,err := CreateUserDto.ToModel()
    if err != nil {
        newErrorResponse(g, http.StatusBadRequest, err.Error())
        return
    }

    token, err := ep.services.User.RegisterUser(user)
    if err != nil {
        newErrorResponse(g, http.StatusInternalServerError, err.Error())
        return
    }

    g.JSON(http.StatusCreated, gin.H{"token": token})
}

// @Summary Login for user
// @Schemes
// @Description Login in api
// @Tags User
// @Accept json
// @Produce json
// @Param userDto body model.UserDto true "Login user dt"
// @Success 200 {string} string "token"
// @Failure 400,404 {object} errorResponse
// @Failure 500 {object} errorResponse
// @Failure default {object} errorResponse
// @Router /v1/account/login [POST]
func (ep *Endpoints) LoginUser(g *gin.Context){
    var userDto model.UserDto

    if err := g.BindJSON(&userDto); err != nil {
        newErrorResponse(g, http.StatusBadRequest, err.Error())
        return
    }

    user := userDto.ToModel()

    token, err := ep.services.User.LoginUser(user)
    if err != nil {
        newErrorResponse(g, http.StatusInternalServerError, err.Error())
        return
    }

    g.JSON(http.StatusOK, gin.H{"token": token})
}

// @Summary Get user data
// @Schemes
// @Description Get user data by token
// @Security ApiKeyAuth
// @Tags User
// @Accept json
// @Produce json
// @Success 200 {object} model.UserResponse "user responce"
// @Failure 400,404,401 {object} errorResponse
// @Failure 500 {object} errorResponse
// @Failure default {object} errorResponse
// @Router /v1/account [GET]
func (ep *Endpoints) GetUserData(g *gin.Context){
    tokenString := g.GetHeader("token")
    if tokenString == ""{
        newErrorResponse(g, http.StatusUnauthorized, "token nil")
        return
    }
    _, err := ep.services.User.GetUsernameFromToken(tokenString)
    if err != nil{
        newErrorResponse(g, http.StatusUnauthorized, err.Error())
    }

    user, err := ep.services.User.GetUserData(tokenString)
    if err != nil {
        newErrorResponse(g, http.StatusInternalServerError, err.Error())
        return
    }
    userResponse := user.ToResponse()
    g.JSON(http.StatusOK, userResponse)
}

// @Summary Get user data
// @Schemes
// @Description Get user data by token
// @Security ApiKeyAuth
// @Tags User
// @Accept json
// @Produce json
// @Param username query string true "username query"
// @Param offset query int false "offset from first responses"
// @Param limit query int false "restriction on return of publications"
// @Success 200 {object} []model.UserResponse "users response"
// @Failure 400,404,401 {object} errorResponse
// @Failure 500 {object} errorResponse
// @Failure default {object} errorResponse
// @Router /v1/accounts [GET]
func (ep *Endpoints) GetUsersByUsername(g *gin.Context){
    var offset , limit int
    var username string

    queryOffset,exists := g.GetQuery("offset")
    if exists{
        queryOffset, err := strconv.Atoi(queryOffset)
        if err != nil{
            newErrorResponse(g, http.StatusBadRequest, "offset is not int")
            return
        }
        offset = queryOffset
        if offset < 0{
            newErrorResponse(g, http.StatusBadRequest, "offset is not positive")
            return
        }
    }else{offset = 0}

    queryLimit,exists := g.GetQuery("limit")
    if exists{
        queryLimit, err := strconv.Atoi(queryLimit)
        if err != nil{
            newErrorResponse(g, http.StatusBadRequest, "offset is not int")
            return
        }
        limit = queryLimit
        if limit < 0{
            newErrorResponse(g, http.StatusBadRequest, "offset is not positive")
            return
        }
    }else{limit = 10000}

    usernameQuery ,exists := g.GetQuery("username")
    if exists{
        username = strings.TrimSpace(usernameQuery)
    }else{
        newErrorResponse(g, http.StatusBadRequest, "username is not provided")
    }


    usersResp, err := ep.services.User.GetUsersWithQuery_ToResponse(username, offset, limit)
    if err != nil{
        newErrorResponse(g,http.StatusInternalServerError, err.Error())
    }
    g.JSON(http.StatusOK, usersResp)
}