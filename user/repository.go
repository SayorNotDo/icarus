package user

import (
	"errors"
	"fmt"
	database "icarus/database/mariadb"
	"log"
)

// UserRepository will process some user instance operation
// interface of test
type UserRepository interface {
	Insert(user User) (insertUser User, err error)
	Select(user User) (selectUser User, found bool)
	Updates(user User, updateInfo map[string]interface{}) (err error)
	Delete(uid uint32) (deleted bool)
	QueryByField(filed string, value interface{}) (user User)
	//Exec(query Query, action Query, limit int, mode int) (ok bool)
}

func NewUserRepository() UserRepository {
	return &userRepository{}
}

type userRepository struct {
	// mu sync.RWMutex
}

const (
	ReadOnlyMode = iota
	ReadWriteMode
)

func (r *userRepository) QueryByField(field string, value interface{}) (user User) {
	database.Db.Model(&User{}).Where(fmt.Sprintf("%s = '%s'", field, value)).Find(&user)
	return
}

func (r *userRepository) Delete(uid uint32) (deleted bool) {
	log.Println(uid)
	return false
}

func (r *userRepository) Select(user User) (u User, found bool) {
	result := database.Db.Model(&User{}).Where(&user).First(&u)
	if result.Error == nil {
		found = true
	}
	return
}

func (r *userRepository) Insert(user User) (insertUser User, err error) {
	uid := user.UID
	// validate if the user is already registered
	if uid == 0 {
		database.Db.Model(&User{}).Create(&user)
		return user, nil
	}
	return User{}, errors.New("user create failed")
}

func (r *userRepository) Updates(user User, updateInfo map[string]interface{}) (err error) {
	tx := database.Db.Model(&User{}).Where(&user).Updates(updateInfo)
	log.Println("update user")
	if tx.Error != nil {
		return tx.Error
	}
	return
}
