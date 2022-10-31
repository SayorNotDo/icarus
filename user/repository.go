package user

import (
	"errors"
	"log"
	"sync"
)

// Query represent `Guest` and action of query.
type Query func(User) bool

// UserRepository will process some user instance operation
// interface of test
type UserRepository interface {
	Exec(query Query, action Query, limit int, mode int) (ok bool)
	Select(query Query) (user User, found bool)
	SelectMany(query Query, limit int) (results []User)
	InsertOrUpdate(user User) (updateUser User, err error)
	Delete(query Query, limit int) (deleted bool)
}

func NewUserRepository(source map[int64]User) UserRepository {
	return &userMemoryRepository{source: source}
}

type userMemoryRepository struct {
	source map[int64]User
	mu     sync.RWMutex
}

const (
	ReadOnlyMode = iota
	ReadWriteMode
)

func (r *userMemoryRepository) Exec(query Query, action Query, actionLimit int, mode int) (ok bool) {
	loops := 0
	if mode == ReadOnlyMode {
		r.mu.RLock()
		defer r.mu.RUnlock()
	} else {
		r.mu.Lock()
		defer r.mu.Unlock()
	}

	for _, user := range r.source {
		ok = query(user)
		if ok {
			if action(user) {
				loops++
				if actionLimit >= loops {
					break
				}
			}
		}
	}
	return
}

func (r *userMemoryRepository) Select(query Query) (u User, found bool) {
	found = r.Exec(query, func(s User) bool {
		u = s
		return true
	}, 1, ReadOnlyMode)
	if !found {
		u = User{}
	}
	return
}

func (r *userMemoryRepository) SelectMany(query Query, limit int) (results []User) {
	r.Exec(query, func(s User) bool {
		results = append(results, s)
		return true
	}, limit, ReadOnlyMode)
	return
}

func (r *userMemoryRepository) InsertOrUpdate(user User) (updateUser User, err error) {
	uid := user.UID
	if uid == 0 {
		var lastUID int64
		// find the biggest uid in order to avoid duplicated
		r.mu.RLock()
		for _, item := range r.source {
			if item.UID > lastUID {
				lastUID = item.UID
			}
		}
		r.mu.RUnlock()
		uid = lastUID + 1
		user.UID = uid

		// map-specific thing
		r.mu.Lock()
		r.source[uid] = user
		r.mu.Unlock()
		log.Printf("create user: %v", user)
		return user, nil
	}
	_, exists := r.Select(func(user User) bool {
		return user.UID == uid
	})
	if !exists {
		return User{}, errors.New("failed to update a nonexistent user")
	}
	// or comment these and r.source[id] = user for pure replace
	//if user.Username != "" {
	//	current.Username = user.Username
	//}
	// map specific thing
	r.mu.Lock()
	r.source[uid] = user
	r.mu.Unlock()
	return user, nil
}

func (r *userMemoryRepository) Delete(query Query, limit int) (deleted bool) {
	//TODO implement me
	panic("implement me")
}
