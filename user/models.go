package user

import (
	"golang.org/x/crypto/bcrypt"
	"time"
)

type User struct {
	UID            int64     `json:"uid" gorm:"primaryKey"`
	Username       string    `json:"username" gorm:"unique; not null; type:varchar(32)"`
	HashedPassword []byte    `json:"-" gorm:"-"`
	ChineseName    string    `json:"chineseName" gorm:"type:varchar(32)""`
	RoleId         int8      `json:"roleId" gorm:"default:5"`
	EmployeeId     string    `json:"employeeId" gorm:"<-:false; type:varchar(32)"`
	Position       string    `json:"position" gorm:"type:varchar(255)"`
	Email          string    `json:"email" gorm:"type:varchar(255)"`
	Phone          string    `json:"phone" gorm:"type:varchar(32)"`
	JoinDate       time.Time `json:"joinDate" gorm:"autoCreateTime"`
	LastLoginTime  time.Time `json:"lastLoginTime" gorm:"autoUpdateTime:milli"`
	Status         bool      `json:"status" gorm:"default:1"`
	Department     string    `json:"department" gorm:"type:varchar(255)"`
}

type Department struct {
	DID     int16  `json:"did" gorm:"primaryKey"`
	Name    string `json:"name" gorm:"not null; type:varchar(255)"`
	Center  string `json:"center" gorm:"type:varchar(255)"`
	Company string `json:"company" gorm:"not null; type:varchar(255)"`
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
