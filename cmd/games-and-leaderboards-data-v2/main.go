package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"os"

	"github.com/alexmerren/speedruncom-scraper/internal"
	"github.com/buger/jsonparser"
)

func main() {
	if err := generateHybridData(); err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		os.Exit(1)
	}
}

func generateHybridData() error {
	client := internal.NewSrcomV1Client()

	gamesIdListFile, closeFunc, _ := internal.NewCsvReader(internal.GamesIdListFilenameV1)
	defer closeFunc()

	leaderboardsDataFile, closeFunc, _ := internal.NewCsvWriter(internal.LeaderboardsDataFilenameV1)
	leaderboardsDataFile.Write(internal.FileHeaders[internal.LeaderboardsDataFilenameV1])
	defer closeFunc()

	gamesDataFile, closeFunc, _ := internal.NewCsvWriter(internal.GamesDataFilenameV1)
	gamesDataFile.Write(internal.FileHeaders[internal.GamesDataFilenameV1])
	defer closeFunc()

	categoriesDataFile, closeFunc, _ := internal.NewCsvWriter(internal.CategoriesDataFilenameV1)
	categoriesDataFile.Write(internal.FileHeaders[internal.CategoriesDataFilenameV1])
	defer closeFunc()

	levelsDataFile, closeFunc, _ := internal.NewCsvWriter(internal.LevelsDataFilenameV1)
	levelsDataFile.Write(internal.FileHeaders[internal.LevelsDataFilenameV1])
	defer closeFunc()

	variablesDataFile, closeFunc, _ := internal.NewCsvWriter(internal.VariablesDataFilenameV1)
	variablesDataFile.Write(internal.FileHeaders[internal.VariablesDataFilenameV1])
	defer closeFunc()

	valuesDataFile, closeFunc, _ := internal.NewCsvWriter(internal.ValuesDataFilenameV1)
	valuesDataFile.Write(internal.FileHeaders[internal.ValuesDataFilenameV1])
	defer closeFunc()

	gamesIdListFile.Read()

	for {
		record, err := gamesIdListFile.Read()
		if err != nil && errors.Is(err, io.EOF) {
			break
		}
		gameId := record[0]
		gameResponse, err := client.GetGame(gameId)
		if err != nil {
			continue
		}

		// // Process categories, levels, variables, values, and game (taken from cmd/games-data/main.go)
		// numCategories, err := processCategory(categoriesDataFile, gameResponse, gameId)
		// if err != nil {
		// 	return err
		// }

		// numLevels, err := processLevel(levelsDataFile, gameResponse, gameId)
		// if err != nil {
		// 	return err
		// }

		// err = processVariableValue(variablesDataFile, valuesDataFile, gameResponse, gameId)
		// if err != nil {
		// 	return err
		// }

		// err = processGame(gamesDataFile, gameResponse, numCategories, numLevels, gameId)
		// if err != nil {
		// 	return err
		// }
		fmt.Println(gameId)
		leaderboard, err := generateLeaderboardMap(gameResponse)
		if err != nil {
			return err
		}

		t, _ := json.MarshalIndent(leaderboard, "", "\t")
		fmt.Println(string(t))

		// Process leaderboards for both categories and levels.
		// Iterate over leaderboardMap here.
	}

	return nil
}

// Processing functions taken from cmd/leaderboards-data/main.go
// func processLeaderboard(leaderboardsOutputFile *csv.Writer, leaderboardResponse []byte) error {
// 	_, err := jsonparser.ArrayEach(leaderboardResponse, func(value []byte, dataType jsonparser.ValueType, offset int, _ error) {
// 		place, _ := jsonparser.GetInt(value, "place")
// 		runData, _, _, _ := jsonparser.Get(value, "run")
// 		runId, _ := jsonparser.GetString(runData, "id")
// 		gameId, _ := jsonparser.GetString(runData, "game")
// 		categoryId, _ := jsonparser.GetString(runData, "category")
// 		levelId, _ := jsonparser.GetString(runData, "level")
// 		runDate, _ := jsonparser.GetString(runData, "date")
// 		runPrimaryTime, _ := jsonparser.GetFloat(runData, "times", "primary_t")
// 		runPlatform, _ := jsonparser.GetString(runData, "system", "platform")
// 		runEmulated, _ := jsonparser.GetBoolean(runData, "system", "emulated")
// 		runVerifiedDate, _ := jsonparser.GetString(runData, "status", "verify-date")
// 		runExaminer, _ := jsonparser.GetString(runData, "status", "examiner")

