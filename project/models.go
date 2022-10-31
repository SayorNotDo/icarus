package project

import "time"

type Project struct {
	PID            int64     `json:"pid" form:"pid"`
	CreateTime     time.Time `json:"createTime" form:"lastUpdateTime"`
	LastUpdateTime time.Time `json:"lastUpdateTime" form:"lastUpdateTime"`
	Name           string    `json:"name" form:"name"`
	Manager        string    `json:"manager" form:"manager"`
	Description    string    `json:"description" form:"description"`
	Status         int8      `json:"status" form:"status"`
	StartTime      time.Time `json:"startTime" form:"startTime"`
	FinishTime     time.Time `json:"finishTime" form:"finishTime"`
	Reference      string    `json:"reference" form:"reference"`
}
