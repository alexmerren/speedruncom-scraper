package http_client

func WithRetry(retryDelay, retryMaximum int) func(*HttpClient) {
	return func(c *HttpClient) {
		c.retryDelay = retryDelay
		c.retryMaximum = retryMaximum
		c.doRetry = true
	}
}

// TODO Set the httpClient to have a cache around the roundtripper
func WithCacheRequests() func(*HttpClient) {
	return func(c *HttpClient) {}
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
