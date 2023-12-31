package srcomv2

import "github.com/alexmerren/speedruncom-scraper/pkg/requests"

const (
	baseApiUrl = "https://www.speedrun.com/api/v2/%s?_r=%s"

	userDataFunction               = "GetUserLeaderboard"
	userSummaryFunction            = "GetUserSummary"
	gameDataFunction               = "GetGameData"
	gameSummaryFunction            = "GetGameSummary"
	gameListFunction               = "GetGameList"
	gameLeaderboardFunction        = "GetGameLeaderboard2"
	gameWorldRecordHistoryFunction = "GetGameRecordHistory"
)

func GetUserData(userID string) ([]byte, error) {
	data := map[string]interface{}{
		"userId":    userID,
		"levelType": 1,
	}

	URL, err := formatHeader(data, userDataFunction)
	if err != nil {
		return nil, err
	}

	return requests.RequestSrcom(URL)
}

func GetGameData(gameID string) ([]byte, error) {
	data := map[string]interface{}{
		"gameId": gameID,
	}

	URL, err := formatHeader(data, gameDataFunction)
	if err != nil {
		return nil, err
	}

	return requests.RequestSrcom(URL)
}

func GetGameSummary(gameURL string) ([]byte, error) {
	data := map[string]interface{}{
		"gameUrl": gameURL,
	}

	URL, err := formatHeader(data, gameSummaryFunction)
	if err != nil {
		return nil, err
	}

	return requests.RequestSrcom(URL)
}

func GetGameList(pageNumber int) ([]byte, error) {
	data := map[string]interface{}{
		"page": pageNumber,
	}

	URL, err := formatHeader(data, gameListFunction)
	if err != nil {
		return nil, err
	}

	return requests.RequestSrcom(URL)
}

func GetGameCategoryLeaderboard(gameID, categoryID string, pageNumber int) ([]byte, error) {
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

	URL, err := formatHeader(data, gameLeaderboardFunction)
	if err != nil {
		return nil, err
	}

	return requests.RequestSrcom(URL)
}

func GetGameCategoryVariableValueLeaderboard(
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

	URL, err := formatHeader(data, gameLeaderboardFunction)
	if err != nil {
		return nil, err
	}

	return requests.RequestSrcom(URL)
}

func GetGameCategoryWorldRecordHistory(gameID, categoryID string) ([]byte, error) {
	data := map[string]interface{}{
		"params": map[string]interface{}{
			"gameId":     gameID,
			"categoryId": categoryID,
		},
	}

	URL, err := formatHeader(data, gameWorldRecordHistoryFunction)
	if err != nil {
		return nil, err
	}

	return requests.RequestSrcom(URL)
}

func GetGameCategoryLevelWorldRecordHistory(gameID, categoryID, levelID string) ([]byte, error) {
	data := map[string]interface{}{
		"params": map[string]interface{}{
			"gameId":     gameID,
			"categoryId": categoryID,
			"levelId":    levelID,
		},
	}

	URL, err := formatHeader(data, gameWorldRecordHistoryFunction)
	if err != nil {
		return nil, err
	}

	return requests.RequestSrcom(URL)
}
