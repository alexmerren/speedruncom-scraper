package srcom_api

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"strings"
)

const (
	v2ApiUrl                 = "https://www.speedrun.com/api/v2/%s"
	v2UserDataQuery          = "GetUserLeaderboard?_r=%s"
	v2UserSummaryQuery       = "GetUserSummary?_r=%s"
	v2GameDataQuery          = "GetGameData?_r=%s"
	v2GameSummaryQuery       = "GetGameSummary?_r=%s"
	v2GameListQuery          = "GetGameList?_r=%s"
	v2GameLeaderboardQuery   = "GetGameLeaderboard2?_r=%s"
	v2GameRecordHistoryQuery = "GetGameRecordHistory?_r=%s"
)

type SrcomV2Client struct {
	client HttpClient
}

func NewSrcomV2Client(client HttpClient) *SrcomV2Client {
	return &SrcomV2Client{
		client: client,
	}
}

func (c *SrcomV2Client) GetUser(userID string) ([]byte, error) {
	data := map[string]interface{}{
		"userId":    userID,
		"levelType": 1,
	}

	url, err := formatUrl(v2UserDataQuery, data)
	if err != nil {
		return nil, err
	}

	return c.client.Get(url)
}

func (c *SrcomV2Client) GetGame(gameID string) ([]byte, error) {
	data := map[string]interface{}{
		"gameId": gameID,
	}

	url, err := formatUrl(v2GameDataQuery, data)
	if err != nil {
		return nil, err
	}

	return c.client.Get(url)
}

func (c *SrcomV2Client) GetGameSummary(gameurl string) ([]byte, error) {
	data := map[string]interface{}{
		"gameUrl": gameurl,
	}

	url, err := formatUrl(v2GameSummaryQuery, data)
	if err != nil {
		return nil, err
	}

	return c.client.Get(url)
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

func (c *SrcomV2Client) GetGameCategoryLeaderboard(gameID, categoryID string, pageNumber int) ([]byte, error) {
	data := map[string]interface{}{
		"params": map[string]interface{}{
			"gameId":     gameID,
			"categoryId": categoryID,
			"timer":      0,
			"video":      0,
			"obsolete":   0,
		},
		"page": pageNumber,
	}

	url, err := formatUrl(v2GameLeaderboardQuery, data)
	if err != nil {
		return nil, err
	}

	return c.client.Get(url)
}

func (c *SrcomV2Client) GetGameCategoryVariableValueLeaderboard(
	gameID, categoryID, variableID, valueID string,
	pageNumber int,
) ([]byte, error) {
	data := map[string]interface{}{
		"params": map[string]interface{}{
			"gameId":     gameID,
			"categoryId": categoryID,
			"values": []map[string]interface{}{
				{
					"variableId": variableID,
					"valueIds":   []string{valueID},
				},
			},
			"timer":    0,
			"video":    0,
			"obsolete": 0,
		},
		"page": pageNumber,
	}

	url, err := formatUrl(v2GameLeaderboardQuery, data)
	if err != nil {
		return nil, err
	}

	return c.client.Get(url)
}

func (c *SrcomV2Client) GetGameCategoryWorldRecordHistory(gameID, categoryID string) ([]byte, error) {
	data := map[string]interface{}{
		"params": map[string]interface{}{
			"gameId":     gameID,
			"categoryId": categoryID,
		},
	}

	url, err := formatUrl(v2GameRecordHistoryQuery, data)
	if err != nil {
		return nil, err
	}

	return c.client.Get(url)
}

func (c *SrcomV2Client) GetGameCategoryLevelWorldRecordHistory(gameID, categoryID, levelID string) ([]byte, error) {
	data := map[string]interface{}{
		"params": map[string]interface{}{
			"gameId":     gameID,
			"categoryId": categoryID,
			"levelId":    levelID,
		},
	}

	url, err := formatUrl(v2GameRecordHistoryQuery, data)
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
