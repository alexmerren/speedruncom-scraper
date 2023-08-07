package srcomv1

import "fmt"

const (
	baseApiUrl                                   = "https://www.speedrun.com/api/v1/%s"
	gameListFunction                             = "games?_bulk=yes&max=1000&orderby=released&direction=asc&offset=%d"
	gameFunction                                 = "games/%s?embed=levels,cateogires,developers,platforms,genres,variables"
	categoryFunction                             = "categories/%s?embed=game,variables"
	levelFunction                                = "levels/%s?embed=categories,variables"
	userFunction                                 = "users/%s"
	developerFunction                            = "developers/%s"
	gameCategoryLeaderboardFunction              = "leaderboards/%s/category/%s?embed=game,category,level,players,variables"
	gameCategoryVariableValueLeaderboardFunction = "leaderboards/%s/category/%s?%s=%s"

	gameListNumberPerPage = 1000
)

func GetUser(userID string) ([]byte, error) {
	header := fmt.Sprintf(userFunction, userID)
	URL := fmt.Sprintf(baseApiUrl, header)
	return requestSrcom(URL)
}

func GetGame(gameID string) ([]byte, error) {
	header := fmt.Sprintf(gameFunction, gameID)
	URL := fmt.Sprintf(baseApiUrl, header)
	return requestSrcom(URL)
}

func GetCategory(categoryID string) ([]byte, error) {
	header := fmt.Sprintf(categoryFunction, categoryID)
	URL := fmt.Sprintf(baseApiUrl, header)
	return requestSrcom(URL)
}

func GetLevel(levelID string) ([]byte, error) {
	header := fmt.Sprintf(levelFunction, levelID)
	URL := fmt.Sprintf(baseApiUrl, header)
	return requestSrcom(URL)
}

func GetDeveloper(developerID string) ([]byte, error) {
	header := fmt.Sprintf(developerFunction, developerID)
	URL := fmt.Sprintf(baseApiUrl, header)
	return requestSrcom(URL)
}

func GetGameList(pageNumber int) ([]byte, error) {
	header := fmt.Sprintf(gameListFunction, pageNumber*gameListNumberPerPage)
	URL := fmt.Sprintf(baseApiUrl, header)
	return requestSrcom(URL)
}

func GetGameCategoryLeaderboard(gameID, categoryID string) ([]byte, error) {
	header := fmt.Sprintf(gameCategoryLeaderboardFunction,
		gameID, categoryID,
	)
	URL := fmt.Sprintf(baseApiUrl, header)
	return requestSrcom(URL)
}

func GetGameCategoryVariableValueLeaderboard(
	gameID, categoryID, variableID, valueID string,
	pageNumber int,
) ([]byte, error) {
	header := fmt.Sprintf(gameCategoryVariableValueLeaderboardFunction,
		gameID, categoryID, variableID, valueID,
	)
	URL := fmt.Sprintf(baseApiUrl, header)
	return requestSrcom(URL)
}
