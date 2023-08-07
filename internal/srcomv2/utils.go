package srcomv2

import (
	b64 "encoding/base64"
	"encoding/json"
	"fmt"
	"io"
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

	if response.StatusCode == 429 {
		defer response.Body.Close()
		fmt.Println(response.Header, response.Body)
		return nil, fmt.Errorf("Rate limit has been hit.")
	}

	if response.StatusCode != 200 {
		time.Sleep(unsuccessfulRequestSleepTime)
		return requestSrcom(URL)
	}

	fmt.Println(URL)

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

func decodeB64Header(header string) ([]byte, error) {
	if i := len(header) % 4; i != 0 {
		header += strings.Repeat("=", 4-i)
	}

	decodedString, err := b64.StdEncoding.DecodeString(header)
	if err != nil {
		return nil, err
	}

	return decodedString, nil
}

func encodeB64Header(data []byte) string {
	encodedString := b64.StdEncoding.EncodeToString(data)
	return strings.TrimRight(encodedString, "=")
}
