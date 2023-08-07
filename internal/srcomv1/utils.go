package srcomv1

import (
	"fmt"
	"io"
	"net/http"
	"time"
)

const (
	unsuccessfulRequestSleepTime = 5 * time.Second
)

func requestSrcom(URL string) ([]byte, error) {
	response, err := http.DefaultClient.Get(URL)
	if err != nil {
		return nil, err
	}

	if response.StatusCode != 200 {
		time.Sleep(unsuccessfulRequestSleepTime)
		return requestSrcom(URL)
	}

	fmt.Println(URL)

	defer response.Body.Close()
	return io.ReadAll(response.Body)
}
