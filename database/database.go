package database

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"time"

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

func (db *Database) UpdateDatabase() error {
	allSeries, err := getAllSeries()
	if err != nil {
		err = fmt.Errorf("database.UpdateDatabase: %w", err)
		db.log.Println(err)

		return err
	}

	for _, r := range allSeries {
		for _, s := range r.Records.Comicseries {
			db.log.Printf("%v : %v\n", s.Title, s.UUID)
		}
	}

	return nil
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

func getAllSeries() ([]SearchResult, error) {
	url := "https://search.dcuniverseinfinite.com/api/v1/public/engines/search.json"
	reqBody := searchBody{
		engine_key:     engineKey, // engineKey is in creds.go, not synced due to security concerns
		page:           1,
		per_page:       100,
		document_types: []string{"comicseries"},
		filters:        map[string]string{},
		sort_field: map[string]string{
			"comicseries": "first_released",
		},
		sort_direction: map[string]string{
			"comicseries": "desc",
		},
	}

	jsonData, err := json.Marshal(reqBody)
	if err != nil {
		err = fmt.Errorf("database.getAllSeries: %w", err)

		return nil, err
	}

	resp, err := post(url, jsonData)
	if err != nil {
		err = fmt.Errorf("database.getAllSeries: %w", err)

		return nil, err
	}

	searchResults := []SearchResult{}

	var singleResult SearchResult

	err = json.Unmarshal(resp, &singleResult)
	if err != nil {
		err = fmt.Errorf("database.getAllSeries: %w", err)

		return nil, err
	}

	numPages := singleResult.Info.Comicseries.Num_pages
	searchResults = append(searchResults, singleResult)

	for p := 2; p <= numPages; p++ {
		reqBody.page = p

		jsonData, err = json.Marshal(reqBody)
		if err != nil {
			err = fmt.Errorf("database.getAllSeries: %w", err)

			return nil, err
		}

		resp, err = post(url, jsonData)
		if err != nil {
			err = fmt.Errorf("database.getAllSeries: %w", err)

			return nil, err
		}

		err = json.Unmarshal(resp, &singleResult)
		if err != nil {
			err = fmt.Errorf("database.getAllSeries: %w", err)

			return nil, err
		}

		searchResults = append(searchResults, singleResult)
	}

	return searchResults, nil
}

func post(uri string, data []byte) ([]byte, error) {
	resp, err := http.Post(uri, "application/json", bytes.NewBuffer(data))
	if err != nil {
		if resp != nil {
			resp.Body.Close()
		}
		time.Sleep(30 * time.Second)

		resp, err = http.Post(uri, "application/json", bytes.NewBuffer(data))
		if err != nil {
			if resp != nil {
				resp.Body.Close()
			}

			err = fmt.Errorf("database.post: %w", err)

			return nil, err
		}
	}

	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		err = fmt.Errorf("database.post: %w", err)

		return nil, err
	}

	return body, nil
}
