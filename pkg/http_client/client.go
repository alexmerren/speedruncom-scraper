package http_client

import (
	"io"
	"log"
	"math"
	"net/http"
	"time"
)

const (
	defaultRetryDelay   = 0
	defaultRetryMaximum = 0
	defaultRequestDelay = 0
	defaultIsVerbose    = false
	defaultDoDelay      = false
	defaultDoRetry      = false
)

type HttpClient struct {
	httpClient   *http.Client
	retryDelay   int
	retryMaximum int
	requestDelay int
	isVerbose    bool
	doRetry      bool
	doDelay      bool
}

func NewHttpClient(options ...func(*HttpClient)) *HttpClient {
	client := &HttpClient{
		httpClient:   http.DefaultClient,
		retryDelay:   defaultRetryDelay,
		retryMaximum: defaultRetryMaximum,
		requestDelay: defaultRequestDelay,
		isVerbose:    defaultIsVerbose,
		doRetry:      defaultDoRetry,
		doDelay:      defaultDoDelay,
	}

	for _, optionFunc := range options {
		optionFunc(client)
	}

	return client
}

func (c *HttpClient) Get(url string) ([]byte, error) {
	response, err := c.httpClient.Get(url)
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()

	if c.doRetry && response.StatusCode != http.StatusOK {
		response, err = c.retry(url)
		if err != nil {
			return nil, err
		}
	}

	if c.isVerbose {
		log.Printf("%s: %d", url, response.StatusCode)
	}

	if c.doDelay {
		time.Sleep(time.Duration(c.requestDelay) * time.Millisecond)
	}

	return io.ReadAll(response.Body)
}

func (c *HttpClient) retry(url string) (*http.Response, error) {
	iterCount := 0
	for {
		if iterCount == c.retryMaximum {
			return nil, ErrRetryExceededMaximum
		}

		delayTime := c.retryDelay * int(math.Pow(2, float64(iterCount)))
		time.Sleep(time.Duration(delayTime) * time.Millisecond)

		response, err := c.httpClient.Get(url)
		if err != nil {
			return nil, err
		}

		if response.StatusCode == http.StatusOK {
			return response, nil
		}

		iterCount += 1
	}
}
