package database

import (
	"database/sql"
	"errors"
	"fmt"
	"log"
	"os"
	"strings"
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

func (db Database) Close() {
	db.log.Println("closing database")
	defer db.database.Close()
}

func (db Database) RefreshDatabase() error {
	db.log.Println("refreshing database")

	allSeries, err := db.getAllSeries()
	if err != nil {
		err = fmt.Errorf("database.RefreshDatabase: %w", err)
		db.log.Println(err)

		return err
	}

	for i, r := range allSeries {
		t := r.Info.ComicSeries.TotalResultCount
		for j, series := range r.Records.ComicSeries {
			time.Sleep(apiDelay)

			description, err := db.getSeriesDescription(series.UUID)
			if err != nil {
				err = fmt.Errorf("database.RefreshDatabase: %w", err)
				db.log.Println(err)
				if errors.Is(err, apiResponseError{}) {
					db.log.Printf("skipping %v %v\n", series.UUID, series.Title)

					continue
				}

				return err
			}

			series.description = description
			c := i*100 + j + 1

			db.log.Printf("inserting %v/%v\n", c, t)

			err = db.insertSeries(series)
			if err != nil {
				err = fmt.Errorf("database.RefreshDatabase: %w", err)
				db.log.Println(err)

				return err
			}
		}
	}

	db.log.Println("done refreshing database")

	return nil
}

func (db Database) initialSetup() error {
	db.log.Println("setting up database")

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

	db.log.Println("database setup complete")

	return nil
}

func openLog(userHome string) (*log.Logger, error) {
	logPath := userHome + sep + ".dcui" + sep + "logs"

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

func (db Database) insertSeries(series SearchResultRecordsComicseries) error {
	db.log.Printf("upserting series %v\n", series.UUID)

	updateThreshold := time.Now()
	updateThreshold = updateThreshold.AddDate(-1, 0, 0)
	updateThresholdInt := updateThreshold.Unix()

	url := fmt.Sprintf("https://www.dcuniverseinfinite.com/comics/series/%v/%v", series.Slug, series.UUID)
	value := fmt.Sprintf("('%v','%v','%v',%v,%v,%v,%v,'%v')",
		sanitizeSQLString(series.UUID),
		sanitizeSQLString(series.Title),
		sanitizeSQLString(series.description),
		series.BooksCount,
		series.IssueCount,
		series.VolumeCount,
		series.OmnibusCount,
		url)
	qryUpsertSeries := fmt.Sprintf(templates["upsertSeries"], value, updateThresholdInt)

	_, err := db.database.Exec(qryUpsertSeries)
	if err != nil {
		err = fmt.Errorf("database.insertSeries: %w", err)
		db.log.Println(err)
		db.log.Println(qryUpsertSeries)

		return err
	}

	db.log.Println("upserting series genres")

	for _, g := range series.Genres {
		value = fmt.Sprintf("('%v','%v')", sanitizeSQLString(series.UUID), sanitizeSQLString(g))
		qryUpsertSeriesGenre := fmt.Sprintf(templates["upsertSeriesGenre"], value)

		_, err := db.database.Exec(qryUpsertSeriesGenre)
		if err != nil {
			err = fmt.Errorf("database.insertSeries: %w", err)
			db.log.Println(err)
			db.log.Println(qryUpsertSeriesGenre)

			return err
		}
	}

	db.log.Println("upserting series imprints")

	for _, i := range series.Imprints {
		value = fmt.Sprintf("('%v','%v')", sanitizeSQLString(series.UUID), sanitizeSQLString(i))
		qryUpsertSeriesImprint := fmt.Sprintf(templates["upsertSeriesImprint"], value)

		_, err := db.database.Exec(qryUpsertSeriesImprint)
		if err != nil {
			err = fmt.Errorf("database.insertSeries: %w", err)
			db.log.Println(err)
			db.log.Println(qryUpsertSeriesImprint)

			return err
		}
	}

	db.log.Println("series upsert complete")

	return nil
}

func sanitizeSQLString(s string) string {
	s = strings.ReplaceAll(s, "\r", " ")
	s = strings.ReplaceAll(s, "\n", " ")

	return strings.ReplaceAll(s, "'", "''")
}
