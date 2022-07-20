package db

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"sync"

	_ "github.com/go-sql-driver/mysql"
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
	dbPort := "3306"
	dbUser := "root"
	dbPass := "password123!"
	dbName := "belajar_golang_api"

	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s", dbUser, dbPass, dbHost, dbPort, dbName)

	db, err := sql.Open("mysql", dsn)

	if err != nil {
		log.Default().Println(err)
		os.Exit(1)
		db.Close()
	}

	db.SetMaxIdleConns(8)

	return db
}

func DB() *sql.DB {
	return dbInstance.Instance().(*sql.DB)
}
