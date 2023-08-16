package srcomv2

import (
	b64 "encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"strings"
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

	// Still investigating how to respond to a 429, so if we hit it, we will know!
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

func formatHeader(data map[string]interface{}, function string) (string, error) {
	flattenedData, err := getBytes(data)
	if err != nil {
		return "", err
	}

	formattedData := encodeB64Header(flattenedData)
	URL := fmt.Sprintf(baseApiUrl, function, formattedData)
	return URL, nil
}

func getBytes(data interface{}) ([]byte, error) {
	return json.Marshal(data)
}

func encodeB64Header(data []byte) string {
	encodedString := b64.StdEncoding.EncodeToString(data)
	return strings.TrimRight(encodedString, "=")
}
