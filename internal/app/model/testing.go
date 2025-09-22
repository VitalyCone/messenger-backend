package model

import "testing"

func TestCreateUserDto(t *testing.T) *CreateUserDto {
	return &CreateUserDto{
		Username:   "username",
		Password:   "password",
		FirstName:  "firname",
		SecondName: "secname",
		Role:       "user",
	}
}

func TestUsers(t *testing.T) []User{
	return []User{
		{
			Username: "u1",
		},
		{
			Username: "u2",
		},
		{
			Username: "u3",
		},
	}
}

func TestCreateChatDto(t *testing.T) *CreateChatDto {
	return &CreateChatDto{
		Name:        "name",
		CompanionsUsernames: []string{
			"u1",
			"u2",
			"u3",
		},
	}
}

func TestChat(t *testing.T) *Chat {
	return &Chat{
		Name:        "name",
		Users: TestUsers(t),
	}
}
