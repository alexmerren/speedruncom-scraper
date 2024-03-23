package http_client

import "net/http"

func WithRetry(retryDelay, retryMaximum int) func(*HttpClient) {
	return func(c *HttpClient) {
		c.retryDelay = retryDelay
		c.retryMaximum = retryMaximum
		c.doRetry = true
	}
}

func WithCache(cachedTransport http.RoundTripper) func(*HttpClient) {
	return func(c *HttpClient) {
		c.httpClient = &http.Client{
			Transport: cachedTransport,
		}
	}
}

func WithDelay(requestDelay int) func(*HttpClient) {
	return func(c *HttpClient) {
		c.requestDelay = requestDelay
		c.doDelay = true
	}
}

func WithLogging() func(*HttpClient) {
	return func(c *HttpClient) {
		c.isVerbose = true
	}
}
