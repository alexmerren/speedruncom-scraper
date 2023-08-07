package srcomv2

import (
	"fmt"
	"time"
)

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
)

func GetUserData(userID string, data []byte) ([]byte, error) {
	return requestSrcom("")
}

func GetGameData(gameID string, data []byte) ([]byte, error) {
	return requestSrcom("")
}

func GetGameSummary(gameURL string) ([]byte, error) {
	return requestSrcom("")
}

func GetGameList(pageNumber int) ([]byte, error) {
	data := map[string]interface{}{
		"page": pageNumber,
		"vary": time.Now().Unix(),
	}

	flattenedData, err := getBytes(data)
	if err != nil {
		return nil, err
	}

	formattedData := EncodeB64Header(flattenedData)
	URL := fmt.Sprintf(baseApiUrl, gameListFunction, formattedData)

	return requestSrcom(URL)
}

func GetGameCategoryLeaderboard(gameID, categoryID string, pageNumber int) ([]byte, error) {
	return requestSrcom("")
}

func GetGameCategoryVariableValueLeaderboard(
	gameID, categoryID, variableID, valueID string,
	pageNumber int,
) ([]byte, error) {
	return requestSrcom("")
}

func GetGameCategoryWorldRecordHistory(gameID, categoryID string) ([]byte, error) {
	return requestSrcom("")
}

func GetSearch(query string) ([]byte, error) {
	data := map[string]interface{}{
		"query":         query,
		"limit":         5,
		"includeGames":  true,
		"includeNews":   true,
		"includePages":  true,
		"includeSeries": true,
		"includeUsers":  true,
	}

	flattenedData, err := getBytes(data)
	if err != nil {
		return nil, err
	}

	formattedData := EncodeB64Header(flattenedData)
	URL := fmt.Sprintf(baseApiUrl, searchFunction, formattedData)

	return requestSrcom(URL)
}
