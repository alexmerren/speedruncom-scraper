package srcomv1

import (
	"io"
	"net/http"
)

func requestSrcom(URL string) ([]byte, error) {
	response, err := http.DefaultClient.Get(URL)
	if err != nil {
		return nil, err
	}

	defer response.Body.Close()
	return io.ReadAll(response.Body)
}
