package db

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"sync"

	_ "github.com/lib/pq"
)

func init() {
	dbInstance = &DBInstance{initializer: dbInit}
}

var dbInstance *DBInstance

type DBInstance struct {
	initializer func() interface{}
	instance    interface{}
	once        sync.Once
}

func (i *DBInstance) Instance() interface{} {
	i.once.Do(func() {
		i.instance = i.initializer()
	})
	return i.instance
}

func dbInit() interface{} {
	dbHost := "localhost"
	dbPort := "5432"
	dbUser := "admin"
	dbPass := "password"
	dbName := "sandbox_pii"

	dsn := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		dbHost, dbPort, dbUser, dbPass, dbName)

	db, err := sql.Open("postgres", dsn)

	if err != nil {
		log.Default().Println(err)
		os.Exit(1)
		db.Close()
	}

	if err := db.Ping(); err != nil {
		log.Default().Println(err)
		os.Exit(1)
	}

	db.SetMaxIdleConns(8)

	return db
}

func DB() *sql.DB {
	return dbInstance.Instance().(*sql.DB)
}
