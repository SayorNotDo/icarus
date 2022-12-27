package user

import (
	"errors"
	"fmt"
	"github.com/golang-jwt/jwt"
	"github.com/gomodule/redigo/redis"
	"github.com/kataras/iris/v12"
	"golang.org/x/crypto/bcrypt"
	redismanage "icarus/database/redis"
	"log"
	"time"
)

// UserService will deal with `user` model CRUD operation
type UserService interface {
	Create(params map[string]string, hashedPassword []byte) Status
	Update(user User, params map[string]interface{}) (User, error)
	Login(username, password string) (token, refreshToken string, status Status)
	Logout(params map[string]interface{}) (err error)
	Authenticate(oldToken, oldRefreshToken string) (token, refreshToken string, status Status)
	Authorize(username, password string) (authorizeToken string, status Status)
	GetUserInfo(uid uint32, username string) (map[string]interface{}, bool)
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

func (u *userService) GetUserInfo(uid uint32, username string) (map[string]interface{}, bool) {
	user, found := u.repo.Select(User{Username: username, UID: uid})
	if !found {
		return map[string]interface{}{}, false
	}
	userInfo := map[string]interface{}{
		"Username":    user.Username,
		"Uid":         user.UID,
		"Email":       user.Email,
		"ChineseName": user.ChineseName,
		"RoleId":      user.RoleId,
		"EmployeeId":  user.EmployeeId,
		"JoinDate":    user.JoinDate,
		"Position":    user.Position,
		"Phone":       user.Phone,
		"Department":  user.Department,
	}
	return userInfo, true
}

// Create insert a new user
// the password is the client-typed password
// it will be hashed before the insertion to our repository.
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
	if _, err := u.repo.Insert(user); err != nil {
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

func (u *userService) Login(username, password string) (accessToken, refreshToken string, status Status) {
	user, found := u.repo.Select(User{Username: username})
	if !found {
		return "", "", Status{
			err:        errors.New("user is not exist"),
			statusCode: iris.StatusBadRequest,
		}
	}
	if ok, _ := validatePassword(password, user.HashedPassword); !ok {
		return "", "", Status{
			err:        errors.New("username or password is wrong"),
			statusCode: iris.StatusBadRequest,
		}
	}
	accessToken, refreshToken, err := generateToken(user.Username, user.UID)
	if err != nil {
		return "", "", Status{
			err:        err,
			statusCode: iris.StatusInternalServerError,
		}
	}
	lastLoginTime := time.Now()
	if err := u.repo.Updates(user, map[string]interface{}{"refresh_token": refreshToken, "last_login_time": lastLoginTime}); err != nil {
		return "", "", Status{
			err:        err,
			statusCode: iris.StatusInternalServerError,
		}
	}
	return accessToken, refreshToken, Status{
		err:        nil,
		statusCode: iris.StatusOK,
	}
}

// Logout TODO: implement completely
func (u *userService) Logout(params map[string]interface{}) (err error) {
	log.Println("implement me")
	user := User{Username: params["username"].(string)}
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
	user, found := u.repo.Select(User{Username: userClaims["username"].(string), RefreshToken: result})
	if !found {
		return "", "", Status{
			err:        errors.New("authenticate error"),
			statusCode: iris.StatusInternalServerError,
		}
	}
	accessToken, refreshToken, err = generateToken(user.Username, user.UID)
	if err != nil {
		log.Println("error generate accessToken")
		return "", "", Status{
			err:        errors.New("accessToken generate error"),
			statusCode: iris.StatusInternalServerError,
		}
	}
	if err := u.repo.Updates(user, map[string]interface{}{"refresh_token": refreshToken}); err != nil {
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
	user, found := u.repo.Select(User{Username: username})
	if !found {
		return "", Status{
			err:        errors.New("user is not exist"),
			statusCode: iris.StatusBadRequest,
		}
	}
	if ok, _ := validatePassword(password, user.HashedPassword); !ok {
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
	log.Println("_________________________debug____________________")
	return false
}

func generatePassword(password string) ([]byte, error) {
	return bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
}

func validatePassword(password string, hashed []byte) (bool, error) {
	log.Println("validate processing...")
	if err := bcrypt.CompareHashAndPassword(hashed, []byte(password)); err != nil {
		return false, err
	}
	return true, nil
}

func generateToken(username string, uid uint32) (accessToken, refreshToken string, err error) {
	accessToken, err = generateAccessToken(username, uid)
	if err != nil {
		return "", "", err
	}
	refreshToken, err = generateRefreshToken(accessToken)
	if err != nil {
		return "", "", err
	}
	conn := redismanage.Pool.Get()
	defer func(conn redis.Conn) {
		if err := conn.Close(); err != nil {
			log.Println("Redis client closed connection error")
		}
	}(conn)
	if _, err := conn.Do("SET", fmt.Sprintf("%s:%s", username, "refreshToken"), refreshToken, "EX", 60*60*24*7); err != nil {
		return "", "", err
	}
	return accessToken, refreshToken, nil
}

func validateAdministrator(uid uint32) bool {
	return false
}
