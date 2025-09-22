package service

import (
	"context"
	"encoding/json"
	"errors"
	"time"

	"github.com/VitalyCone/websocket-messenger/internal/app/model"
	"github.com/dgrijalva/jwt-go"
	"github.com/redis/go-redis/v9"
	"github.com/sirupsen/logrus"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

const(
    bcryptSalt = 10
	tokenTTL   = 12 * time.Hour
)

type tokenClaims struct {
	jwt.StandardClaims
	UserUsername string `json:"user_username"`
}

type UserService struct {
	db *gorm.DB
	rdb *redis.Client
	tokenKey string
}

func NewUserService(db *gorm.DB, rdb *redis.Client, tokenKey string) *UserService {
	return &UserService{
		db: db,
		rdb: rdb,
		tokenKey: tokenKey,
	}
}

func (s *UserService) RegisterUser(m model.User) (interface{}, error){
	passHash, err := bcrypt.GenerateFromPassword([]byte(m.Password), bcryptSalt)
    if err != nil{
        return nil, err
    }
	m.PasswordHash = string(passHash)

	tx := s.db.Begin()

    resoult := tx.Create(&m)
	if resoult.Error != nil{
		tx.Rollback()
		return nil, resoult.Error
	}

	mByte, err := json.Marshal(m)
	if err != nil{
		return nil, err
	}

	err = s.rdb.Set(context.Background(), "user_" + m.Username, mByte, 0).Err()
    if err != nil {
		tx.Rollback()
        return nil, err
    }

    tokenString, err := createToken(m, s.tokenKey)
    if err != nil {
		tx.Rollback()
		return nil, err
    }

	tx.Commit()
	return tokenString, nil
}

func (s *UserService) LoginUser(m model.User) (interface{}, error){
	var user model.User
	val, err := s.rdb.Get(context.Background(), "user_" + m.Username).Result()
	if err != nil{
		resoult := s.db.Unscoped().Where(model.User{Username:m.Username}).First(&user)
		if resoult.Error != nil{
			return nil, resoult.Error
		}
		logrus.Printf("%s form db", user.Username)
	}else{
		if err := json.Unmarshal([]byte(val), &user); err != nil{
			return nil, err
		}
		logrus.Printf("%s form redis", user.Username)
	}
	if err := verifyPassword(user.PasswordHash, m.Password);err != nil{
		return nil, err
	}

	mByte, err := json.Marshal(user)
	if err != nil{
		return nil, err
	}

	err = s.rdb.Set(context.Background(), "user_" + m.Username, mByte, 0).Err()
    if err != nil {
        return nil, err
    }

	tokenString, err := createToken(m, s.tokenKey)
    if err != nil {
		return nil, err
    }

	return tokenString, nil
}

func (s *UserService) GetUsernameFromToken(tokenString string) (string, error){
	claims, err := verifyToken(tokenString, s.tokenKey)
	if err != nil{
		return  "", err
	}
	return claims.UserUsername, nil
}

func (s *UserService) GetUserData(tokenString string) (model.User, error){
	var user model.User

	claims, err := verifyToken(tokenString, s.tokenKey)
	if err != nil{
		return model.User{}, err
	}

	user, err = getUserByUsername(claims.UserUsername, s.db, s.rdb)
	if err != nil{
		return model.User{}, err
	}
	
	return user, nil
}

func (s *UserService) GetUsersWithQuery_ToResponse(username string, offset,limit int) ([]model.UserResponse, error){
	usersResp := make([]model.UserResponse, 0)
	users := make([]model.User, 0)
	respUsername := "%"+username+"%"

	resoult := s.db.Unscoped().Where("username LIKE ?", respUsername).Offset(offset).Limit(limit).Find(&users)
	if resoult.Error != nil{
		return nil, resoult.Error
	}
	for _, user := range users{
		usersResp = append(usersResp, user.ToResponse())
	}
	return usersResp, nil
}

func createToken(user model.User, tokenSigned string)(string, error){

    claims := jwt.NewWithClaims(jwt.SigningMethodHS256, &tokenClaims{
		jwt.StandardClaims{
			ExpiresAt: time.Now().Add(tokenTTL).Unix(),
			IssuedAt:  time.Now().Unix(),
		},
		user.Username,
	})

    return claims.SignedString([]byte(tokenSigned))
}

func verifyToken(tokenString string, tokenSignedString string) (*tokenClaims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &tokenClaims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("invalid signing method")
		}

		return []byte(tokenSignedString), nil
	})
	if err != nil {
		return nil, err
	}

	claims, ok := token.Claims.(*tokenClaims)
	if !ok {
		return nil, errors.New("token claims are not of type *tokenClaims")
	}

	return claims, nil
}

func verifyPassword(hashPass, password string) error{
    return bcrypt.CompareHashAndPassword([]byte(hashPass), []byte(password))
}

func getUserByUsername(username string, db *gorm.DB, rdb *redis.Client) (model.User, error){
	var user model.User
	val, err := rdb.Get(context.Background(), "user_" + username).Result()
	if err != nil{
		resoult := db.Unscoped().Where(model.User{Username:username}).First(&user)
		if resoult.Error != nil{
			return model.User{}, resoult.Error
		}
		logrus.Printf("%s form db", user.Username)
		mByte, err := json.Marshal(user)
		if err != nil{
			return model.User{}, err
		}
	
		err = rdb.Set(context.Background(), "user_" + username, mByte, 0).Err()
		if err != nil {
			return model.User{}, err
		}
	}else{
		if err := json.Unmarshal([]byte(val), &user); err != nil{
			return model.User{}, err
		}
		logrus.Printf("%s form redis", user.Username)
	}
	return user, nil
}	