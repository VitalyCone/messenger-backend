package model_test

import (
	"testing"

	"github.com/VitalyCone/websocket-messenger/internal/app/model"
	"github.com/stretchr/testify/assert"
)

func TestCreateUserDroValidate(t *testing.T) {
	testCases := []struct {
		name     string
		u func() *model.CreateUserDto
		isValid bool
	}{
		{
			name:"valid",
			u: func() *model.CreateUserDto{
				return model.TestCreateUserDto(t)
			},
			isValid: true,
		},
		{
			name:"don't have username",
			u: func() *model.CreateUserDto{
				m := model.TestCreateUserDto(t)
				m.Username = ""
				return m
			},
			isValid: false,
		},
		{
			name:"don't have password",
			u: func() *model.CreateUserDto{
				m := model.TestCreateUserDto(t)
				m.Password = ""
				return m
			},
			isValid: false,
		},
		{
			name:"username dont alphanum",
			u: func() *model.CreateUserDto{
				m := model.TestCreateUserDto(t)
				m.Username = "wsada241$%#}"
				return m
			},
			isValid: false,
		},
		{
			name:"username min 3",
			u: func() *model.CreateUserDto{
				m := model.TestCreateUserDto(t)
				m.Username = "da"
				return m
			},
			isValid: false,
		},
		{
			name:"username max 32",
			u: func() *model.CreateUserDto{
				m := model.TestCreateUserDto(t)
				m.Username = "hhhhhhhhhhhhhhhhhhhhhhhhhhhhhhhhhhhhhhhhhhhhhhhhhhhhhhhhhhhhhhhhhhhhh"
				return m
			},
			isValid: false,
		},
		{
			name:"password min 3",
			u: func() *model.CreateUserDto{
				m := model.TestCreateUserDto(t)
				m.Password = "da"
				return m
			},
			isValid: false,
		},
		{
			name:"password max 32",
			u: func() *model.CreateUserDto{
				m := model.TestCreateUserDto(t)
				m.Password = "hhhhhhhhhhhhhhhhhhhhhhhhhhhhhhhhhhhhhhhhhhhhhhhhhhhhhhhhhhhhhhhhhhhhh"
				return m
			},
			isValid: false,
		},

		{
			name:"firstname max 50",
			u: func() *model.CreateUserDto{
				m := model.TestCreateUserDto(t)
				m.FirstName = "hhhhhhhhhhhhhhhhhhhhhhhhhhhhhhhhhhhhhhhhhhhhhhhhhhhhhhhhhhhhhhhhhhhhh"
				return m
			},
			isValid: false,
		},
		{
			name:"firstname can be empty",
			u: func() *model.CreateUserDto{
				m := model.TestCreateUserDto(t)
				m.FirstName = ""
				return m
			},
			isValid: true,
		},

		{
			name:"secondname max 50",
			u: func() *model.CreateUserDto{
				m := model.TestCreateUserDto(t)
				m.SecondName = "hhhhhhhhhhhhhhhhhhhhhhhhhhhhhhhhhhhhhhhhhhhhhhhhhhhhhhhhhhhhhhhhhhhhh"
				return m
			},
			isValid: false,
		},
		{
			name:"secondname can be empty",
			u: func() *model.CreateUserDto{
				m := model.TestCreateUserDto(t)
				m.SecondName = ""
				return m
			},
			isValid: true,
		},
		{
			name:"role user",
			u: func() *model.CreateUserDto{
				m := model.TestCreateUserDto(t)
				m.Role = "user"
				return m
			},
			isValid: true,
		},
		{
			name:"role admin",
			u: func() *model.CreateUserDto{
				m := model.TestCreateUserDto(t)
				m.Role = "admin"
				return m
			},
			isValid: true,
		},
		{
			name:"role empty",
			u: func() *model.CreateUserDto{
				m := model.TestCreateUserDto(t)
				m.Role = ""
				return m
			},
			isValid: false,
		},
		{
			name:"role something else",
			u: func() *model.CreateUserDto{
				m := model.TestCreateUserDto(t)
				m.Role = "dasfafwfasfaw"
				return m
			},
			isValid: false,
		},
	}

	for _,tc := range testCases{
		t.Run(tc.name, func(t *testing.T) {
			_, err := tc.u().ToModel()
			if tc.isValid{
				assert.NoError(t, err)
			} else{
				assert.Error(t, err)
			}
		})
	}
}

