package srcom_api

import (
	"log/slog"

	"github.com/alexmerren/httpcache"
	"github.com/alexmerren/speedruncom-scraper/internal/http_client"
)

type HttpClient interface {
	Get(url string) ([]byte, error)
}

var DefaultV1Client = NewSrcomV1Client(
	http_client.NewHttpClient(
		http_client.WithLogging(),
		http_client.WithRetry(100, 5),
		http_client.WithCache(newCache()),
	),
)

var DefaultV2Client = NewSrcomV2Client(
	http_client.NewHttpClient(
		http_client.WithLogging(),
		http_client.WithRetry(1_000, 3),
		http_client.WithCache(newCache()),
	),
)

func newCache() *httpcache.CachedRoundTripper {
	cache, err := httpcache.NewCachedRoundTripper(
		httpcache.WithName("./data/httpcache.db"),
	)

	if err != nil {
		slog.Error("Failed to create HTTP cache", "error", err)
	}

	return cache
}
