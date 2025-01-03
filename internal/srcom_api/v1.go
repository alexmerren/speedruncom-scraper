package srcom_api

import (
	"fmt"
	"strings"
)

const (
	apiUrl                          = "https://www.speedrun.com/api/v1/%s"
	gameQuery                       = "games/%s?embed=levels,categories,variables"
	categoryQuery                   = "categories/%s"
	levelQuery                      = "levels/%s"
	leaderboardByCategoryQuery      = "leaderboards/%s/category/%s"
	leaderboardByLevelCategoryQuery = "leaderboards/%s/level/%s/%s"
	userQuery                       = "users/%s"
	runListQuery                    = "runs?user=%s&embed=players&max=%d&offset=%d"
	gameListQuery                   = "games?_bulk=yes&max=%d&orderby=released&direction=asc&offset=%d"
	platformListQuery               = "platforms?max=%d&offset=%d"
	publisherListQuery              = "publishers?max=%d&offset=%d"
	genreListQuery                  = "genres?max=%d&offset=%d"
	developerListQuery              = "developers?max=%d&offset=%d"
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
	requestQuery := fmt.Sprintf(userQuery, userId)
	return c.get(requestQuery)
}

func (c *SrcomV1Client) GetGame(gameId string) ([]byte, error) {
	requestQuery := fmt.Sprintf(gameQuery, gameId)
	return c.get(requestQuery)
}

func (c *SrcomV1Client) GetCategory(categoryId string) ([]byte, error) {
	requestQuery := fmt.Sprintf(categoryQuery, categoryId)
	return c.get(requestQuery)
}

func (c *SrcomV1Client) GetLevel(levelId string) ([]byte, error) {
	requestQuery := fmt.Sprintf(levelQuery, levelId)
	return c.get(requestQuery)
}

func (c *SrcomV1Client) GetLeaderboardByGameCategory(gameId, categoryId string) ([]byte, error) {
	requestQuery := fmt.Sprintf(leaderboardByCategoryQuery,
		gameId, categoryId,
	)
	return c.get(requestQuery)
}

func (c *SrcomV1Client) GetLeaderboardByGameLevelCategory(gameId, levelId, categoryId string) ([]byte, error) {
	requestQuery := fmt.Sprintf(
		leaderboardByLevelCategoryQuery,
		gameId,
		levelId,
		categoryId,
	)
	return c.get(requestQuery)
}

func (c *SrcomV1Client) GetLeaderboardByVariables(gameId, categoryId string, levelId *string, variables, values []string) ([]byte, error) {
	var requestQuery strings.Builder

	baseQuery := generateLeaderboardByVariablesRequestQuery(gameId, categoryId, levelId)
	if _, err := requestQuery.WriteString(baseQuery); err != nil {
		return nil, err
	}

	queryParametersArray := make([]string, len(variables))
	for index := range len(variables) {
		queryParametersArray[index] = fmt.Sprintf("var-%s=%s", variables[index], values[index])
	}

	_, err := requestQuery.WriteString("?" + strings.Join(queryParametersArray, "&"))
	if err != nil {
		return nil, err
	}

	return c.get(requestQuery.String())
}

func generateLeaderboardByVariablesRequestQuery(gameId, categoryId string, levelId *string) string {
	if levelId != nil {
		return fmt.Sprintf(
			leaderboardByLevelCategoryQuery,
			gameId,
			*levelId,
			categoryId,
		)
	}

	return fmt.Sprintf(leaderboardByCategoryQuery,
		gameId, categoryId,
	)
}

func (c *SrcomV1Client) GetRunsByUser(userId string, pageNumber int) ([]byte, error) {
	requestQuery := fmt.Sprintf(
		runListQuery,
		userId,
		maximumPagination,
		pageNumber*maximumPagination,
	)
	return c.get(requestQuery)
}

func (c *SrcomV1Client) GetGameList(pageNumber int) ([]byte, error) {
	requestQuery := fmt.Sprintf(gameListQuery, bulkEndpointMaximumPagination, pageNumber*bulkEndpointMaximumPagination)
	return c.get(requestQuery)
}

func (c *SrcomV1Client) GetPlatformList(pageNumber int) ([]byte, error) {
	requestQuery := fmt.Sprintf(platformListQuery, maximumPagination, pageNumber*maximumPagination)
	return c.get(requestQuery)
}

func (c *SrcomV1Client) GetPublisherList(pageNumber int) ([]byte, error) {
	requestQuery := fmt.Sprintf(publisherListQuery, maximumPagination, pageNumber*maximumPagination)
	return c.get(requestQuery)
}

func (c *SrcomV1Client) GetGenreList(pageNumber int) ([]byte, error) {
	requestQuery := fmt.Sprintf(genreListQuery, maximumPagination, pageNumber*maximumPagination)
	return c.get(requestQuery)
}

func (c *SrcomV1Client) GetDeveloperList(pageNumber int) ([]byte, error) {
	requestQuery := fmt.Sprintf(developerListQuery, maximumPagination, pageNumber*maximumPagination)
	return c.get(requestQuery)
}

func (c *SrcomV1Client) get(requestQuery string) ([]byte, error) {
	url := fmt.Sprintf(apiUrl, requestQuery)
	return c.client.Get(url)
}
