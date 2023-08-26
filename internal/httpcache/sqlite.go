package httpcache

import (
	"bytes"
	"database/sql"
	"errors"
	"io"
	"net/http"
	"strings"

	_ "github.com/mattn/go-sqlite3"
)

const (
	defaultDatabaseName = "./data/httpcache.db"

	createDatabaseQuery = "CREATE TABLE IF NOT EXISTS responses (url TEXT PRIMARY KEY, body TEXT)"
	insertRequestQuery  = "INSERT INTO responses (url, body) VALUES (?, ?)"
	selectRequestQuery  = "SELECT body FROM responses WHERE url = ?"
)

var DefaultClient = NewHTTPCacheClient(defaultDatabaseName)

type HTTPCacheClient struct {
	httpClient *http.Client
	database   *sql.DB
}

type cachedRequest struct {
	url  string
	body string
}

func NewHTTPCacheClient(dbname string) *HTTPCacheClient {
	conn, err := sql.Open("sqlite3", dbname)
	if err != nil {
		return nil
	}

	_, err = conn.Exec(createDatabaseQuery)
	if err != nil {
		return nil
	}

	return &HTTPCacheClient{
		httpClient: http.DefaultClient,
		database:   conn,
	}
}

func (h *HTTPCacheClient) Do(request *http.Request) (*http.Response, error) {
	cachedResponse := &cachedRequest{}
	row := h.database.QueryRow(selectRequestQuery, request.URL.RequestURI())
	err := row.Scan(&cachedResponse.body)
	if err == nil {
		return &http.Response{
			StatusCode: 200,
			Body:       io.NopCloser(strings.NewReader(cachedResponse.body)),
		}, nil
	}

	// If there is no cached response and there is some other error, we want to
	// return that other than continue the HTTP request.
	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		return nil, err
	}

	response, err := h.httpClient.Do(request)
	if err != nil {
		return nil, err
	}

	// Return nothing but the 'bad' status code, so that our bad request
	// handling logic can do exponential backoff and re-requesting.
	if response.StatusCode != 200 {
		return &http.Response{
			StatusCode: response.StatusCode,
		}, nil
	}

	// We read the body once to store, and then reset the response to the
	// beginning of the original response. This is done so we can effectively
	// read the response body twice.
	responseBody, err := io.ReadAll(response.Body)
	if err != nil {
		return nil, err
	}
	response.Body = io.NopCloser(bytes.NewBuffer(responseBody))

	_, err = h.database.Exec(insertRequestQuery, request.URL.RequestURI(), responseBody)
	if err != nil {
		return nil, err
	}

	return response, nil
}

func (h *HTTPCacheClient) Get(url string) (*http.Response, error) {
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}
	return h.Do(req)
}
