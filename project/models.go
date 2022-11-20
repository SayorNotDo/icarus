package project

import (
	"errors"
	"icarus/utils"
	"log"
	"reflect"
)

type Project struct {
	PID            uint16           `json:"p_id" gorm:"primaryKey; autoIncrement"`
	CreateTime     *utils.LocalTime `json:"create_time" gorm:"autoCreateTime"`
	LastUpdateTime *utils.LocalTime `json:"last_update_time" gorm:"autoUpdateTime:milli"`
	Name           string           `json:"name" gorm:"unique; not null; type:varchar(256)"`
	Designation    string           `json:"designation" gorm:"unique; type:varchar(256)"`
	Rank           uint8            `json:"rank" gorm:"default:1"`
	Description    string           `json:"description" gorm:"type:text"`
	Status         uint8            `json:"status" gorm:"default:1"`
	Reference      string           `json:"reference" gorm:"type:text"`
	StartTime      *utils.LocalTime `json:"start_time"`
	FinishTime     *utils.LocalTime `json:"finish_time"`
}

type ProjectMember struct {
	PID            uint16           `json:"p_id" gorm:"index:pro_member"`
	UID            uint32           `json:"uid" gorm:"index:pro_member"`
	CreateTime     *utils.LocalTime `json:"create_time" gorm:"autoCreateTime"`
	LastUpdateTime *utils.LocalTime `json:"last_update_time" gorm:"autoUpdateTime:milli"`
	Character      string           `json:"character"`
	JoinDate       *utils.LocalTime `json:"join_date"`
	LeaveDate      *utils.LocalTime `json:"leave_date"`
	Status         bool             `json:"status" gorm:"default:1"`
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

func projectValidate(params map[string]interface{}) (err error) {
	for key, param := range params {
		log.Printf("key: %v param type: %v", key, reflect.TypeOf(param))
		switch key {
		case "designation", "name", "description", "reference":
			if _, ok := param.(string); !ok {
				return errors.New("parameter type error")
			}
		default:
			delete(params, key)
		}
	}
	return nil
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
