package httpcache

import "net/http"

type HTTPCacher interface {
	Do(request *http.Request) (response *http.Response, err error)
	Get(url string) (response *http.Response, err error)
}
