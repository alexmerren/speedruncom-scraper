package main

import (
	"os"
)

func main() {
	os.Exit(1)
	// if err := generateHybridData(); err != nil {
	// 	fmt.Fprintf(os.Stderr, "error: %v\n", err)
	// 	os.Exit(1)
	// }
}

// func generateHybridData() error {
// 	client := internal.NewSrcomV1Client()

// 	gamesIdListFile, closeFunc, _ := internal.NewCsvReader(internal.GamesIdListFilenameV1)
// 	defer closeFunc()

// 	leaderboardsDataFile, closeFunc, _ := internal.NewCsvWriter(internal.LeaderboardsDataFilenameV1)
// 	leaderboardsDataFile.Write(internal.FileHeaders[internal.LeaderboardsDataFilenameV1])
// 	defer closeFunc()

// 	gamesDataFile, closeFunc, _ := internal.NewCsvWriter(internal.GamesDataFilenameV1)
// 	gamesDataFile.Write(internal.FileHeaders[internal.GamesDataFilenameV1])
// 	defer closeFunc()

// 	categoriesDataFile, closeFunc, _ := internal.NewCsvWriter(internal.CategoriesDataFilenameV1)
// 	categoriesDataFile.Write(internal.FileHeaders[internal.CategoriesDataFilenameV1])
// 	defer closeFunc()

// 	levelsDataFile, closeFunc, _ := internal.NewCsvWriter(internal.LevelsDataFilenameV1)
// 	levelsDataFile.Write(internal.FileHeaders[internal.LevelsDataFilenameV1])
// 	defer closeFunc()

// 	variablesDataFile, closeFunc, _ := internal.NewCsvWriter(internal.VariablesDataFilenameV1)
// 	variablesDataFile.Write(internal.FileHeaders[internal.VariablesDataFilenameV1])
// 	defer closeFunc()

// 	valuesDataFile, closeFunc, _ := internal.NewCsvWriter(internal.ValuesDataFilenameV1)
// 	valuesDataFile.Write(internal.FileHeaders[internal.ValuesDataFilenameV1])
// 	defer closeFunc()

// 	gamesIdListFile.Read()

// 	for {
// 		record, err := gamesIdListFile.Read()
// 		if err != nil && errors.Is(err, io.EOF) {
// 			break
// 		}
// 		gameId := record[0]
// 		gameResponse, err := client.GetGame(gameId)
// 		if err != nil {
// 			continue
// 		}

// 		// Process categories, levels, variables, values, and game (taken from cmd/games-data/main.go)
// 		numCategories, err := processCategory(categoriesDataFile, response, gameId)
// 		if err != nil {
// 			return err
// 		}

// 		numLevels, err := processLevel(levelsDataFile, response, gameId)
// 		if err != nil {
// 			return err
// 		}

// 		err = processVariableValue(variablesDataFile, valuesDataFile, response, gameId)
// 		if err != nil {
// 			return err
// 		}

// 		err = processGame(gamesDataFile, response, numCategories, numLevels, gameId)
// 		if err != nil {
// 			return err
// 		}

// 		// Process leaderboards for both categories and levels.

// 	}

// 	return nil
// }

// // Processing functions taken from cmd/games-data/main.go
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
