package srcom_api

import (
	"github.com/alexmerren/httpcache"
	"github.com/alexmerren/speedruncom-scraper/internal/http_client"
)

type HttpClient interface {
	Get(url string) ([]byte, error)
}

var cachedRoundTripper = httpcache.NewCachedRoundTripper(
	httpcache.WithName("./data/httpcache.db"),
)

var DefaultV1Client = NewSrcomV1Client(
	http_client.NewHttpClient(
		http_client.WithLogging(),
		http_client.WithRetry(100, 5),
		http_client.WithCache(cachedRoundTripper),
	),
)

var DefaultV2Client = NewSrcomV2Client(
	http_client.NewHttpClient(
		http_client.WithLogging(),
		http_client.WithRetry(1_000, 3),
		http_client.WithCache(cachedRoundTripper),
	),
)
