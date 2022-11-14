package user

import (
	"errors"
	"log"
)

// UserService will deal with `user` model CRUD operation
type UserService interface {
	Create(params map[string]string) (User, error)
	Update(params map[string]string) (User, error)
	Login(username, password string) (token string, err error)
	//GetAll() []User
	//GetByID(uid int64) (User, bool)
	//DeleteByID(uid int64) bool
}

func NewUserService(repo UserRepository) UserService {
	return &userService{
		repo: repo,
	}
}

type userService struct {
	repo UserRepository
}

// Create insert a new user
// the password is the client-typed password
// it will be hashed before the insertion to our repository.
func (u *userService) Create(params map[string]string) (User, error) {
	username := params["username"]
	password := params["password"]
	email := params["email"]
	phone := params["phone"]
	log.Printf("username: %s", username)
	if password == "" || username == "" {
		return User{}, errors.New("username or password is empty")
	}
	log.Println("validate if the user is already registered.")
	user := User{
		Username: username,
		Email:    email,
		Phone:    phone,
	}
	_, found := u.repo.Select(user)
	log.Printf("validate result: %v", found)
	if found {
		return User{}, errors.New("user is already exist")
	}
	hashed, err := GeneratePassword(password)
	log.Printf("Get hashed password: %s", hashed)
	if err != nil {
		log.Println("generate password error!")
		return User{}, err
	}
	user.HashedPassword = hashed
	log.Printf("user's Info: %v", user)
	return u.repo.Insert(user)
}

func (u *userService) Update(params map[string]string) (User, error) {
	log.Println("try to update user info")
	chineseName := params["chinese_name"]
	employeeId := params["employee_id"]
	position := params["position"]
	department := params["department"]
	phone := params["phone"]
	log.Printf("chinese_name: %s, employee_id: %s, position: %s, department: %s, phone: %s", chineseName, employeeId, position, department, phone)
	return User{}, nil
}

func (u *userService) Login(username, password string) (token string, err error) {
	user, found := u.repo.Select(User{Username: username})
	if !found {
		return "", errors.New("username or password is wrong")
	}
	res, _ := ValidatePassword(password, user.HashedPassword)
	if res {
		token, err := generateToken(user.Username, user.UID)
		return token, err
	}
	return "", errors.New("username or password is wrong")
}

//func (u *userService) GetAll() []User {
//	return u.repo.SelectMany(func(_ User) bool {
//		return true
//	}, -1)
//}
//
//func (u *userService) GetByID(uid int64) (User, bool) {
//	return u.repo.Select(func(s User) bool {
//		return s.UID == uid
//	})
//}
//
//func (u userService) DeleteByID(uid int64) bool {
//	return u.repo.Delete(func(s User) bool {
//		return s.UID == uid
//	}, 1)
//}
