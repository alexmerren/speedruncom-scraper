package srcomv1

import (
	"io"
	"log"
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

	if response.StatusCode == 429 {
		defer response.Body.Close()
		log.Fatal(response.Header, response.Body)
	}

	if response.StatusCode != 200 {
		time.Sleep(unsuccessfulRequestSleepTime)
		return requestSrcom(URL)
	}

	log.Print(URL)
	defer response.Body.Close()

	return io.ReadAll(response.Body)
}
