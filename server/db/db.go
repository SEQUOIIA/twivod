package db

import (
	"database/sql"
	"log"
	"os"
)

type Dbc struct {
	Dbc * sql.DB
	ps preppedStatements
}

type preppedStatements struct {
	addTwitchUser * sql.Stmt
}

func NewDbc(path string) *Dbc {
	d := &Dbc{}

	var (
		createDatabase bool
		err error
	)

	if _, err = os.Stat(path); os.IsNotExist(err) {
		createDatabase = true
	} else {
		createDatabase = false
	}

	d.Dbc, err = sql.Open("sqlite3", path)
	if err != nil {
		log.Fatal(err)
	}

	d.Dbc.SetMaxOpenConns(1)

	_, err = d.Dbc.Exec("PRAGMA journal_mode=WAL")
	if err != nil {
		log.Fatal(err)
	}

	if createDatabase {
		err = d.createDbc()
		if err != nil {
			log.Fatal(err)
		}
	}

	return d
}

func (d * Dbc) GetDbc() *sql.DB {
	return d.Dbc
}

func (d * Dbc) createDbc() error {
	// Create twitch_user table
	err := d.createTwitchUserTable()
	if err != nil {
		return err
	}

	return nil
}

func (d * Dbc) createTwitchUserTable() error{
	_, err := d.Dbc.Exec("CREATE TABLE twitch_user " +
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