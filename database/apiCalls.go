package database

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

func (db *Database) getAllSeries() ([]SearchResult, error) {
	const recordsPerPage = 100

	reqBody := SearchBody{
		EngineKey:     engineKey, // engineKey is in creds.go, not synced due to security concerns
		Page:          1,
		PerPage:       recordsPerPage,
		DocumentTypes: []string{"comicseries"},
		Filters:       map[string]string{},
		SortField: map[string]string{
			"comicseries": "first_released",
		},
		SortDirection: map[string]string{
			"comicseries": "desc",
		},
	}

	singleResult, err := db.requestSeries(reqBody)
	if err != nil {
		err = fmt.Errorf("database.getAllSeries: %w", err)
		db.log.Println(err)

		return nil, err
	}

	searchResults := []SearchResult{}

	numPages := singleResult.Info.ComicSeries.NumPages
	searchResults = append(searchResults, singleResult)

	for p := 2; p <= numPages; p++ {
		reqBody.Page = p

		singleResult, err = db.requestSeries(reqBody)
		if err != nil {
			err = fmt.Errorf("database.getAllSeries: %w", err)
			db.log.Println(err)

			return nil, err
		}

		searchResults = append(searchResults, singleResult)
	}

	return searchResults, nil
}

func (db *Database) requestSeries(reqBody SearchBody) (SearchResult, error) {
	const uri = "https://search.dcuniverseinfinite.com/api/v1/public/engines/search.json"

	var searchResult SearchResult

	jsonData, err := json.Marshal(reqBody)
	if err != nil {
		err = fmt.Errorf("database.requestSeries: %w", err)
		db.log.Println(err)

		return searchResult, err
	}

	resp, err := post(uri, jsonData)
	if err != nil {
		err = fmt.Errorf("database.requestSeries: %w", err)
		db.log.Println(err)
		db.log.Println(string(jsonData))

		return searchResult, err
	}

	err = json.Unmarshal(resp, &searchResult)
	if err != nil {
		err = fmt.Errorf("database.requestSeries: %w", err)
		db.log.Println(err)
		db.log.Println(string(resp))

		return searchResult, err
	}

	return searchResult, nil
}

func post(uri string, data []byte) ([]byte, error) {
	const retryDelay = 30 * time.Second

	resp, err := http.Post(uri, "application/json", bytes.NewBuffer(data))
	if err != nil {
		if resp != nil {
			resp.Body.Close()
		}

		time.Sleep(retryDelay)

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
