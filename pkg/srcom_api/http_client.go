package srcom_api

type HttpClient interface {
	Get(url string) ([]byte, error)
}
