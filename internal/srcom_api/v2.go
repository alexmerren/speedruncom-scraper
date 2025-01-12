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

	requestQuery, err := generateRequestQuery(v2GameListQuery, data)
	if err != nil {
		return nil, err
	}

	return c.get(requestQuery)
}

func (c *SrcomV2Client) GetGameRecordHistory(gameId, categoryId string, levelId *string, variables, values []string) ([]byte, error) {
	data := generateGameRecordHistoryData(gameId, categoryId, levelId, variables, values)
	requestQuery, err := generateRequestQuery(v2GameRecordQuery, data)
	if err != nil {
		return nil, err
	}

	return c.get(requestQuery)
}

func generateGameRecordHistoryData(gameId, categoryId string, levelId *string, variables, values []string) map[string]any {
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

func generateRequestQuery(query string, data map[string]any) (string, error) {
	dataBytes, err := json.Marshal(data)
	if err != nil {
		return "", err
	}

	encodedData := base64.StdEncoding.EncodeToString(dataBytes)
	trimmedData := strings.TrimRight(encodedData, "=")

	return fmt.Sprintf(query, trimmedData), nil
}

func (c *SrcomV2Client) get(requestQuery string) ([]byte, error) {
	url := fmt.Sprintf(v2ApiUrl, requestQuery)
	return c.client.Get(url)
}