// 		playerIDArray := []string{}
// 		jsonparser.ArrayEach(runData, func(value []byte, dataType jsonparser.ValueType, offset int, err error) {
// 			playerID, _ := jsonparser.GetString(value, "id")
// 			playerIDArray = append(playerIDArray, string(playerID))
// 		}, "players")
// 		runPlayers := strings.Join(playerIDArray, ",")

// 		runValuesArray := []string{}
// 		jsonparser.ObjectEach(runData, func(key []byte, value []byte, dataType jsonparser.ValueType, offset int) error {
// 			runValuesArray = append(runValuesArray, fmt.Sprintf("%s=%s", string(key), string(value)))
// 			return nil
// 		}, "values")
// 		runValues := strings.Join(runValuesArray, ",")

// 		leaderboardsOutputFile.Write([]string{
// 			runId,
// 			gameId,
// 			categoryId,
// 			levelId,
// 			strconv.Itoa(int(place)),
// 			runDate,
// 			strconv.FormatFloat(runPrimaryTime, 'f', -1, 64),
// 			runPlatform,
// 			strconv.FormatBool(runEmulated),
// 			runPlayers,
// 			runExaminer,
// 			runVerifiedDate,
// 			runValues,
// 		})
// 	}, "data", "runs")

// 	return err
// }

type leaderboardMap map[string]map[string]map[string][]string
type leaderboardMapSub1 map[string]map[string][]string
type leaderboardMapSub2 map[string][]string

func generateLeaderboardMap(gameResponse []byte) (leaderboardMap, error) {
	result := make(leaderboardMap)
	result["per-level"] = make(leaderboardMapSub1)
	result["per-game"] = make(leaderboardMapSub1)

	_, err := jsonparser.ArrayEach(gameResponse, func(value []byte, dataType jsonparser.ValueType, offset int, err error) {
		categoryId, _ := jsonparser.GetString(value, "id")
		categoryType, _ := jsonparser.GetString(value, "type")
		result[categoryType][categoryId] = make(leaderboardMapSub2)
	}, "data", "categories", "data")
	if err != nil {
		return nil, err
	}

	variableScopeToCategoryType := map[string][]string{
		"global":       {"per-game", "per-level"},
		"full-game":    {"per-game"},
		"all-levels":   {"per-level"},
		"single-level": {},
	}

	_, err = jsonparser.ArrayEach(gameResponse, func(variableData []byte, dataType jsonparser.ValueType, offset int, err error) {
		variableId, _ := jsonparser.GetString(variableData, "id")
		valueIds := make([]string, 0)
		jsonparser.ObjectEach(variableData, func(key []byte, value []byte, dataType jsonparser.ValueType, offset int) error {
			valueIds = append(valueIds, string(key))
			return nil
		}, "values", "choices")

		categoryType, _ := jsonparser.GetString(variableData, "scope", "type")
		categoryId, err := jsonparser.GetString(variableData, "category")

		for _, categoryMapping := range variableScopeToCategoryType[categoryType] {
			categoryIds := make([]string, 0)
			if err != nil {
				for key := range result[categoryMapping] {
					categoryIds = append(categoryIds, key)
				}
			} else {
				categoryIds = append(categoryIds, categoryId)
			}

			for _, categoryId := range categoryIds {
				if _, ok := result[categoryMapping][categoryId]; ok {
					result[categoryMapping][categoryId][variableId] = valueIds
				}
			}
		}
	}, "data", "variables", "data")
	if err != nil {
		return nil, err
	}

	return result, nil
}

