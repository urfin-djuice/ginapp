package db

import (
	"fmt"
	"log"
	"oko/pkg/env"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres" //nolint
)

const (
	tryPause = 5
)

var conn *gorm.DB

// Opening a database and save the reference to `Database` struct.
func Connect() {
	conString := env.GetEnvOrPanic("POSTGRES_URL")

	// open a db connection

	db, err := gorm.Open("postgres", conString)
	if err != nil {
		log.Println("db err: ", err)
		conn = nil
		return
	}
	db.DB().SetMaxIdleConns(0)
	conn = db
}

// Using this function to get a connection, you can create your connection pool here.
func GetDB() *gorm.DB {
	if conn == nil {
		Connect()
	} else if err := conn.DB().Ping(); err != nil {
		conn = nil
	}
	if conn == nil {
		sig := make(chan os.Signal, 1)
		signal.Notify(sig, syscall.SIGTERM, syscall.SIGINT, syscall.SIGKILL) //nolint
		var i int
	reconnect:
		for {
			select {
			case <-sig:
				break reconnect
			case <-time.After(tryPause * time.Second):
				Connect()
				i++
				if conn != nil {
					break reconnect
				}
			}
		}
		if conn == nil {
			log.Panicln(fmt.Sprintf("Connection failed after %d tries", i))
		}
	}
	return conn
}
