package database

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

const (
	retryDelay  = 30 * time.Second
	apiDelay    = 50 * time.Millisecond
	httpTimeout = time.Minute
	startPage   = 1
)

type apiResponseError struct {
	statusCode int
	status     string
}

func (e apiResponseError) Error() string {
	return fmt.Sprintf("%v %v", e.statusCode, e.status)
}

func (e apiResponseError) Is(target error) bool {
	return target == apiResponseError{}
}

func (db Database) getAllSeries() ([]SearchResult, error) {
	db.log.Println("getting all series from DCUI API")

	const recordsPerPage = 100
	reqBody := SearchBody{
		EngineKey:     engineKey, // engineKey is in creds.go, not synced due to security concerns
		Page:          startPage,
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

	db.log.Printf("retrieving first %v records\n", recordsPerPage)

	singleResult, err := db.requestSeries(reqBody)
	if err != nil {
		err = fmt.Errorf("database.getAllSeries: %w", err)
		db.log.Println(err)

		return nil, err
	}

	searchResults := []SearchResult{singleResult}
	numPages := singleResult.Info.ComicSeries.NumPages
	nextPage := startPage + 1

	for p := nextPage; p <= numPages; p++ {
		time.Sleep(apiDelay)
		db.log.Printf("retrieving records %v/%v\n", p*recordsPerPage, singleResult.Info.ComicSeries.TotalResultCount)

		reqBody.Page = p

		singleResult, err = db.requestSeries(reqBody)
		if err != nil {
			err = fmt.Errorf("database.getAllSeries: %w", err)
			db.log.Println(err)

			return nil, err
		}

		searchResults = append(searchResults, singleResult)
	}

	db.log.Println("done getting all series")

	return searchResults, nil
}

func (db Database) requestSeries(reqBody SearchBody) (SearchResult, error) {
	db.log.Println("requesting series")

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

	db.log.Println("series retrieved")

	return searchResult, nil
}

func post(uri string, data []byte) ([]byte, error) {
	httpClient := &http.Client{
		Timeout: httpTimeout,
	}

	resp, err := httpClient.Post(uri, "application/json", bytes.NewBuffer(data))
	if err != nil {
		time.Sleep(retryDelay)

		resp, err = httpClient.Post(uri, "application/json", bytes.NewBuffer(data))
		if err != nil {
			err = fmt.Errorf("database.post: %w", err)

			return nil, err
		}
	}

	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		err = apiResponseError{
			statusCode: resp.StatusCode,
			status:     resp.Status,
		}

		return nil, err
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		err = fmt.Errorf("database.post: %w", err)

		return nil, err
	}

	return body, nil
}

func (db Database) getSeriesDescription(uuid string) (string, error) {
	db.log.Println("getting series description")

	uri := fmt.Sprintf("https://www.dcuniverseinfinite.com/api/comics/1/series/%v/?trans=en", uuid)

	resp, err := get(uri)
	if err != nil {
		err = fmt.Errorf("database.getSeriesDescription: %w", err)
		db.log.Println(err)
		db.log.Println(uri)

		return "", err
	}

	var seriesDetail SeriesDetail

	err = json.Unmarshal(resp, &seriesDetail)
	if err != nil {
		err = fmt.Errorf("database.getSeriesDescription: %w", err)
		db.log.Println(err)
		db.log.Println(string(resp))

		return "", err
	}

	db.log.Println("series description retrieved")

	return seriesDetail.Description, nil
}

func get(uri string) ([]byte, error) {
	httpClient := &http.Client{
		Timeout: httpTimeout,
	}

	req, err := http.NewRequestWithContext(context.Background(), http.MethodGet, uri, nil)
	if err != nil {
		err = fmt.Errorf("database.get: %w", err)

		return nil, err
	}

	req.Header.Add("X-Consumer-Key", xConsumerKey)

	resp, err := httpClient.Do(req)
	if err != nil {
		time.Sleep(retryDelay)

		resp, err = httpClient.Do(req)
		if err != nil {
			err = fmt.Errorf("database.get: %w", err)

			return nil, err
		}
	}

	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		err = apiResponseError{
			statusCode: resp.StatusCode,
			status:     resp.Status,
		}

		return nil, err
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		err = fmt.Errorf("database.get: %w", err)

		return nil, err
	}

	return body, nil
}
