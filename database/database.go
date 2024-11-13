package database

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	_ "github.com/mattn/go-sqlite3" // Required to use sqlite3 driver
)

const (
	userRWX = 0o700
	sep     = string(os.PathSeparator)
)

type Database struct {
	database *sql.DB
	log      *log.Logger
}

func New() (Database, error) {
	var dcuiDB Database

	userHome, err := os.UserHomeDir()
	if err != nil {
		err = fmt.Errorf("database.New: %w", err)

		return dcuiDB, err
	}

	logger, err := openLog(userHome)
	if err != nil {
		err = fmt.Errorf("database.New: %w", err)

		return dcuiDB, err
	}

	dcuiDB.log = logger

	dcuiDB.log.Println("opening database")

	dbase, err := openDB(userHome)
	if err != nil {
		err = fmt.Errorf("database.New: %w", err)
		dcuiDB.log.Println(err)

		return dcuiDB, err
	}

	dcuiDB.database = dbase

	dcuiDB.log.Println("performing initial setup")

	err = dcuiDB.initialSetup()
	if err != nil {
		err = fmt.Errorf("database.New: %w", err)
		dcuiDB.log.Println(err)

		return dcuiDB, err
	}

	dcuiDB.log.Println("database opened")

	return dcuiDB, nil
}

func (db *Database) Close() {
	db.log.Println("closing database")
	defer db.database.Close()
}

func (db *Database) initialSetup() error {
	rows, err := db.database.Query(queries["pingDatabase"])
	if rows != nil {
		defer rows.Close()
	}
	if err != nil || rows.Err() != nil {
		_, err = db.database.Exec(queries["createDatabase"])
		if err != nil {
			err = fmt.Errorf("database.initialSetup: %w", err)
			db.log.Println(err)

			return err
		}
	}

	return nil
}

func openLog(userHome string) (*log.Logger, error) {
	logPath := userHome + sep + "logs"

	logFile, err := os.OpenFile(logPath+sep+"dcui-scraper.log",
		os.O_WRONLY|os.O_CREATE|os.O_APPEND, userRWX)
	if err != nil {
		err = fmt.Errorf("database.openLog: %w", err)

		return nil, err
	}

	logger := log.New(logFile, "database: ", log.LstdFlags)

	return logger, nil
}

func openDB(userHome string) (*sql.DB, error) {
	databasePath := userHome + sep + ".dcui"

	_, err := os.Stat(databasePath)
	if err != nil {
		if os.IsNotExist(err) {
			err = os.MkdirAll(databasePath, userRWX)
			if err != nil {
				err = fmt.Errorf("database.openDB: %w", err)

				return nil, err
			}
		} else {
			err = fmt.Errorf("database.openDB: %w", err)

			return nil, err
		}
	}

	databaseFile := databasePath + sep + "dcui.db"

	dbase, err := sql.Open("sqlite3", databaseFile)
	if err != nil {
		err = fmt.Errorf("database.openDB: %w", err)

		return nil, err
	}

	return dbase, nil
}
