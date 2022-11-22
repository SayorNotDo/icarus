package project

import (
	"errors"
	"reflect"
	"time"
)

type Project struct {
	PID            uint16    `json:"p_id" gorm:"primaryKey; autoIncrement"`
	CreateTime     time.Time `json:"create_time" gorm:"autoCreateTime:milli"`
	LastUpdateTime time.Time `json:"last_update_time" gorm:"autoUpdateTime:milli"`
	Name           string    `json:"name" gorm:"unique; not null; type:varchar(256)"`
	Designation    string    `json:"designation" gorm:"unique; type:varchar(256)"`
	Rank           uint8     `json:"rank" gorm:"default:1"`
	Description    string    `json:"description" gorm:"type:text"`
	Status         uint8     `json:"status" gorm:"default:1"`
	Reference      string    `json:"reference" gorm:"type:text"`
	StartTime      time.Time `json:"start_time"`
	FinishTime     time.Time `json:"finish_time"`
}

type ProjectMember struct {
	PID            uint16    `json:"p_id" gorm:"index:pro_member"`
	UID            uint32    `json:"uid" gorm:"index:pro_member"`
	CreateTime     time.Time `json:"create_time" gorm:"autoCreateTime:milli"`
	LastUpdateTime time.Time `json:"last_update_time" gorm:"autoUpdateTime:milli"`
	Character      string    `json:"character"`
	JoinDate       time.Time `json:"join_date"`
	LeaveDate      time.Time `json:"leave_date"`
	Status         bool      `json:"status" gorm:"default:1"`
}

type Tabler interface {
	TableName() string
}

func (Project) TableName() string {
	return "project"
}

func (ProjectMember) TableName() string {
	return "project_member"
}

func projectValidate(params map[string]interface{}) (map[string]interface{}, error) {
	for key, param := range params {
		switch key {
		case "designation", "name", "description", "reference":
			if reflect.TypeOf(param).Kind() != reflect.String {
				return nil, errors.New("parameter error")
			}
		default:
			delete(params, key)
		}
	}
	return params, nil
}

// func projectMemberValidate(params map[string]interface{}) (err error) {
// 	for key, param := range params {
// 		switch key {
// 		case "p_id", "uid":
// 			if _, ok := param.(int64); !ok {
// 				return fmt.Errorf("field %s must be int64 type", key)
// 			}
// 		case "character":
// 			if _, ok := param.(string); !ok {
// 				return fmt.Errorf("filed %s must be string type", key)
// 			}
// 		case "join_date", "leave_date":
// 			if _, ok := param.(*utils.LocalTime); !ok {
// 				return fmt.Errorf("filed %s must be time.Time type", key)
// 			}
// 		case "status":
// 			if _, ok := param.(bool); !ok {
// 				return fmt.Errorf("filed %s must be bool type", key)
// 			}
// 		}
// 	}
// 	return nil
// }
