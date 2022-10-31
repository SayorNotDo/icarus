package user

import (
	"icaru/database"
)

func init() {
	// return current database
	database.Db.Migrator().CurrentDatabase()

	// Migrate: run auto migration for given models, will only add missing field, won't delete/change current data
	database.Db.AutoMigrate(&User{}, &Department{})
	database.Db.Migrator().RenameTable("users", "user")
	database.Db.Migrator().RenameTable("departments", "department")
}
