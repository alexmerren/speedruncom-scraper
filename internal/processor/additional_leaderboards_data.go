package processor

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/alexmerren/speedruncom-scraper/internal/repository"
	"github.com/alexmerren/speedruncom-scraper/internal/srcom_api"
	"github.com/alistanis/cartesian"
	"github.com/buger/jsonparser"
)

type AdditionalLeaderboardsDataProcessor struct {
	GameId                     string
	AdditionalLeaderboardsFile *repository.WriteRepository
	Client                     *srcom_api.SrcomV1Client
}

func (p *AdditionalLeaderboardsDataProcessor) Process() error {
	response, err := p.Client.GetGame(p.GameId)
	if err != nil {
		return err
	}

	leaderboardCombinations, err := generateLeaderboardCombinations(response)
	if err != nil {
		return err
	}

	for _, combination := range leaderboardCombinations {
		leaderboardResponse, err := p.Client.GetLeaderboardByVariables(combination.gameId, combination.categoryId, combination.levelId, combination.variableIds, combination.valueIds)
		if err != nil {
			return err
		}

		err = p.processLeaderboard(leaderboardResponse)
		if err != nil {
			return err
		}
	}

	return nil
}

type combination struct {
	gameId      string
	categoryId  string
	levelId     *string
	variableIds []string
	valueIds    []string
}

func (c *combination) isValid() bool {
	return len(c.variableIds) == len(c.valueIds)
}

// Use similar logic to [leaderboards_data.go] when parsing categories and levels.
// This can be quite opaque IMO, so in [adr_02.md] the logic is written in Python.
func generateLeaderboardCombinations(response []byte) ([]*combination, error) {
	gameId, err := jsonparser.GetString(response, "data", "id")
	if err != nil {
		return nil, err
	}

	variableData, _, _, err := jsonparser.Get(response, "data", "variables", "data")
	if err != nil {
		return nil, err
	}

	combinations := make([]*combination, 0)

	_, err = jsonparser.ArrayEach(response, func(value []byte, dataType jsonparser.ValueType, offset int, _ error) {
		categoryId, _ := jsonparser.GetString(value, "id")
		categoryType, _ := jsonparser.GetString(value, "type")

		if string(categoryType) == "per-game" {

			applicableVariables, applicableValues := findApplicableVariablesAndValues(variableData, categoryId, nil)
			newCombinations := computeCartesianProduct(gameId, categoryId, nil, applicableVariables, applicableValues)
			combinations = append(combinations, newCombinations...)

		}

		if string(categoryType) == "per-level" {

			_, err = jsonparser.ArrayEach(response, func(value []byte, dataType jsonparser.ValueType, offset int, _ error) {
				levelId, _ := jsonparser.GetString(value, "id")

				applicableVariables, applicableValues := findApplicableVariablesAndValues(variableData, categoryId, &levelId)
				newCombinations := computeCartesianProduct(gameId, categoryId, &levelId, applicableVariables, applicableValues)
				combinations = append(combinations, newCombinations...)

			}, "data", "levels", "data")
		}

	}, "data", "categories", "data")

	for _, combination := range combinations {
		if isValid := combination.isValid(); !isValid {
			return nil, fmt.Errorf("combination is invalid: %+v", combination)
		}
	}

	return combinations, nil
}

// findApplicableVariablesAndValues will apply the criteria from [variableIsApplicable]
// and [variableIsApplicableWithLevel] to determine if a variable (and it's corresponding
// values) should apply to the current category and level. If the variable matches the
// criteria, we should return a list of all the variables, and a list of list of
// the possible values of those variables. i.e. For applicable variables and their
// values such that the input is like:
//
//	{
//		variable1: [value1, value2],
//	 	variable2: [value3, value4]
//	}
//
// we should return [variable1, variable2] [[value1, value2], [value3, value4]]
func findApplicableVariablesAndValues(variablesData []byte, categoryId string, levelId *string) ([]string, [][]string) {
	variables := make([]string, 0)
	values := make([][]string, 0)

	_, err := jsonparser.ArrayEach(variablesData, func(value []byte, dataType jsonparser.ValueType, offset int, _ error) {
		if levelId == nil && !variableIsApplicable(value, categoryId) {
			return
		}

		if levelId != nil && !variableIsApplicableWithLevel(value, categoryId, *levelId) {
			return
		}

		variableId, _ := jsonparser.GetString(value, "id")
		valueIds := make([]string, 0)
		jsonparser.ObjectEach(value, func(key []byte, value []byte, dataType jsonparser.ValueType, offset int) error {
			valueIds = append(valueIds, string(key))
			return nil
		}, "values", "values")

		variables = append(variables, variableId)
		values = append(values, valueIds)
	})
	if err != nil {
		return nil, nil
	}

	return variables, values
}

