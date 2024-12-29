package srcom_api

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"strings"
)

const (
	v2ApiUrl          = "https://www.speedrun.com/api/v2/%s"
	v2GameListQuery   = "GetGameList?_r=%s"
	v2GameRecordQuery = "GetGameRecordHistory?_r=%s"
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
	data := map[string]any{
		"page": pageNumber,
	}

	url, err := formatUrl(v2GameListQuery, data)
	if err != nil {
		return nil, err
	}

	return c.client.Get(url)
}

func (c *SrcomV2Client) GetGameRecordHistory(gameId, categoryId string, levelId *string, variables, values []string) ([]byte, error) {
	data := generateGameRecordHistoryQuery(gameId, categoryId, levelId, variables, values)
	url, err := formatUrl(v2GameRecordQuery, data)
	if err != nil {
		return nil, err
	}

	return c.client.Get(url)
}

func generateGameRecordHistoryQuery(gameId, categoryId string, levelId *string, variables, values []string) map[string]any {
	variablesAndValues := make([]map[string]any, len(variables))

	for index := range len(variables) {
		variablesAndValues[index] = map[string]any{
			"variableId": variables[index],
			"valueIds":   []string{values[index]},
		}
	}

	return map[string]any{
		"params": map[string]any{
			"categoryId": categoryId,
			"gameId":     gameId,
			"levelId":    levelId,
			"values":     variablesAndValues,
		},
	}
}

func formatUrl(query string, data map[string]any) (string, error) {
	encodedData, err := encodeData(data)
	if err != nil {
		return "", err
	}
	formattedQuery := fmt.Sprintf(query, encodedData)

	return fmt.Sprintf(v2ApiUrl, formattedQuery), nil
}

func encodeData(data map[string]any) (string, error) {
	bytes, err := json.Marshal(data)
	if err != nil {
		return "", err
	}

	encodedString := base64.StdEncoding.EncodeToString(bytes)
	return strings.TrimRight(encodedString, "="), nil
}
