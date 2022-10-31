package user

import (
	"errors"
	"log"
)

// UserService will deal with `user` model CRUD operation
type UserService interface {
	GetAll() []User
	GetByID(uid int64) (User, bool)
	DeleteByID(uid int64) bool
	Create(password string, user User) (User, error)
}

func NewUserService(repo UserRepository) UserService {
	return &userService{
		repo: repo,
	}
}

type userService struct {
	repo UserRepository
}

// Create insert a new User
// the password is the client-typed password
// it will be hashed before the insertion to our repository.
func (u *userService) Create(password string, user User) (User, error) {
	log.Printf("user: %v", user)
	log.Printf("uid: %d, password: %s, username: %s", user.UID, password, user.Username)
	if user.UID > 0 || password == "" || user.Username == "" {
		return User{}, errors.New("unable to create this user")
	}
	hashed, err := GeneratePassword(password)
	if err != nil {
		return User{}, err
	}
	user.HashedPassword = hashed
	return u.repo.InsertOrUpdate(user)
}

func (u *userService) GetAll() []User {
	return u.repo.SelectMany(func(_ User) bool {
		return true
	}, -1)
}

func (u *userService) GetByID(uid int64) (User, bool) {
	return u.repo.Select(func(s User) bool {
		return s.UID == uid
	})
}

func (u userService) DeleteByID(uid int64) bool {
	return u.repo.Delete(func(s User) bool {
		return s.UID == uid
	}, 1)
}
