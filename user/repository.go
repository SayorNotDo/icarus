package user

import (
	"icarus/database"
	"log"
)

// Query represent `Guest` and action of query.
type Query func(User) bool

// UserRepository will process some user instance operation
// interface of test
type UserRepository interface {
	Insert(user User) (updateUser User, err error)
	Select(user User) (selectUser User, found bool)
	//Exec(query Query, action Query, limit int, mode int) (ok bool)
	//SelectMany(query Query, limit int) (results []User)
	//Delete(query Query, limit int) (deleted bool)
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

//func (r *userMemoryRepository) Exec(query Query, action Query, actionLimit int, mode int) (ok bool) {
//	loops := 0
//	if mode == ReadOnlyMode {
//		r.mu.RLock()
//		defer r.mu.RUnlock()
//	} else {
//		r.mu.Lock()
//		defer r.mu.Unlock()
//	}
//
//	for _, user := range r.source {
//		ok = query(user)
//		if ok {
//			if action(user) {
//				loops++
//				if actionLimit >= loops {
//					break
//				}
//			}
//		}
//	}
//	return
//}

func (r *userRepository) Select(user User) (u User, found bool) {
	log.Println("Query User info from database")
	result := database.Db.Model(&User{}).Where("username = ?", user.Username).First(&u)
	log.Printf("Query Result: %v", result.Error)
	if result.Error == nil {
		found = true
	}
	return
}

//func (r *userMemoryRepository) SelectMany(query Query, limit int) (results []User) {
//	r.Exec(query, func(s User) bool {
//		results = append(results, s)
//		return true
//	}, limit, ReadOnlyMode)
//	return
//}

func (r *userRepository) Insert(user User) (updateUser User, err error) {
	uid := user.UID
	log.Printf("get user'UID: %d ", uid)
	// validate if the user is already registered
	if uid == 0 {
		log.Println("new user for create")
		// find the biggest UID in database avoid duplicated
		var lastUser User
		database.Db.Model(&User{}).First(&lastUser)
		user.UID = lastUser.UID + 1
		log.Printf("new user's UID: %d: ", user.UID)
		database.Db.Model(&User{}).Create(&user)
		log.Println("create user success!")
		return user, nil
	}
	return User{}, nil
}

func (r *userRepository) Delete(query Query, limit int) (deleted bool) {
	//TODO implement me
	panic("implement me")
}
