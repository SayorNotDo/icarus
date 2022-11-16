package trino

import (
	"database/sql"
	"log"

	_ "github.com/trinodb/trino-go-client/trino"
)

var TrinoDb *sql.DB

func init() {
	dsn := ""
	TrinoDb, err := sql.Open("trino", dsn)
	if err != nil {
		log.Printf("connect error: %v", err)
		return
	}
	log.Println(TrinoDb)
}
