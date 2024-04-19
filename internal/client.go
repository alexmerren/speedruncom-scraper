package internal

import (
	"github.com/alexmerren/httpcache"
	"github.com/alexmerren/speedruncom-scraper/pkg/http_client"
	"github.com/alexmerren/speedruncom-scraper/pkg/srcom_api"
)

var cachedRoundTripper = httpcache.NewCachedRoundTripper(
	httpcache.WithName("./data/httpcache.db"),
)

func NewSrcomV1Client() *srcom_api.SrcomV1Client {
	return srcom_api.NewSrcomV1Client(
		http_client.NewHttpClient(
			http_client.WithRetry(250, 3),
			http_client.WithCache(cachedRoundTripper),
		),
	)
}

func NewSrcomV2Client() *srcom_api.SrcomV2Client {
	return srcom_api.NewSrcomV2Client(
		http_client.NewHttpClient(
			http_client.WithDelay(500),
			http_client.WithRetry(500, 2),
			http_client.WithCache(cachedRoundTripper),
		),
	)
}
