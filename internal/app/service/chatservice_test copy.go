package service

import (
	"testing"

	"github.com/VitalyCone/websocket-messenger/internal/app/model"
	"github.com/stretchr/testify/assert"
)

func TestCreateChat(t *testing.T) {
	testCases := []struct {
		name     string
		s *Service
		u func() *model.CreateChatDto
		isValid bool
	}{
		{
			name:"valid",
			u: func() *model.CreateChatDto{
				return model.TestCreateChatDto(t)
			},
			isValid: true,
		},
		{
			name:"invalid",
			u: func() *model.CreateChatDto{
				m :=  model.TestCreateChatDto(t)
				return m
			},
			isValid: false,
		},
		{
			name:"cant create chat with empty companions",
			u: func() *model.CreateChatDto{
				m :=  model.TestCreateChatDto(t)
				m.CompanionsUsernames = []string{}
				return m
			},
			isValid: false,
		},
		{
			name:"can create chat with one person",
			u: func() *model.CreateChatDto{
				m :=  model.TestCreateChatDto(t)
				m.CompanionsUsernames = []string{m.CompanionsUsernames[1]}
				return m
			},
			isValid: true,
		},
		{
			name:"cant create second chat with one person",
			u: func() *model.CreateChatDto{
				m :=  model.TestCreateChatDto(t)
				m.CompanionsUsernames = []string{m.CompanionsUsernames[1]}
				return m
			},
			isValid: false,
		},
		{
			name:"cant create second chat with many person",
			u: func() *model.CreateChatDto{
				m :=  model.TestCreateChatDto(t)
				return m
			},
			isValid: false,
		},
		{
			name:"cant create second chat with many person with wrong order",
			u: func() *model.CreateChatDto{
				m :=  model.TestCreateChatDto(t)
				m.CompanionsUsernames[1], m.CompanionsUsernames[2] = m.CompanionsUsernames[2], m.CompanionsUsernames[1]
				return m
			},
			isValid: false,
		},
	}

	service, tr, err := TestService(t)
	if err != nil{
		t.Fatal(err)
		return
	}
	defer tr(service.db)
	
	m := model.TestCreateChatDto(t)

	service.User.RegisterUser(model.User{Username: m.CompanionsUsernames[0]})
	service.User.RegisterUser(model.User{Username: m.CompanionsUsernames[1]})
	service.User.RegisterUser(model.User{Username: m.CompanionsUsernames[2]})
	for _,tc := range testCases{
		t.Run(tc.name, func(t *testing.T) {
			
			mod := tc.u().ToModel(m.CompanionsUsernames[0])
			err = service.Chat.CreateChat(&mod)
			if tc.isValid{
				assert.NoError(t, err)
			} else{
				assert.Error(t, err)
			}
		})
	}
}

func TestGetChats(t *testing.T) {
	testCases := []struct {
		name     string
		s *Service
		u func() *model.User
		isValid bool
	}{
		{
			name:"valid",
			u: func() *model.User{
				return &model.TestUsers(t)[0]
			},
			isValid: true,
		},
	}

	service, tr, err := TestService(t)
	if err != nil{
		t.Fatal(err)
		return
	}
	defer tr(service.db)
	
	testChat := model.TestChat(t)
	users := model.TestUsers(t)
	for _, user := range users{
		_, res := service.User.RegisterUser(user)
		if res != nil{
			t.Fatal(err)
			return
		}
	}
	testChat.Users = users
	err = service.Chat.CreateChat(testChat)
	if err != nil{
		t.Fatal(err)
		return
	}
	// equal := []*model.Chat{testChat}

	for _,tc := range testCases{
		t.Run(tc.name, func(t *testing.T) {
			
			username := tc.u().Username
			_, err := service.Chat.GetChats(username)
			if tc.isValid{
				assert.NoError(t, err)
			} else{
				assert.Error(t, err)
			}
			// assert.Equal(t, equal, m, "test struct")
			// log.Printf("%+v", m[0])

		})
	}

}