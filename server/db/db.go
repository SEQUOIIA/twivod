package db

import (
	"database/sql"
	"log"
	"os"
)

type Db struct {
	db * sql.DB
}

func NewDb(path string) *Db {
	d := &Db{}

	var (
		createDatabase bool
		err error
	)

	if _, err = os.Stat(path); os.IsNotExist(err) {
		createDatabase = true
	} else {
		createDatabase = false
	}

	d.db, err = sql.Open("sqlite3", path)
	if err != nil {
		log.Fatal(err)
	}

	d.db.SetMaxOpenConns(1)

	_, err = d.db.Exec("PRAGMA journal_mode=WAL")
	if err != nil {
		log.Fatal(err)
	}

	if createDatabase {
		err = d.createDb()
		if err != nil {
			log.Fatal(err)
		}
	}

	return d
}

func (d * Db) GetDb() *sql.DB {
	return d.db
}

func (d * Db) createDb() error {
	// Create twitch_user table
	err := d.createTwitchUserTable()
	if err != nil {
		return err
	}

	return nil
}

func (d * Db) createTwitchUserTable() error{
	_, err := d.db.Exec("CREATE TABLE twitch_user " +
		"(id INTEGER PRIMARY KEY AUTOINCREMENT NOT NULL, " +
		"twitch_id TEXT, " +
		"name TEXT, " +
		"bio TEXT," +
		"created_at TEXT," +
		"display_name TEXT," +
		"logo TEXT," +
		"type TEXT," +
		"updated_at TEXT);")
	if err != nil {
		return err
	}
	return nil
}