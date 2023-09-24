package requests

import (
	"fmt"
	"io"
	"log"
	"math"
	"net/http"
	"time"

	"github.com/alexmerren/httpcache"
	"golang.org/x/exp/slices"
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
		return nil, fmt.Errorf("unrecoverable error for url %s", URL)
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
		if iterationNumber == 5 {
			return nil, fmt.Errorf("srcom: maximum retry iterations exceeded for url %s", URL)
		}

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
