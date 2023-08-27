package srcomv2

import (
	b64 "encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"math"
	"net/http"
	"slices"
	"strings"
	"time"

	"github.com/alexmerren/speedruncom-scraper/internal/httpcache"
)

const (
	exponentialBackoffStartInt   = 500
	exponentialBackoffMultiplier = 2
)

var unrecovableStatusCodes = []int{
	http.StatusNotFound,
	http.StatusBadRequest,
	http.StatusGatewayTimeout,
	http.StatusServiceUnavailable,
}

func RequestSrcom(URL string) ([]byte, error) {
	response, err := httpcache.DefaultClient.Get(URL)
	if err != nil {
		return nil, err
	}

	if slices.Contains(unrecovableStatusCodes, response.StatusCode) {
		return nil, fmt.Errorf("srcom: unrecoverable error for url %s", URL)
	}

	if response.StatusCode != http.StatusOK {
		response, err = retryWithExponentialBackoff(URL)
		if err != nil {
			return nil, err
		}
	}

	log.Print(URL)
	defer response.Body.Close()

	return io.ReadAll(response.Body)
}

func retryWithExponentialBackoff(URL string) (*http.Response, error) {
	iterationNumber := 0
	for {
		backoffTime := exponentialBackoff(iterationNumber)
		log.Printf("Sleeping for %s", backoffTime)
		time.Sleep(backoffTime)

		response, err := httpcache.DefaultClient.Get(URL)
		if err != nil {
			return nil, err
		}

		if response.StatusCode == http.StatusOK {
			return response, nil
		}

		iterationNumber += 1
	}
}

func exponentialBackoff(iteration int) time.Duration {
	newTime := exponentialBackoffStartInt * math.Pow(exponentialBackoffMultiplier, float64(iteration))
	return time.Duration(newTime) * time.Millisecond
}

func formatHeader(data interface{}, function string) (string, error) {
	flattenedData, err := getBytes(data)
	if err != nil {
		return "", err
	}

	formattedData := encodeB64Header(flattenedData)
	URL := fmt.Sprintf(baseApiUrl, function, formattedData)
	return URL, nil
}

func getBytes(data interface{}) ([]byte, error) {
	return json.Marshal(data)
}

func encodeB64Header(data []byte) string {
	encodedString := b64.StdEncoding.EncodeToString(data)
	return strings.TrimRight(encodedString, "=")
}
