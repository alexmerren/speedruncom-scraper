package srcom_api

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"strings"
)

const (
	v2ApiUrl        = "https://www.speedrun.com/api/v2/%s"
	v2GameListQuery = "GetGameList?_r=%s"
)

type SrcomV2Client struct {
	client HttpClient
}

func NewSrcomV2Client(client HttpClient) *SrcomV2Client {
	return &SrcomV2Client{
		client: client,
	}
}

func (c *SrcomV2Client) GetGameList(pageNumber int) ([]byte, error) {
	data := map[string]interface{}{
		"page": pageNumber,
	}

	url, err := formatUrl(v2GameListQuery, data)
	if err != nil {
		return nil, err
	}

	return c.client.Get(url)
}

func formatUrl(query string, data map[string]interface{}) (string, error) {
	encodedData, err := encodeData(data)
	if err != nil {
		return encodedData, err
	}
	formattedQuery := fmt.Sprintf(query, encodedData)

	return fmt.Sprintf(v2ApiUrl, formattedQuery), nil
}

func encodeData(data map[string]interface{}) (string, error) {
	bytes, err := json.Marshal(data)
	if err != nil {
		return "", err
	}

	encodedString := base64.StdEncoding.EncodeToString(bytes)
	return strings.TrimRight(encodedString, "="), nil
}
