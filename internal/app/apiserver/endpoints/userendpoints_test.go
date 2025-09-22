package endpoints

// import (
// 	"bytes"
// 	"net/http"
// 	"net/http/httptest"
// 	"testing"

// 	"github.com/VitalyCone/websocket-messenger/internal/app/apiserver/chat"
// 	"github.com/VitalyCone/websocket-messenger/internal/app/model"
// 	"github.com/VitalyCone/websocket-messenger/internal/app/service"
// 	mock_service "github.com/VitalyCone/websocket-messenger/internal/app/service/mocks"
// 	"github.com/gin-gonic/gin"
// 	"github.com/golang/mock/gomock"
// 	"github.com/redis/go-redis/v9"
// 	"github.com/stretchr/testify/assert"

// 	// "github.com/stretchr/testify/mock"
// 	"gorm.io/driver/postgres"
// 	"gorm.io/gorm"
// 	// "github.com/VitalyCone/websocket-messenger/internal/app/model"
// 	// "github.com/stretchr/testify/assert"
// )

// func TestUserEndpoints_RegisterUser(t *testing.T){
// 	type mockBehavior func(r *mock_service.MockUser, user model.User)

// 	tests := []struct {
// 		name                 string
// 		inputBody            string
// 		inputUser            model.User
// 		mockBehavior         mockBehavior
// 		expectedStatusCode   int
// 		expectedResponseBody string
// 	}{
// 		{
// 			name:      "Ok",
// 			inputBody: `{"username": "username", "first_name": "Test Name", "password": "qwerty", "role":"user"}`,
// 			inputUser: model.User{
// 				Username: "username",
// 				FirstName:     "Test Name",
// 				Password: "qwerty",
// 			},
// 			mockBehavior: func(r *mock_service.MockUser, user model.User) {
// 				r.EXPECT().RegisterUser(user).Return(gomock.Any().String(), nil)
// 			},
// 			expectedStatusCode:   http.StatusCreated,
// 		},
// 		{
// 			name:      "Wrong Input",
// 			inputBody: `{"username": "username"}`,
// 			inputUser: model.User{},
// 			mockBehavior: func(r *mock_service.MockUser, user model.User) {},
// 			expectedStatusCode:   http.StatusBadRequest,
// 		},
// 		{
// 			name:      "Service Error",
// 			inputBody: `{"username": "username", "first_name": "Test Name", "password": "qw"}`,
// 			inputUser: model.User{
// 				Username: "username",
// 				FirstName:     "Test Name",
// 				Password: "qw",
// 			},
// 			mockBehavior: func(r *mock_service.MockUser, user model.User) {
// 				r.EXPECT().RegisterUser(user).Return(gomock.All(), nil)
// 			},
// 			expectedStatusCode:   http.StatusBadRequest,
// 		},
// 	}

// 	for _, test := range tests {
// 		t.Run(test.name, func(t *testing.T) {
// 			// Init Dependencies
// 			c := gomock.NewController(t)
// 			defer c.Finish()
			
// 			repo := mock_service.NewMockUser(c)
// 			test.mockBehavior(repo, test.inputUser)
			
// 			db, err := gorm.Open(postgres.Open("host=localhost port=5432 user=postgres dbname=websocket_messenger_test password=postgres sslmode=disable"), &gorm.Config{})
// 			if err != nil {
// 				t.Fatal(err)
// 			}
			
// 			if err := db.AutoMigrate(model.User{}); err != nil{
// 				t.Fatal(err)
// 			}
// 			defer db.Migrator().DropTable(model.User{})

// 			rdb := redis.NewClient(
// 				&redis.Options{
// 					Addr:     "127.0.0.1:7080",
// 					Password: "", // no password set
// 					DB:       0,  // use default DB
// 				})

// 			service := service.NewService(db, rdb, "GNd0c64~Q?naTXb}p1{V|lbh&~`#`&@&")
// 			r := gin.New()

// 			hub := chat.NewHub()
// 			go hub.Run()

// 			ep := NewEndpoints(service, r)

// 			r.POST("/account/register", ep.RegisterUser)

// 			// Create Request
// 			w := httptest.NewRecorder()
// 			req:= httptest.NewRequest("POST", "/account/register",
// 				bytes.NewBufferString(test.inputBody))

// 			// Make Request
// 			r.ServeHTTP(w, req)

// 			// Assert
// 			assert.Equal(t, w.Code, test.expectedStatusCode)
// 		})
// 	}
// }