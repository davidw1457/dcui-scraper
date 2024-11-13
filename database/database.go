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

func (db *Database) Dummy() {
	// TODO: Get rid of this func
	fmt.Println("this is a dummy func to get around unused vars during dev. Delete me")
}

func (db *Database) initialSetup() error {
	// TODO: Setup DCUI DB
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
