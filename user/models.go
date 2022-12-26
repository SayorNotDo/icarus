package user

import (
	"fmt"
	"time"
)

type User struct {
	UID            uint32    `json:"uid" gorm:"primaryKey; autoIncrement"`
	Username       string    `json:"username" gorm:"unique; not null; type:varchar(32)"`
	HashedPassword []byte    `json:"hashed_password"`
	ChineseName    string    `json:"chinese_name" gorm:"type:varchar(32)"`
	CreateTime     time.Time `json:"create_time" gorm:"autoCreateTime"`
	LastUpdateTime time.Time `json:"last_update_time" gorm:"autoUpdateTIme:milli"`
	RoleId         uint8     `json:"role_id" gorm:"default:5"`
	EmployeeId     string    `json:"employee_id" gorm:"type:varchar(32)"`
	Position       string    `json:"position" gorm:"type:varchar(256)"`
	Email          string    `json:"email" gorm:"type:varchar(256)" validate:"email"`
	Phone          string    `json:"phone" gorm:"type:varchar(32)"`
	JoinDate       time.Time `json:"join_date"`
	LeaveDate      time.Time `json:"leave_date"`
	LastLoginTime  time.Time `json:"last_login_time"`
	Status         bool      `json:"status" gorm:"default:true"`
	Department     string    `json:"department" gorm:"type:varchar(256)"`
	RefreshToken   string    `json:"refresh_token" gorm:"type:text"`
}

type Tabler interface {
	TableName() string
}

func (User) TableName() string {
	return "user"
}

type Department struct {
	DID     int16  `json:"d_id" gorm:"primaryKey"`
	Name    string `json:"name" gorm:"not null; type:varchar(256)"`
	Center  string `json:"center" gorm:"type:varchar(256)"`
	Company string `json:"company" gorm:"not null; type:varchar(256)"`
}

func (Department) TableName() string {
	return "department"
}

type Serializer struct {
	UID            uint32    `json:"uid"`
	JoinDate       time.Time `json:"join_date" gorm:"autoCreateTime"`
	LastLoginTime  time.Time `json:"last_login_time" gorm:"autoUpdateTime:milli"`
	Username       string    `json:"username"`
	HashedPassword []byte    `json:"-"`
}

func (u User) Serializer() Serializer {
	return Serializer{
		UID:            u.UID,
		Username:       u.Username,
		HashedPassword: u.HashedPassword,
		JoinDate:       u.JoinDate,
		LastLoginTime:  u.LastLoginTime,
	}
}

func userValidate(params map[string]interface{}) (err error) {
	for key, param := range params {
		switch key {
		case "chinese_name", "position", "email", "department", "employee_id", "phone":
			if _, ok := param.(string); !ok {
				return fmt.Errorf("filed %s must be string type", key)
			}
		case "role_id":
			if _, ok := param.(int8); !ok {
				return fmt.Errorf("filed %s must be int8 type", key)
			}
		case "status":
			if _, ok := param.(bool); !ok {
				return fmt.Errorf("field %s must be bool type", key)
			}
		}
	}
	return nil
}
