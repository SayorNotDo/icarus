package user

import (
	"errors"
	"fmt"
	"github.com/golang-jwt/jwt"
	"github.com/gomodule/redigo/redis"
	"github.com/kataras/iris/v12"
	"github.com/kataras/iris/v12/context"
	"golang.org/x/crypto/bcrypt"
	redismanage "icarus/database/redis"
	"log"
	"time"
)

// UserService will deal with `user` model CRUD operation
type UserService interface {
	Create(params map[string]string, hashedPassword []byte) Status
	Update(user *User, params map[string]interface{}) (User, error)
	Login(username, password string) (tokenMap context.Map, status Status)
	Logout(params map[string]interface{}) (err error)
	Authenticate(oldToken, oldRefreshToken string) (token, refreshToken string, status Status)
	Authorize(username, password string) (authorizeToken string, status Status)
	GetUserInfo(ctx iris.Context) *UserInfo
	updateLoginTime(user *User) bool
	DeleteByID(uid uint32) bool
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

func (u *userService) GetUserInfo(ctx iris.Context) *UserInfo {
	uid, username := ParseUserinfo(ctx)
	user := u.repo.Select(&User{Username: username, UID: uid})
	userInfo := BuildUserInfo(user)
	return userInfo
}

func (u *userService) Create(params map[string]string, hashedPassword []byte) Status {
	for key, val := range params {
		if queryResult := u.repo.QueryByField(key, val); queryResult.UID != 0 {
			return Status{
				err:        errors.New(fmt.Sprintf("%s is already in used", key)),
				statusCode: iris.StatusBadRequest,
			}
		}
	}
	user := User{
		Username:       params["username"],
		Email:          params["email"],
		Phone:          params["phone"],
		HashedPassword: hashedPassword,
	}
	if err := u.repo.Insert(&user); err != nil {
		return Status{
			err:        err,
			statusCode: iris.StatusInternalServerError,
		}
	}
	return Status{
		err:        nil,
		statusCode: iris.StatusCreated,
	}
}

func (u *userService) Update(user *User, params map[string]interface{}) (User, error) {
	if user := u.repo.Select(user); user == nil {
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

func (u *userService) Login(username, password string) (tokenMap context.Map, status Status) {
	user := u.repo.Select(&User{Username: username})
	if user == nil {
		return nil, Status{
			err:        errors.New("user is not exist"),
			statusCode: iris.StatusBadRequest,
		}
	}
	if ok := validatePassword(password, user.HashedPassword); !ok {
		return nil, Status{
			err:        errors.New("username or password is wrong"),
			statusCode: iris.StatusBadRequest,
		}
	}
	tokenMap, err := generateToken(user.Username, user.UID)
	if err != nil {
		return nil, Status{
			err:        err,
			statusCode: iris.StatusInternalServerError,
		}
	}
	if err := u.repo.Updates(user, map[string]interface{}{"refresh_token": tokenMap["refreshToken"]}); err != nil {
		return nil, Status{
			err:        err,
			statusCode: iris.StatusInternalServerError,
		}
	}
	u.updateLoginTime(user)
	return tokenMap, Status{
		err:        nil,
		statusCode: iris.StatusOK,
	}
}

func (u *userService) updateLoginTime(user *User) bool {
	currentTime := time.Now()
	if err := u.repo.Updates(user, map[string]interface{}{"last_login_time": currentTime}); err != nil {
		return false
	}
	return true
}

func (u *userService) Logout(params map[string]interface{}) (err error) {
	log.Println("implement me")
	user := &User{Username: params["username"].(string)}
	err = u.repo.Updates(user, params)
	return
}

func (u *userService) Authenticate(oldAccessToken, oldRefreshToken string) (accessToken, refreshToken string, status Status) {
	parseToken, err := jwt.Parse(oldAccessToken, func(t *jwt.Token) (interface{}, error) { return secret, nil })
	if err != nil && err.Error() != "Token is expired" {
		return "", "", Status{
			err:        err,
			statusCode: iris.StatusUnauthorized,
		}
	}
	userClaims := parseToken.Claims.(jwt.MapClaims)
	conn := redismanage.Pool.Get()
	defer func(conn redis.Conn) {
		if err := conn.Close(); err != nil {
			log.Println("get user credentials failed: ", err)
		}
	}(conn)
	result, err := redismanage.RedisString(conn.Do("GET", fmt.Sprintf("%s:%s", userClaims["username"], "refreshToken")))
	if err != nil || oldRefreshToken != result {
		return "", "", Status{
			err:        errors.New("refreshToken is invalid"),
			statusCode: iris.StatusUnauthorized,
		}
	}
	user := u.repo.Select(&User{Username: userClaims["username"].(string), RefreshToken: result})
	if user == nil {
		return "", "", Status{
			err:        errors.New("authenticate error"),
			statusCode: iris.StatusInternalServerError,
		}
	}
	tokenMap, err := generateToken(user.Username, user.UID)
	if err != nil {
		log.Println("error generate accessToken")
		return "", "", Status{
			err:        errors.New("accessToken generate error"),
			statusCode: iris.StatusInternalServerError,
		}
	}
	if err := u.repo.Updates(user, map[string]interface{}{"refresh_token": tokenMap["refreshToken"]}); err != nil {
		return "", "", Status{
			err:        err,
			statusCode: iris.StatusInternalServerError,
		}
	}
	return accessToken, oldRefreshToken, Status{
		err:        nil,
		statusCode: iris.StatusOK,
	}
}

func (u *userService) Authorize(username, password string) (token string, status Status) {
	user := u.repo.Select(&User{Username: username})
	if user == nil {
		return "", Status{
			err:        errors.New("user is not exist"),
			statusCode: iris.StatusBadRequest,
		}
	}
	if ok := validatePassword(password, user.HashedPassword); !ok {
		return "", Status{
			err:        errors.New("username or password is wrong"),
			statusCode: iris.StatusBadRequest,
		}
	}
	token, err := authorizeToken(user.Username, user.UID)
	if err != nil {
		return "", Status{
			err:        err,
			statusCode: iris.StatusInternalServerError,
		}
	}
	return token, Status{
		err:        nil,
		statusCode: iris.StatusOK,
	}
}

func (u *userService) DeleteByID(uid uint32) bool {
	return false
}

func generatePassword(password string) ([]byte, error) {
	return bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
}

func validatePassword(password string, hashed []byte) bool {
	if err := bcrypt.CompareHashAndPassword(hashed, []byte(password)); err != nil {
		return false
	}
	return true
}

func generateToken(username string, uid uint32) (token context.Map, err error) {
	tokenMap := make(context.Map)
	accessToken, err := generateAccessToken(username, uid)
	if err != nil {
		return tokenMap, err
	}
	tokenMap["accessToken"] = accessToken
	refreshToken, err := generateRefreshToken(accessToken)
	if err != nil {
		return tokenMap, err
	}
	tokenMap["refreshToken"] = refreshToken
	tokenMap["tokenType"] = "Bearer"
	conn := redismanage.Pool.Get()
	defer func(conn redis.Conn) {
		if err := conn.Close(); err != nil {
			log.Println("Redis client closed connection error")
		}
	}(conn)
	if _, err := conn.Do("SET", fmt.Sprintf("%s:%s", username, "refreshToken"), refreshToken, "EX", 60*60*24*7); err != nil {
		return tokenMap, err
	}
	return tokenMap, nil
}

func validateAdministrator(uid uint32) bool {
	return false
}
