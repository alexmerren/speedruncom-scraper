package srcomv2

const (
	baseApiUrl = "https://www.speedrun.com/api/v2/%s?_r=%s"

	searchFunction                               = "GetSearch"
	userDataFunction                             = "GetUserLeaderboard"
	userSummaryFunction                          = "GetUserSummary"
	gameDataFunction                             = "GetGameData"
	gameSummaryFunction                          = "GetGameSummary"
	gameListFunction                             = "GetGameList"
	gameCategoryLeaderboardFunction              = "GetGameLeaderboard2"
	gameCategoryVariableValueLeaderboardFunction = "GetGameLeaderboard2"
	gameCategoryWorldRecordHistory               = "GetGameRecordHistory"

	searchResultLimit = 15
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

	return RequestSrcom(URL)
}

func GetGameData(gameID string) ([]byte, error) {
	data := map[string]interface{}{
		"gameId": gameID,
	}

	URL, err := formatHeader(data, gameDataFunction)
	if err != nil {
		return nil, err
	}

	return RequestSrcom(URL)
}

func GetGameSummary(gameURL string) ([]byte, error) {
	data := map[string]interface{}{
		"gameUrl": gameURL,
	}

	URL, err := formatHeader(data, gameSummaryFunction)
	if err != nil {
		return nil, err
	}

	return RequestSrcom(URL)
}

func GetGameList(pageNumber int) ([]byte, error) {
	data := map[string]interface{}{
		"page": pageNumber,
	}

	URL, err := formatHeader(data, gameListFunction)
	if err != nil {
		return nil, err
	}

	return RequestSrcom(URL)
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

	URL, err := formatHeader(data, gameListFunction)
	if err != nil {
		return nil, err
	}

	return RequestSrcom(URL)
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

	URL, err := formatHeader(data, gameListFunction)
	if err != nil {
		return nil, err
	}

	return RequestSrcom(URL)
}

func GetGameCategoryWorldRecordHistory(gameID, categoryID string) ([]byte, error) {
	data := map[string]interface{}{
		"params": map[string]interface{}{
			"gameId":     gameID,
			"categoryId": categoryID,
		},
	}

	URL, err := formatHeader(data, gameListFunction)
	if err != nil {
		return nil, err
	}

	return RequestSrcom(URL)
}

func GetGameCategoryLevelWorldRecordHistory(gameID, categoryID, levelID string) ([]byte, error) {
	data := map[string]interface{}{
		"params": map[string]interface{}{
			"gameId":     gameID,
			"categoryId": categoryID,
			"levelId":    levelID,
		},
	}

	URL, err := formatHeader(data, gameListFunction)
	if err != nil {
		return nil, err
	}

	return RequestSrcom(URL)
}

func GetSearch(query string) ([]byte, error) {
	data := map[string]interface{}{
		"query":         query,
		"limit":         searchResultLimit,
		"includeGames":  true,
		"includeNews":   true,
		"includePages":  true,
		"includeSeries": true,
		"includeUsers":  true,
	}

	URL, err := formatHeader(data, searchFunction)
	if err != nil {
		return nil, err
	}

	return RequestSrcom(URL)
}
