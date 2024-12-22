package srcom_api

import (
	"fmt"
)

const (
	apiUrl                            = "https://www.speedrun.com/api/v1/%s"
	gameQuery                         = "games/%s?embed=levels,categories,developers,platforms,genres,variables,publishers"
	categoryQuery                     = "categories/%s"
	levelQuery                        = "levels/%s"
	developerQuery                    = "developers/%s"
	leaderboardWithCategoryQuery      = "leaderboards/%s/category/%s"
	leaderboardWithLevelCategoryQuery = "leaderboards/%s/level/%s/%s"
	userQuery                         = "users/%s"
	runListQuery                      = "runs?user=%s&embed=players&max=%d&offset=%d"
	gameListQuery                     = "games?_bulk=yes&max=%d&orderby=released&direction=asc&offset=%d"
	platformListQuery                 = "platforms?max=%d&offset=%d"
	publisherListQuery                = "publishers?max=%d&offset=%d"
	genreListQuery                    = "genres?max=%d&offset=%d"
)

const (
	bulkEndpointMaximumPagination = 1000
	maximumPagination             = 200
)

type SrcomV1Client struct {
	client HttpClient
}

func NewSrcomV1Client(client HttpClient) *SrcomV1Client {
	return &SrcomV1Client{
		client: client,
	}
}

func (c *SrcomV1Client) GetUser(userId string) ([]byte, error) {
	requestPath := fmt.Sprintf(userQuery, userId)
	return c.get(requestPath)
}

func (c *SrcomV1Client) GetGame(gameId string) ([]byte, error) {
	requestPath := fmt.Sprintf(gameQuery, gameId)
	return c.get(requestPath)
}

func (c *SrcomV1Client) GetCategory(categoryId string) ([]byte, error) {
	requestPath := fmt.Sprintf(categoryQuery, categoryId)
	return c.get(requestPath)
}

func (c *SrcomV1Client) GetLevel(levelId string) ([]byte, error) {
	requestPath := fmt.Sprintf(levelQuery, levelId)
	return c.get(requestPath)
}

func (c *SrcomV1Client) GetDeveloper(developerId string) ([]byte, error) {
	requestPath := fmt.Sprintf(developerQuery, developerId)
	return c.get(requestPath)
}

func (c *SrcomV1Client) GetGameCategoryLeaderboard(gameId, categoryId string) ([]byte, error) {
	requestPath := fmt.Sprintf(leaderboardWithCategoryQuery,
		gameId, categoryId,
	)
	return c.get(requestPath)
}

func (c *SrcomV1Client) GetGameLevelCategoryLeaderboard(gameId, levelId, categoryId string) ([]byte, error) {
	requestPath := fmt.Sprintf(
		leaderboardWithLevelCategoryQuery,
		gameId,
		levelId,
		categoryId,
	)
	return c.get(requestPath)
}

func (c *SrcomV1Client) GetRunsByUser(userId string, pageNumber int) ([]byte, error) {
	requestPath := fmt.Sprintf(
		runListQuery,
		userId,
		maximumPagination,
		pageNumber*maximumPagination,
	)
	return c.get(requestPath)
}

func (c *SrcomV1Client) GetGameList(pageNumber int) ([]byte, error) {
	requestPath := fmt.Sprintf(gameListQuery, bulkEndpointMaximumPagination, pageNumber*bulkEndpointMaximumPagination)
	return c.get(requestPath)
}

func (c *SrcomV1Client) GetPlatformList(pageNumber int) ([]byte, error) {
	requestPath := fmt.Sprintf(platformListQuery, maximumPagination, pageNumber*maximumPagination)
	return c.get(requestPath)
}

func (c *SrcomV1Client) GetPublisherList(pageNumber int) ([]byte, error) {
	requestPath := fmt.Sprintf(publisherListQuery, maximumPagination, pageNumber*maximumPagination)
	return c.get(requestPath)
}

func (c *SrcomV1Client) GetGenreList(pageNumber int) ([]byte, error) {
	requestPath := fmt.Sprintf(genreListQuery, maximumPagination, pageNumber*maximumPagination)
	return c.get(requestPath)
}

func (c *SrcomV1Client) get(requestPath string) ([]byte, error) {
	url := fmt.Sprintf(apiUrl, requestPath)
	return c.client.Get(url)
}