func variableIsApplicableWithLevel(variableData []byte, categoryId string, levelId string) bool {
	isSubcategory, _ := jsonparser.GetBoolean(variableData, "is-subcategory")
	if !isSubcategory {
		return false
	}

	scopeType, _ := jsonparser.GetString(variableData, "scope", "type")
	if scopeType != "global" && scopeType != "full-game" && scopeType != "single-level" {
		return false
	}

	scopeLevel, _ := jsonparser.GetString(variableData, "scope", "level")
	if scopeType == "single-level" && scopeLevel != levelId {
		return false
	}

	category, _, _, err := jsonparser.Get(variableData, "category")
	if err != nil {
		return false
	}

	if category != nil && string(category) != categoryId {
		return false
	}

	return true
}

func variableIsApplicable(variableData []byte, categoryId string) bool {
	isSubcategory, _ := jsonparser.GetBoolean(variableData, "is-subcategory")
	if !isSubcategory {
		return false
	}

	scopeType, _ := jsonparser.GetString(variableData, "scope", "type")
	if scopeType != "global" && scopeType != "full-game" {
		return false
	}

	category, _, _, err := jsonparser.Get(variableData, "category")
	if err != nil {
		return false
	}

	if category != nil && string(category) != categoryId {
		return false
	}

	return true
}

// https://en.wikipedia.org/wiki/Cartesian_product
func computeCartesianProduct(gameId, categoryId string, levelId *string, variables []string, values [][]string) []*combination {
	nonEmptyValues := filter(values, func(valueIds []string) bool {
		return len(valueIds) > 0
	})

	// If no variables apply to a category/level combination, we have no values
	// to compute the cartesian product of. We can just return a single combination
	// of the game, category, and level ID (if applicable).
	if nonEmptyValues == nil {
		return []*combination{{
			gameId:      gameId,
			categoryId:  categoryId,
			levelId:     levelId,
			variableIds: nil,
			valueIds:    nil,
		}}
	}

	combinations := make([]*combination, 0)

	for _, product := range cartesian.Product(nonEmptyValues...) {
		combinations = append(combinations, &combination{
			gameId:      gameId,
			categoryId:  categoryId,
			levelId:     levelId,
			variableIds: variables,
			valueIds:    product,
		})
	}

	return combinations
}

// https://stackoverflow.com/questions/37562873/most-idiomatic-way-to-select-elements-from-an-array-in-golang
func filter[T any](ss []T, test func(T) bool) (result []T) {
	for _, s := range ss {
		if test(s) {
			result = append(result, s)
		}
	}
	return
}

func (p *AdditionalLeaderboardsDataProcessor) processLeaderboard(leaderboardResponse []byte) error {
	_, err := jsonparser.ArrayEach(leaderboardResponse, func(value []byte, dataType jsonparser.ValueType, offset int, _ error) {
		place, _ := jsonparser.GetInt(value, "place")
		runData, _, _, _ := jsonparser.Get(value, "run")
		runId, _ := jsonparser.GetString(runData, "id")
		gameId, _ := jsonparser.GetString(runData, "game")
		categoryId, _ := jsonparser.GetString(runData, "category")
		levelId, _ := jsonparser.GetString(runData, "level")
		runDate, _ := jsonparser.GetString(runData, "date")
		runPrimaryTime, _ := jsonparser.GetFloat(runData, "times", "primary_t")
		runPlatform, _ := jsonparser.GetString(runData, "system", "platform")
		runEmulated, _ := jsonparser.GetBoolean(runData, "system", "emulated")
		runVerifiedDate, _ := jsonparser.GetString(runData, "status", "verify-date")
		runExaminer, _ := jsonparser.GetString(runData, "status", "examiner")

		playerIDArray := []string{}
		jsonparser.ArrayEach(runData, func(value []byte, dataType jsonparser.ValueType, offset int, err error) {
			playerID, _ := jsonparser.GetString(value, "id")
			playerIDArray = append(playerIDArray, string(playerID))
		}, "players")
		runPlayers := strings.Join(playerIDArray, ",")

		runValuesArray := []string{}
		jsonparser.ObjectEach(runData, func(key []byte, value []byte, dataType jsonparser.ValueType, offset int) error {
			runValuesArray = append(runValuesArray, fmt.Sprintf("%s=%s", string(key), string(value)))
			return nil
		}, "values")
		runValues := strings.Join(runValuesArray, ",")

		p.AdditionalLeaderboardsFile.Write([]string{
			runId,
			gameId,
			categoryId,
			levelId,
			strconv.Itoa(int(place)),
			runDate,
			strconv.FormatFloat(runPrimaryTime, 'f', -1, 64),
			runPlatform,
			strconv.FormatBool(runEmulated),
			runPlayers,
			runExaminer,
			runVerifiedDate,
			runValues,
		})
	}, "data", "runs")

	return err
}
