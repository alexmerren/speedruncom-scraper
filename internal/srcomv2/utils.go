package srcomv2

import (
	b64 "encoding/base64"
	"encoding/json"
	"io"
	"net/http"
	"strings"
)

func requestSrcom(URL string) ([]byte, error) {
	response, err := http.DefaultClient.Get(URL)
	if err != nil {
		return nil, err
	}

	defer response.Body.Close()
	return io.ReadAll(response.Body)
}

func getBytes(data interface{}) ([]byte, error) {
	return json.Marshal(data)
}

func DecodeB64Header(header string) ([]byte, error) {
	if i := len(header) % 4; i != 0 {
		header += strings.Repeat("=", 4-i)
	}
	decodedString, err := b64.StdEncoding.DecodeString(header)
	if err != nil {
		return nil, err
	}
	return decodedString, nil
}

func EncodeB64Header(data []byte) string {
	encodedString := b64.StdEncoding.EncodeToString(data)
	return strings.TrimRight(encodedString, "=")
}
