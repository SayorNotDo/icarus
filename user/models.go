package user

import (
	"golang.org/x/crypto/bcrypt"
	"time"
)

type User struct {
	UID            int64     `json:"uid" gorm:"primaryKey"`
	Username       string    `json:"username" gorm:"unique; not null"`
	HashedPassword []byte    `json:"-" gorm:"-"`
	ChineseName    string    `json:"chineseName"`
	RoleId         int8      `json:"roleId" gorm:"default:5"`
	EmployeeId     string    `json:"employeeId" gorm:"<-:false"`
	Position       string    `json:"position"`
	Email          string    `json:"email"`
	Phone          string    `json:"phone"`
	JoinDate       time.Time `json:"joinDate" gorm:"autoCreateTime"`
	LastLoginTime  time.Time `json:"lastLoginTime" gorm:"autoUpdateTime:milli"`
	Status         bool      `json:"status" gorm:"default:1"`
	Department     string    `json:"department"`
}

type Department struct {
	DID     int16  `json:"did" gorm:"primaryKey"`
	Name    string `json:"name" gorm:"not null"`
	Center  string `json:"center"`
	Company string `json:"company" gorm:"not null"`
}

func GeneratePassword(password string) ([]byte, error) {
	return bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
}

func ValidatePassword(password string, hashed []byte) (bool, error) {
	if err := bcrypt.CompareHashAndPassword(hashed, []byte(password)); err != nil {
		return false, err
	}
	return true, nil
}
