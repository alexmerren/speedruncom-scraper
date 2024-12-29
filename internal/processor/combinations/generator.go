package combinations

import (
	"fmt"

	"github.com/alistanis/cartesian"
	"github.com/buger/jsonparser"
)

// GenerateLeaderboardCombinations returns a list of valid category/level/variables/values
// combinations for a given game ID. The function requires the response of [SrcomV1Client.GetGame].
func GenerateLeaderboardCombinations(response []byte) ([]*Combination, error) {
	gameId, err := jsonparser.GetString(response, "data", "id")
	if err != nil {
		return nil, err
	}

	variableData, _, _, err := jsonparser.Get(response, "data", "variables", "data")
	if err != nil {
		return nil, err
	}

	combinations := make([]*Combination, 0)

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

		// This is a dumb way to get the keys of the map, but the API model is
		// forcing my hand here.
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
func computeCartesianProduct(gameId, categoryId string, levelId *string, variables []string, values [][]string) []*Combination {
	nonEmptyValues := filter(values, func(valueIds []string) bool {
		return len(valueIds) > 0
	})

	// If no variables apply to a category/level combination, we have no values
	// to compute the cartesian product of. We can just return a single combination
	// of the game, category, and level ID (if applicable).
	if nonEmptyValues == nil {
		return []*Combination{{
			GameId:      gameId,
			CategoryId:  categoryId,
			LevelId:     levelId,
			VariableIds: nil,
			ValueIds:    nil,
		}}
	}

	combinations := make([]*Combination, 0)

	for _, product := range cartesian.Product(nonEmptyValues...) {
		combinations = append(combinations, &Combination{
			GameId:      gameId,
			CategoryId:  categoryId,
			LevelId:     levelId,
			VariableIds: variables,
			ValueIds:    product,
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
