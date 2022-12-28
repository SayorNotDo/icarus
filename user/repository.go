package user

import (
	"errors"
	"fmt"
	database "icarus/database/mariadb"
	. "icarus/repository"
	"log"
)

// UserRepository will process some user instance operation
type UserRepository interface {
	Insert(user *User) error
	Select(user *User) *User
	Updates(user *User, updateInfo map[string]interface{}) error
	Delete(uid uint32) (deleted bool)
	QueryByField(filed string, value interface{}) *User
}

func NewUserRepository() UserRepository {
	return &userRepository{}
}

type userRepository struct {
	base *BaseRepository
}

func (r *userRepository) Delete(uid uint32) (deleted bool) {
	//TODO implement me
	panic("implement me")
}

func (r *userRepository) QueryByField(field string, value interface{}) *User {
	user := &User{}
	database.Db.Model(&User{}).Where(fmt.Sprintf("%s = '%s'", field, value)).Find(&user)
	return user
}

func (r *userRepository) Select(user *User) *User {
	ret := &User{}
	if err := r.base.Select("user", &user, ret); err != nil {
		return nil
	}
	return ret
}

func (r *userRepository) Insert(user *User) error {
	uid := user.UID
	if uid != 0 {
		return errors.New("user create failed")
	}
	if err := r.base.Insert("user", user); err != nil {
		return err
	}
	return nil
}

func (r *userRepository) Updates(user *User, updateInfo map[string]interface{}) (err error) {
	tx := database.Db.Model(&User{}).Where(user).Updates(updateInfo)
	log.Println("update user")
	if tx.Error != nil {
		return tx.Error
	}
	return
}
