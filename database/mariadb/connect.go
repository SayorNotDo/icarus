package database

import (
	"fmt"
	"icarus/utils"
	"log"
	"os"
	"time"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var Db *gorm.DB

/*
Database Initialize
init() means func run while package initializing
*/
func init() {
	// ------------------------------logger setting------------------------------
	newLogger := logger.New(
		log.New(os.Stdout, "\r\n", log.LstdFlags), // io writer
		logger.Config{
			SlowThreshold:             time.Second,   //Slow SQL threshold
			LogLevel:                  logger.Silent, //Log level
			IgnoreRecordNotFoundError: true,          //Ignore ErrRecordNotFound error for logger
			Colorful:                  false,         //Disable color
		},
	)
	// ------------------------------database connection------------------------------
	var err error
	dsn := fmt.Sprintf("%s:%s@tcp(%s:3306)/%s?parseTime=true&charset=utf8mb4&collation=utf8mb4_unicode_ci",
		utils.GetEnv("MYSQL_USER", "root"),
		utils.GetEnv("MYSQL_PASSWORD", "test1234"),
		utils.GetEnv("MYSQL_HOST", "localhost"),
		utils.GetEnv("MYSQL_DATABASE", "icarus"),
	)
	Db, err = gorm.Open(mysql.New(mysql.Config{
		DSN:                       dsn,
		DefaultStringSize:         256,
		DisableDatetimePrecision:  true,
		DontSupportRenameIndex:    true,
		DontSupportRenameColumn:   true,
		SkipInitializeWithVersion: false,
	}), &gorm.Config{
		Logger: newLogger,
	})
	if err != nil {
		log.Fatalf("error occur while connecting to the MySQL database: %v", err)
	}
	if Db.Error != nil {
		log.Fatalf("Database error %v", Db.Error)
	}
	// ------------------------------connection poll setting------------------------------
	sqlDB, err := Db.DB()

	if err != nil {
		log.Fatalf("Database error %v", err.Error())
	}

	// SetMaxIdleConns: maximum of idle connection pool
	sqlDB.SetMaxIdleConns(10)

	// SetMaxOpenConns: maximum of database connection opening
	sqlDB.SetMaxOpenConns(100)

	// SerConnMaxLifetime: max time of reusable connection
	sqlDB.SetConnMaxLifetime(time.Hour)

	log.Print("----------------------------connect success----------------------------")
}