// Processing functions taken from cmd/games-data/main.go
// func processCategory(categoriesDataFile *csv.Writer, gameResponse []byte, gameId string) (int, error) {
// 	numCategories := 0
// 	_, err := jsonparser.ArrayEach(gameResponse, func(value []byte, dataType jsonparser.ValueType, offset int, err error) {
// 		numCategories += 1
// 		categoryId, _ := jsonparser.GetString(value, "id")
// 		categoryName, _ := jsonparser.GetString(value, "name")
// 		categoryRules, _ := jsonparser.GetString(value, "rules")
// 		categoryNumPlayers, _ := jsonparser.GetInt(value, "players", "value")
// 		categoryType, _ := jsonparser.GetString(value, "type")

// 		categoriesDataFile.Write([]string{
// 			gameId,
// 			categoryId,
// 			categoryName,
// 			internal.FormatCsvString(categoryRules),
// 			categoryType,
// 			strconv.Itoa(int(categoryNumPlayers)),
// 		})
// 	}, "data", "categories", "data")

// 	return numCategories, err
// }

// func processLevel(levelsDataFile *csv.Writer, gameResponse []byte, gameId string) (int, error) {
// 	numLevels := 0
// 	_, err := jsonparser.ArrayEach(gameResponse, func(value []byte, dataType jsonparser.ValueType, offset int, err error) {
// 		numLevels += 1
// 		levelId, _ := jsonparser.GetString(value, "id")
// 		levelName, _ := jsonparser.GetString(value, "name")
// 		levelRules, _ := jsonparser.GetString(value, "rules")

// 		levelsDataFile.Write([]string{gameId, levelId, levelName, internal.FormatCsvString(levelRules)})
// 	}, "data", "levels", "data")

// 	return numLevels, err
// }

// func processVariableValue(variablesDataFile, valuesDataFile *csv.Writer, gameResponse []byte, gameId string) error {
// 	_, err := jsonparser.ArrayEach(gameResponse, func(value []byte, dataType jsonparser.ValueType, offset int, _ error) {
// 		variableId, _ := jsonparser.GetString(value, "id")
// 		variableName, _ := jsonparser.GetString(value, "name")
// 		variableCategory, _ := jsonparser.GetString(value, "category")
// 		variableScope, _ := jsonparser.GetString(value, "scope", "type")
// 		variableIsSubcategory, _ := jsonparser.GetBoolean(value, "is-subcategory")
// 		variableDefault, _ := jsonparser.GetString(value, "values", "default")

// 		variablesDataFile.Write([]string{
// 			gameId,
// 			variableId,
// 			variableName,
// 			variableCategory,
// 			variableScope,
// 			strconv.FormatBool(variableIsSubcategory),
// 			variableDefault,
// 		})

// 		jsonparser.ObjectEach(value, func(key []byte, value []byte, dataType jsonparser.ValueType, offset int) error {
// 			valueId := string(key)
// 			valueLabel, _ := jsonparser.GetString(value, "label")
// 			valueRules, _ := jsonparser.GetString(value, "rules")

// 			valuesDataFile.Write([]string{
// 				gameId,
// 				variableId,
// 				valueId,
// 				valueLabel,
// 				internal.FormatCsvString(valueRules),
// 			})
// 			return nil
// 		}, "values", "values")
// 	}, "data", "variables", "data")

// 	return err
// }

// func processGame(gamesDataFile *csv.Writer, gameResponse []byte, numCategories, numLevels int, gameId string) error {
// 	gameData, _, _, err := jsonparser.Get(gameResponse, "data")
// 	if err != nil {
// 		return err
// 	}

// 	gameName, _ := jsonparser.GetString(gameData, "names", "international")
// 	gameURL, _ := jsonparser.GetString(gameData, "abbreviation")
// 	gameReleaseDate, _ := jsonparser.GetString(gameData, "release-date")
// 	gameCreatedDate, _ := jsonparser.GetString(gameData, "created")

// 	gamesDataFile.Write([]string{
// 		gameId,
// 		gameName,
// 		gameURL,
// 		gameReleaseDate,
// 		gameCreatedDate,
// 		strconv.Itoa(numCategories),
// 		strconv.Itoa(numLevels),
// 	})

// 	return nil
// }
