package srcomv1

import (
	"fmt"

	"github.com/alexmerren/speedruncom-scraper/internal/srcomv2"
)

const (
	baseApiUrl                           = "https://www.speedrun.com/api/v1/%s"
	gameListFunction                     = "games?_bulk=yes&max=1000&orderby=released&direction=asc&offset=%d"
	gameFunction                         = "games/%s?embed=levels,categories,developers,platforms,genres,variables"
	categoryFunction                     = "categories/%s"
	levelFunction                        = "levels/%s"
	userFunction                         = "users/%s/personal-bests?embed=game,players"
	runsFunction                         = "runs?user=%s&embed=players&max=200&offset=%d"
	developerFunction                    = "developers/%s"
	gameCategoryLeaderboardFunction      = "leaderboards/%s/category/%s?embed=game,category,level,players,variables"
	gameCategoryLevelLeaderboardFunction = "leaderboards/%s/level/%s/%s"

	gameListNumberPerPage = 1000
	runsListNumberPerPage = 200
)

func GetUserRuns(userID string, pageNumber int) ([]byte, error) {
	header := fmt.Sprintf(runsFunction, userID, pageNumber*runsListNumberPerPage)
	URL := fmt.Sprintf(baseApiUrl, header)
	return srcomv2.RequestSrcom(URL)
}

func GetUser(userID string) ([]byte, error) {
	header := fmt.Sprintf(userFunction, userID)
	URL := fmt.Sprintf(baseApiUrl, header)
	return srcomv2.RequestSrcom(URL)
}

func GetGame(gameID string) ([]byte, error) {
	header := fmt.Sprintf(gameFunction, gameID)
	URL := fmt.Sprintf(baseApiUrl, header)
	return srcomv2.RequestSrcom(URL)
}

func GetCategory(categoryID string) ([]byte, error) {
	header := fmt.Sprintf(categoryFunction, categoryID)
	URL := fmt.Sprintf(baseApiUrl, header)
	return srcomv2.RequestSrcom(URL)
}

func GetLevel(levelID string) ([]byte, error) {
	header := fmt.Sprintf(levelFunction, levelID)
	URL := fmt.Sprintf(baseApiUrl, header)
	return srcomv2.RequestSrcom(URL)
}

func GetDeveloper(developerID string) ([]byte, error) {
	header := fmt.Sprintf(developerFunction, developerID)
	URL := fmt.Sprintf(baseApiUrl, header)
	return srcomv2.RequestSrcom(URL)
}

func GetGameList(pageNumber int) ([]byte, error) {
	header := fmt.Sprintf(gameListFunction, pageNumber*gameListNumberPerPage)
	URL := fmt.Sprintf(baseApiUrl, header)
	return srcomv2.RequestSrcom(URL)
}

func GetGameCategoryLeaderboard(gameID, categoryID string) ([]byte, error) {
	header := fmt.Sprintf(gameCategoryLeaderboardFunction,
		gameID, categoryID,
	)
	URL := fmt.Sprintf(baseApiUrl, header)
	return srcomv2.RequestSrcom(URL)
}

func GetGameCategoryLevelLeaderboard(gameID, categoryID, levelID string) ([]byte, error) {
	header := fmt.Sprintf(gameCategoryLevelLeaderboardFunction,
		gameID, levelID, categoryID,
	)
	URL := fmt.Sprintf(baseApiUrl, header)
	return srcomv2.RequestSrcom(URL)
}
