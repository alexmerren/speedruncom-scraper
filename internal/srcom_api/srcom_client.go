package srcom_api

import (
	"github.com/alexmerren/httpcache"
	"github.com/alexmerren/speedruncom-scraper/internal/http_client"
)

var cachedRoundTripper = httpcache.NewCachedRoundTripper(
	httpcache.WithName("./data/httpcache.db"),
)

func NewSrcomClient() *SrcomV1Client {
	return NewSrcomV1Client(
		http_client.NewHttpClient(
			http_client.WithLogging(),
			http_client.WithRetry(100, 5),
			http_client.WithCache(cachedRoundTripper),
		),
	)
}

func NewSrcomClientV2() *SrcomV2Client {
	return NewSrcomV2Client(
		http_client.NewHttpClient(
			http_client.WithLogging(),
			http_client.WithRetry(100, 5),
			http_client.WithCache(cachedRoundTripper),
		),
	)
}
