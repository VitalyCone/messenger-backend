package service

import (
	"encoding/json"
	"testing"

	"github.com/VitalyCone/websocket-messenger/internal/app/model"
	"github.com/go-redis/redismock/v9"
	"github.com/stretchr/testify/assert"
	"gorm.io/gorm"
)

const testTokenKey = "test_secret_key"

func TestRegisterUser_Success(t *testing.T) {
	db := setupTestDB()
	rdb, mock := redismock.NewClientMock()
	service := NewUserService(db, rdb, testTokenKey)

	user := model.User{Model: gorm.Model{ID:1}, Username: "alice", Password: "password123"}
	// 	mock.ExpectSet("user_bob", user2Byte, 0).SetVal("OK")
	// mock.ExpectGet("user_bob").SetVal(string(user2Byte))
	userBytes, _ := json.Marshal(user)
	mock.ExpectSet("user_alice", userBytes, 0).SetVal("OK")

	token, err := service.RegisterUser(user)
	assert.NoError(t, err)
	assert.NotEmpty(t, token)
}

func TestRegisterUser_EmptyPassword(t *testing.T) {
	db := setupTestDB()
	rdb, _ := redismock.NewClientMock()
	service := NewUserService(db, rdb, testTokenKey)

	user := model.User{Username: "bob", Password: ""}
	_, err := service.RegisterUser(user)
	assert.Error(t, err)
}

func TestRegisterUser_DuplicateUser(t *testing.T) {
	db := setupTestDB()
	rdb, _ := redismock.NewClientMock()
	service := NewUserService(db, rdb, testTokenKey)

	user := model.User{Username: "carol", Password: "pass"}
	db.Create(&user)
	_, err := service.RegisterUser(user)
	assert.Error(t, err)
}

func TestLoginUser_Success_DB(t *testing.T) {
	db := setupTestDB()
	rdb, mock := redismock.NewClientMock()
	service := NewUserService(db, rdb, testTokenKey)

	pass := "mypassword"
	user := model.User{Username: "dave", Password: pass}
	service.RegisterUser(user)

	// Redis miss, fallback to DB
	mock.ExpectGet("user_dave").RedisNil()

	loginUser := model.User{Username: "dave", Password: pass}
	token, err := service.LoginUser(loginUser)
	assert.NoError(t, err)
	assert.NotEmpty(t, token)
}

func TestLoginUser_Success_Redis(t *testing.T) {
	db := setupTestDB()
	rdb, mock := redismock.NewClientMock()
	service := NewUserService(db, rdb, testTokenKey)

	user := model.User{Username: "eve", Password: "secret"}
	service.RegisterUser(user)

	// Simulate user in Redis
	var dbUser model.User
	db.Where("username = ?", "eve").First(&dbUser)
	userBytes, _ := json.Marshal(dbUser)
	mock.ExpectGet("user_eve").SetVal(string(userBytes))
	mock.ExpectSet("user_eve", userBytes, 0).SetVal("OK")

	loginUser := model.User{Username: "eve", Password: "secret"}
	token, err := service.LoginUser(loginUser)
	assert.NoError(t, err)
	assert.NotEmpty(t, token)
}

func TestLoginUser_UserNotFound(t *testing.T) {
	db := setupTestDB()
	rdb, mock := redismock.NewClientMock()
	service := NewUserService(db, rdb, testTokenKey)

	mock.ExpectGet("user_ghost").RedisNil()
	loginUser := model.User{Username: "ghost", Password: "pass"}
	_, err := service.LoginUser(loginUser)
	assert.Error(t, err)
}

func TestLoginUser_WrongPassword(t *testing.T) {
	db := setupTestDB()
	rdb, mock := redismock.NewClientMock()
	service := NewUserService(db, rdb, testTokenKey)

	user := model.User{Username: "frank", Password: "rightpass"}
	service.RegisterUser(user)
	mock.ExpectGet("user_frank").RedisNil()

	loginUser := model.User{Username: "frank", Password: "wrongpass"}
	_, err := service.LoginUser(loginUser)
	assert.Error(t, err)
}

func TestGetUsernameFromToken_Success(t *testing.T) {
	db := setupTestDB()
	rdb, _ := redismock.NewClientMock()
	service := NewUserService(db, rdb, testTokenKey)

	user := model.User{Username: "grace", Password: "pw"}
	token, _ := service.RegisterUser(user)

	username, err := service.GetUsernameFromToken(token.(string))
	assert.NoError(t, err)
	assert.Equal(t, "grace", username)
}

func TestGetUsernameFromToken_Invalid(t *testing.T) {
	db := setupTestDB()
	rdb, _ := redismock.NewClientMock()
	service := NewUserService(db, rdb, testTokenKey)

	_, err := service.GetUsernameFromToken("invalidtoken")
	assert.Error(t, err)
}

func TestGetUserData_Success(t *testing.T) {
	db := setupTestDB()
	rdb, mock := redismock.NewClientMock()
	service := NewUserService(db, rdb, testTokenKey)

	user := model.User{Username: "henry", Password: "pw"}
	service.RegisterUser(user)
	mock.ExpectGet("user_henry").RedisNil()

	token, _ := service.LoginUser(model.User{Username: "henry", Password: "pw"})
	data, err := service.GetUserData(token.(string))
	assert.NoError(t, err)
	assert.Equal(t, "henry", data.Username)
}

func TestGetUserData_InvalidToken(t *testing.T) {
	db := setupTestDB()
	rdb, _ := redismock.NewClientMock()
	service := NewUserService(db, rdb, testTokenKey)

	_, err := service.GetUserData("badtoken")
	assert.Error(t, err)
}

func TestGetUsersWithQuery_ToResponse_Empty(t *testing.T) {
	db := setupTestDB()
	rdb, _ := redismock.NewClientMock()
	service := NewUserService(db, rdb, testTokenKey)

	resp, err := service.GetUsersWithQuery_ToResponse("nobody", 0, 10)
	assert.NoError(t, err)
	assert.Len(t, resp, 0)
}

func TestGetUsersWithQuery_ToResponse_Success(t *testing.T) {
	db := setupTestDB()
	rdb, _ := redismock.NewClientMock()
	service := NewUserService(db, rdb, testTokenKey)

	user1 := model.User{Username: "ivan", Password: "pw"}
	user2 := model.User{Username: "ivanov", Password: "pw"}
	service.RegisterUser(user1)
	service.RegisterUser(user2)

	resp, err := service.GetUsersWithQuery_ToResponse("ivan", 0, 10)
	assert.NoError(t, err)
	assert.True(t, len(resp) >= 2)
}
