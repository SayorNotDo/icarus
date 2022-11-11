package user

import (
	"errors"
	"log"
)

// UserService will deal with `user` model CRUD operation
type UserService interface {
	Create(password string, user User) (User, error)
	Update(user User) (User, error)
	Login(username, password string) (token []byte, err error)
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
func (u *userService) Create(password string, user User) (User, error) {
	log.Printf("user: %v", user)
	log.Printf("uid: %d, password: %s, username: %s", user.UID, password, user.Username)
	if user.UID > 0 || password == "" || user.Username == "" {
		return User{}, errors.New("unable to create this user")
	}
	log.Println("validate if the user is already registered.")
	_, found := u.repo.Select(user)
	log.Printf("validate result: %v", found)
	if found == true {
		return User{}, errors.New("user already exist")
	}
	hashed, err := GeneratePassword(password)
	log.Printf("Get hashed password: %s", hashed)
	if err != nil {
		log.Println("generate password error!")
		return User{}, err
	}
	user.HashedPassword = hashed
	log.Printf("user's Info: %v", user)
	return u.repo.InsertOrUpdate(user)
}

func (u *userService) Update(user User) (User, error) {
	//TODO implement me
	panic("implement me")
}

func (u *userService) Login(username, password string) (token []byte, err error) {
	user, found := u.repo.Select(User{Username: username})
	if found != true {
		return nil, errors.New("user is not exist")
	}
	res, _ := ValidatePassword(password, user.HashedPassword)
	if res == true {
		token, err := generateToken(Signer, user.Username, user.UID)
		log.Printf("generate func result: %v, %v", token, err)
		return token, err
	}
	return nil, errors.New("username or password wrong")
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
