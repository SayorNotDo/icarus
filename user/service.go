package user

import (
	"errors"
	"fmt"
	"github.com/kataras/iris/v12"
	redismanage "icarus/database/redis"
	"log"

	"github.com/golang-jwt/jwt"
)

// UserService will deal with `user` model CRUD operation
type UserService interface {
	Create(params map[string]string) (User, Status)
	Update(user User, params map[string]interface{}) (User, error)
	Login(username, password string) (token, refreshToken string, err error)
	Logout(params map[string]interface{}) (err error)
	Authenticate(oldToken, oldRefreshToken string) (token, refreshToken string, err error)
	DeleteByID(uid uint32) bool
	//GetAll() []User
	//GetByID(uid int64) (User, bool)
}

func NewUserService(repo UserRepository) UserService {
	return &userService{
		repo: repo,
	}
}

type userService struct {
	repo UserRepository
}

type Status struct {
	err        error
	statusCode int16
}

// Create insert a new user
// the password is the client-typed password
// it will be hashed before the insertion to our repository.
func (u *userService) Create(params map[string]string) (User, Status) {
	username := params["username"]
	password := params["password"]
	email := params["email"]
	phone := params["phone"]
	if password == "" || username == "" {
		return User{},
			Status{
				err:        errors.New("username or password is empty"),
				statusCode: iris.StatusBadRequest,
			}
	}
	for key, val := range params {
		if key == "password" {
			continue
		}
		if queryResult := u.repo.Query(key, val); queryResult.UID != 0 {
			return User{}, Status{
				err:        errors.New(fmt.Sprintf("%s is already in used", key)),
				statusCode: iris.StatusBadRequest,
			}
		}
	}
	hashed, err := GeneratePassword(password)
	if err != nil {
		return User{}, Status{
			err:        err,
			statusCode: iris.StatusInternalServerError,
		}
	}
	user := User{
		Username:       username,
		Email:          email,
		Phone:          phone,
		HashedPassword: hashed,
	}
	res, err := u.repo.Insert(user)
	if err != nil {
		return User{}, Status{
			err:        err,
			statusCode: iris.StatusInternalServerError,
		}
	}
	return res, Status{
		err:        err,
		statusCode: iris.StatusCreated,
	}
}

func (u *userService) Update(user User, params map[string]interface{}) (User, error) {
	if _, found := u.repo.Select(user); !found {
		return User{}, errors.New("user is not exist")
	}
	if err := userValidate(params); err != nil {
		return User{}, err
	}
	if err := u.repo.Updates(user, params); err != nil {
		return User{}, err
	}
	return User{}, nil
}

func (u *userService) Login(username, password string) (token, refreshToken string, err error) {
	user, found := u.repo.Select(User{Username: username})
	if !found {
		return "", "", errors.New("username or password is wrong")
	}
	res, _ := ValidatePassword(password, user.HashedPassword)
	if res {
		token, err := generateAccessToken(user.Username, user.UID)
		if err == nil {
			if refreshToken, err := generateRefreshToken(token); err == nil {
				// refresh token store in redis
				conn := redismanage.Pool.Get()
				defer conn.Close()
				// conn.Do("SET", fmt.Sprintf("%s:%s", username, "accesstoken"), token, "EX", 60*5)
				conn.Do("SET", fmt.Sprintf("%s:%s", username, "refreshToken"), refreshToken, "EX", 60*60*24*7)
				// update refreshToken in mariadb
				u.repo.Updates(user, map[string]interface{}{"refresh_token": refreshToken})
				return token, refreshToken, err
			}
		}
	}
	return "", "", errors.New("username or password is wrong")
}

func (u *userService) Logout(params map[string]interface{}) (err error) {
	log.Println("implement me")
	user := User{Username: params["username"].(string)}
	err = u.repo.Updates(user, params)
	return
}

func (u *userService) Authenticate(oldToken, oldRefreshToken string) (token, refreshToken string, err error) {
	parseToken, err := jwt.Parse(oldToken, func(t *jwt.Token) (interface{}, error) { return secret, nil })
	if err != nil && err.Error() != "Token is expired" {
		return "", "", err
	}
	userClaims := parseToken.Claims.(jwt.MapClaims)
	conn := redismanage.Pool.Get()
	defer conn.Close()
	result, err := redismanage.RedisString(conn.Do("GET", fmt.Sprintf("%s:%s", userClaims["username"], "refreshToken")))
	if err != nil || oldRefreshToken != result {
		return "", "", errors.New("refreshToken is invalid")
	}
	selectUser, found := u.repo.Select(User{Username: userClaims["username"].(string), RefreshToken: result})
	if found {
		token, err = generateAccessToken(selectUser.Username, selectUser.UID)
		if err != nil {
			log.Println("Error generateAccessToken...")
			return "", "", errors.New("generate token error")
		}
		return token, oldRefreshToken, nil
	}
	return "", "", errors.New("authenticate error")
}

func (u *userService) DeleteByID(uid uint32) bool {
	log.Println("_________________________debug____________________")
	return false
}
