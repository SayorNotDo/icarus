package user

import (
	"icarus/database"
	"log"
	"sync"
)

// Query represent `Guest` and action of query.
type Query func(User) bool

// UserRepository will process some user instance operation
// interface of test
type UserRepository interface {
	InsertOrUpdate(user User) (updateUser User, err error)
	Select(user User) (selectUser User, found bool)
	//Exec(query Query, action Query, limit int, mode int) (ok bool)
	//SelectMany(query Query, limit int) (results []User)
	//Delete(query Query, limit int) (deleted bool)
}

func NewUserRepository() UserRepository {
	return &userRepository{}
}

type userRepository struct {
	mu sync.RWMutex
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
	log.Println("Query User info from database.")
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

func (r *userRepository) InsertOrUpdate(user User) (updateUser User, err error) {
	//uid := user.UID
	//if uid == 0 {
	//	var lastUID int64
	//	// find the biggest uid in order to avoid duplicated
	//	r.mu.RLock()
	//	for _, item := range r.source {
	//		if item.UID > lastUID {
	//			lastUID = item.UID
	//		}
	//	}
	//	r.mu.RUnlock()
	//	uid = lastUID + 1
	//	user.UID = uid
	//
	//	// map-specific thing
	//	r.mu.Lock()
	//	r.source[uid] = user
	//	r.mu.Unlock()
	//	return user, nil
	//}
	//_, exists := r.Select(func(user User) bool {
	//	return user.UID == uid
	//})
	//if !exists {
	//	return User{}, errors.New("failed to update a nonexistent user")
	//}
	//// or comment these and r.source[id] = user for pure replace
	////if user.Username != "" {
	////	current.Username = user.Username
	////}
	//// map specific thing
	//r.mu.Lock()
	//r.source[uid] = user
	//r.mu.Unlock()
	//return user, nil
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
	} else {
		log.Println("try to update user info")
		return User{}, nil
	}
}

func (r *userRepository) Delete(query Query, limit int) (deleted bool) {
	//TODO implement me
	panic("implement me")
}
