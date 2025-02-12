package config

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/mattn/go-sqlite3"
)

func dbSetup() error {
	db, err := sql.Open("sqlite3", BaseDirectory+"/htb-cli.db")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	createTableSQL := `
	CREATE TABLE IF NOT EXISTS machines (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		name TEXT NOT NULL,
		os TEXT NOT NULL,
		difficulty TEXT NOT NULL,
		star REAL NOT NULL,
		status TEXT NOT NULL,
		release_date TIMESTAMP NOT NULL,
		user BOOLEAN NOT NULL,
		root BOOLEAN NOT NULL
	);`

	_, err = db.Exec(createTableSQL)
	if err != nil {
		return fmt.Errorf("error during table creation: %v", err)
	}

	GlobalConfig.Logger.Info("Machines' table successfully created !")
	return nil
}
