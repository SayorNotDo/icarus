package repository

import (
	database "icarus/database/mariadb"
)

type BaseRepository struct {
}

func (e *BaseRepository) Insert(tableName string, insertRecord interface{}) error {
	tx := database.Db.Table(tableName).Create(insertRecord)
	if err := tx.Error; err != nil || tx.RowsAffected == 0 {
		return err
	}
	return nil
}

func (e *BaseRepository) Select(tableName string, query interface{}, model interface{}) error {
	tx := database.Db.Table(tableName).Where(query).Take(&model)
	if err := tx.Error; err != nil {
		return err
	}
	return nil
}
