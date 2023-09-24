package srcomv2

import (
	b64 "encoding/base64"
	"encoding/json"
	"fmt"
	"strings"
)

func formatHeader(data interface{}, function string) (string, error) {
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
