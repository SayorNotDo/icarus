package database

import "icaru/user"

func init() {
	// return current database
	Db.Migrator().CurrentDatabase()

	// Migrate: run auto migration for given models, will only add missing field, won't delete/change current data
	Db.AutoMigrate(&user.User{})
}
