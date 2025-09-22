package model

import (
	"strings"
	"time"

	"github.com/go-playground/validator"
	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	Avatar       []byte  `json:"avatar"`
	Username     string  `json:"username" gorm:"index;unique"`
	PasswordHash string  `json:"password_hash"`
	Password     string  `json:"password" gorm:"-"`
	FirstName    string  `json:"firstname"`
	SecondName   string  `json:"secondname"`
	Role         string  `json:"role" gorm:"not null;default:user"`
	Balance      float32 `json:"balance"`
	Chats        []*Chat    `json:"chats" gorm:"many2many:user_chats;constraint:OnUpdate:CASCADE,OnDelete:SET NULL;"`
    Messages     []Message  `json:"messages" gorm:"foreignKey:SenderID"`
}


func (m *User) ToResponse() UserResponse {
	return UserResponse{
		ID:         m.ID,
		Avatar:     m.Avatar,
		Username:   m.Username,
		FirstName:  m.FirstName,
		SecondName: m.SecondName,
		Balance:    m.Balance,
		Role:       m.Role,
		CreatedAt:  m.CreatedAt,
		UpdatedAt:  m.UpdatedAt,
	}
}
type CreateUserDto struct {
	Username   string `json:"username" form:"username" validate:"required,alphanum,min=3,max=32"`
	Password   string `json:"password" form:"password" validate:"required,min=3,max=32"`
	FirstName  string `json:"firstname" form:"firstname" validate:"max=50"`
	SecondName string `json:"secondname" form:"secondname" validate:"max=50"`
	Role       string `json:"role" form:"role" validate:"required,oneof=user admin"` //"user"/"admin"
}

func (c *CreateUserDto) ToModel() (User, error) {
	validate := validator.New()
	if err := validate.Struct(c); err != nil {
		return User{}, err
	}
	return User{
		Username:   strings.ToLower(c.Username),
		Password:   c.Password,
		FirstName:  c.FirstName,
		SecondName: c.SecondName,
		Role:       c.Role,
	}, nil
}

type ModifyUserDto struct {
	Username    string  `json:"username" form:"username"`
	OldPassword string  `json:"oldPassword" form:"oldPassword"`
	NewPassword string  `json:"newPassword" form:"newPassword"`
	Avatar      []byte  `json:"avatar" form:"avatar"`
	FirstName   string  `json:"firstname" form:"firstname"`
	SecondName  string  `json:"secondname" form:"secondname"`
	Balance     float32 `json:"balance"`
}

func (u *ModifyUserDto) ToModel(passHash string) User {
	return User{}
}

type UserDto struct {
	Username string `json:"username" form:"username"`
	Password string `json:"password" form:"password"`
}

func (u *UserDto) ToModel() User{
	return User{
		Username: u.Username,
		Password: u.Password,
	}
}

type UserResponse struct {
	ID         uint      `json:"id"`
	Avatar     []byte    `json:"avatar"`
	Username   string    `json:"username"`
	FirstName  string    `json:"firstname"`
	SecondName string    `json:"secondname"`
	Balance    float32   `json:"balance"`
	Role       string    `json:"role"`
	CreatedAt  time.Time `json:"createdAt"`
	UpdatedAt  time.Time `json:"updatedAt"`
	DeletedAt *time.Time `json:"deletedAt"`
}
