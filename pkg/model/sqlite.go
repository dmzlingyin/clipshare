//go:build server
// +build server

package model

import (
	"database/sql"
	"errors"
	"os"

	"github.com/dmzlingyin/clipshare/pkg/log"
	_ "github.com/mattn/go-sqlite3"
)

var DB *sql.DB

func init() {
	if _, err := os.Stat("./users.db"); os.IsNotExist(err) {
		DB, err = sql.Open("sqlite3", "./users.db")
		if err != nil {
			log.ErrorLogger.Fatal(err)
			DB.Close()
		}

		sqlStmt := "create table users (id integer not null primary key autoincrement, name text, password text)"
		_, err = DB.Exec(sqlStmt)
		if err != nil {
			log.ErrorLogger.Fatal(err)
			DB.Close()
		}
	} else {
		db, err := sql.Open("sqlite3", "./users.db")
		if err != nil {
			log.ErrorLogger.Fatal(err)
			db.Close()
		}
		DB = db
	}
}

func Check(username, password string) error {
	id := userID(username)
	if id == -1 {
		return errors.New("user does't exit")
	}

	stmt, err := DB.Prepare("select password from users where id = ?")
	if err != nil {
		log.ErrorLogger.Println(err)
		return err
	}
	defer stmt.Close()

	var passwd string
	err = stmt.QueryRow(id).Scan(&passwd)
	if err != nil {
		log.ErrorLogger.Println(err)
		return err
	}
	if passwd != password {
		return errors.New("password wrong")
	}
	return nil
}

func Register(username, password string) error {
	if id := userID(username); id != -1 {
		return errors.New("user already exist")
	}

	stmt, err := DB.Prepare("insert into users(name, password) values(?, ?)")
	if err != nil {
		return err
	}
	defer stmt.Close()

	_, err = stmt.Exec(username, password)
	if err != nil {
		return err
	}
	return nil
}

func userID(username string) int {
	id := -1
	stmt, err := DB.Prepare("select id from users where name = ?")
	if err != nil {
		log.ErrorLogger.Println(err)
		return id
	}
	defer stmt.Close()

	err = stmt.QueryRow(username).Scan(&id)
	if err != nil {
		log.ErrorLogger.Println(err)
		return id
	}
	return id
}
