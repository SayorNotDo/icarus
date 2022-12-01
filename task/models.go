package task

import "time"

type Task struct {
	TID            uint32    `json:"t_id" gorm:"primaryKey; autoIncrement"`
	TPID           uint16    `json:"tp_id"`
	Name           string    `json:"name" gorm:"unique; not null; type:varchar(256)"`
	CreateTime     time.Time `json:"create_time" gorm:"autoCreateTime"`
	LastUpdateTime time.Time `json:"last_update_time" gorm:"autoUpdateTIme:milli"`
	Executor       string    `json:"executor"`
	StartTime      time.Time `json:"start_time"`
	FinishTime     time.Time `json:"finish_time"`
	Status         uint8     `json:"status" gorm:"default:1"`
}

type TaskContent struct {
	TCID        uint32 `json:"tc_id" gorm:"index: t_content"`
	TID         uint32 `json:"t_id" gorm:"index:t_content"`
	CID         uint32 `json:"c_id"`
	Description string `json:"description" gorm:"type:text"`
	Reference   string `json:"reference" gorm:"type:text"`
}

type Tabler interface {
	TableName() string
}

func (TaskContent) TableName() string {
	return "task_content"
}

func (Task) TableName() string {
	return "task"
}
