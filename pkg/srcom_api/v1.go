package srcom_api

import (
	"fmt"
)

const (
	v1ApiUrl                                = "https://www.speedrun.com/api/v1/%s"
	v1GameListQuery                         = "games?_bulk=yes&max=1000&orderby=released&direction=asc&offset=%d"
	v1GameDataQuery                         = "games/%s?embed=levels,categories,developers,platforms,genres,variables"
	v1GameCategoryQuery                     = "categories/%s"
	v1GameLevelQuery                        = "levels/%s"
	v1GameDeveloperQuery                    = "developers/%s"
	v1GameLeaderboardWithCategoryQuery      = "leaderboards/%s/category/%s?embed=game,category,level,players,variables"
	v1GameLeaderboardWithLevelCategoryQuery = "leaderboards/%s/level/%s/%s"
	v1UserDataQuery                         = "users/%s"
	v1UserRunsQuery                         = "runs?user=%s&embed=players&max=%d&offset=%d"
)

const (
	v1GameListMaxPerPage = 1000
	v1RunsListMaxPerPage = 200
)

type SrcomV1Client struct {
	client HttpClient
}

func NewSrcomV1Client(client HttpClient) *SrcomV1Client {
	return &SrcomV1Client{
		client: client,
	}
}

func (c *SrcomV1Client) GetRunsByUser(userId string, pageNumber int) ([]byte, error) {
	header := fmt.Sprintf(
		v1UserRunsQuery,
		userId,
		v1RunsListMaxPerPage,
		pageNumber*v1RunsListMaxPerPage,
	)
	url := fmt.Sprintf(v1ApiUrl, header)

	return c.client.Get(url)
}

func (c *SrcomV1Client) GetUser(userId string) ([]byte, error) {
	header := fmt.Sprintf(v1UserDataQuery, userId)
	url := fmt.Sprintf(v1ApiUrl, header)

	return c.client.Get(url)
}

func (c *SrcomV1Client) GetGame(gameId string) ([]byte, error) {
	header := fmt.Sprintf(v1GameDataQuery, gameId)
	url := fmt.Sprintf(v1ApiUrl, header)

	return c.client.Get(url)
}

func (c *SrcomV1Client) GetCategory(categoryId string) ([]byte, error) {
	header := fmt.Sprintf(v1GameCategoryQuery, categoryId)
	url := fmt.Sprintf(v1ApiUrl, header)

	return c.client.Get(url)
}

func (c *SrcomV1Client) GetLevel(levelId string) ([]byte, error) {
	header := fmt.Sprintf(v1GameLevelQuery, levelId)
	url := fmt.Sprintf(v1ApiUrl, header)

	return c.client.Get(url)
}

func (c *SrcomV1Client) GetDeveloper(developerId string) ([]byte, error) {
	header := fmt.Sprintf(v1GameDeveloperQuery, developerId)
	url := fmt.Sprintf(v1ApiUrl, header)

	return c.client.Get(url)
}

func (c *SrcomV1Client) GetGameList(pageNumber int) ([]byte, error) {
	header := fmt.Sprintf(v1GameListQuery, pageNumber*v1GameListMaxPerPage)
	url := fmt.Sprintf(v1ApiUrl, header)

	return c.client.Get(url)
}

func (c *SrcomV1Client) GetGameCategoryLeaderboard(gameId, categoryId string) ([]byte, error) {
	header := fmt.Sprintf(v1GameLeaderboardWithCategoryQuery,
		gameId, categoryId,
	)
	url := fmt.Sprintf(v1ApiUrl, header)

	return c.client.Get(url)
}

func (c *SrcomV1Client) GetGameLevelCategoryLeaderboard(gameId, levelId, categoryId string) ([]byte, error) {
	header := fmt.Sprintf(
		v1GameLeaderboardWithLevelCategoryQuery,
		gameId,
		levelId,
		categoryId,
	)
	url := fmt.Sprintf(v1ApiUrl, header)

	return c.client.Get(url)
}
