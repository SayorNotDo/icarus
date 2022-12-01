package resource

import (
	"time"
)

type Case struct {
	CID            uint64    `json:"c_id" gorm:"primatyKey; autoIncrement"`
	Name           string    `json:"name" gorm:"type:text"`
	Status         uint8     `json:"status" gorm:"default:1"`
	Creator        string    `json:"creator" gorm:"type:varchar(256)"`
	CreateTime     time.Time `json:"create_time" gorm:"autoCreateTime:milli"`
	LastUpdateTime time.Time `json:"last_update_time" gorm:"autoUpdateTime:milli"`
}

type CaseContainer struct {
	CID            uint64    `json:"c_id" gorm:"index: c_container"`
	ContainerID    uint64    `json:"container_id" gorm:"index: c_container"`
	Name           string    `json:"name" gorm:"type:text"`
	CreateTime     time.Time `json:"create_time" gorm:"autoCreateTime:milli"`
	LastUpdateTime time.Time `json:"last_update_time" gorm:"autoUpdateTime:milli"`
}

type CaseCollection struct {
	CollectionID uint64 `json:"collection_id" gorm:"primatyKey; autoIncrement"`
	Name         string `json:"name" gorm:"type:text"`
	Description  string `json:"description" gorm:"type:text"`
	Reference    string `json:"reference" gorm:"type:text"`
}

type CollectionElement struct {
	CollectionID uint64    `json:"collection_id" gorm:"index: c_collection"`
	CID          uint64    `json:"c_id" gorm:"index: c_collection"`
	CreateTime   time.Time `json:"create_time" gorm:"autoCreateTime:milli"`
}
