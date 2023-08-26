package main

import (
	"bufio"
	"fmt"

	"github.com/alexmerren/speedruncom-scraper/internal/filesystem"
	"github.com/alexmerren/speedruncom-scraper/internal/srcomv1"
	"github.com/alexmerren/speedruncom-scraper/internal/srcomv2"
	"github.com/buger/jsonparser"
)

const (
	allGameIDListV1          = "./data/v1/games-id-list.csv"
	gameOutputFilenameV1     = "./data/v1/games-data.csv"
	categoryOutputFilenameV1 = "./data/v1/categories-data.csv"
	levelOutputFilenameV1    = "./data/v1/level-data.csv"
	variableOutputFileV1     = "./data/v1/variable-data.csv"
	valueOutputFileV1        = "./data/v1/value-data.csv"

	allGameIDListV2          = "./data/v2/games-id-list.csv"
	gameOutputFilenameV2     = "./data/v2/games-data.csv"
	categoryOutputFilenameV2 = "./data/v2/categories-data.csv"
	levelOutputFilenameV2    = "./data/v2/level-data.csv"
	variableOutputFileV2     = "./data/v2/variable-data.csv"
	valueOutputFileV2        = "./data/v2/value-data.csv"
)

func main() {
	getGameDataV1()
}

//nolint:errcheck// Not worth checking for an error for every file write.
func getGameDataV1() {
	inputFile, err := filesystem.OpenInputFile(allGameIDListV1)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer inputFile.Close()

	gameOuptutFile, err := filesystem.CreateOutputFile(gameOutputFilenameV1)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer gameOuptutFile.Close()
	gameOuptutFile.WriteString("#ID,name,URL,releaseDate,createdDate,numCategories,numLevels\n")

	categoryOutputFile, err := filesystem.CreateOutputFile(categoryOutputFilenameV1)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer categoryOutputFile.Close()
	categoryOutputFile.WriteString("#parentGameID,ID,name,rules\n")

	levelOutputFile, err := filesystem.CreateOutputFile(levelOutputFilenameV1)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer levelOutputFile.Close()
	levelOutputFile.WriteString("#parentGameID,ID,name,rules\n")

	variableOutputFile, err := filesystem.CreateOutputFile(variableOutputFileV1)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer variableOutputFile.Close()
	variableOutputFile.WriteString("#parentGameID,ID,name,category,scope,isSubcategory,defaultValue\n")

	valueOutputFile, err := filesystem.CreateOutputFile(valueOutputFileV1)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer valueOutputFile.Close()
	valueOutputFile.WriteString("#parentGameID,variableID,ID,label,rules\n")

	// Scan the input file and get information for each of the game ID's in the
	// input file. We progress to the next line using scanner.Scan()
	scanner := bufio.NewScanner(inputFile)
	scanner.Scan()
	for scanner.Scan() {
		gameID := scanner.Text()
		response, err := srcomv1.GetGame(gameID)
		if err != nil {
			fmt.Println(err)
			continue
		}

		// Step 1. Process each category for a game
		numCategories := 0
		_, err = jsonparser.ArrayEach(response, func(value []byte, dataType jsonparser.ValueType, offset int, err error) {
			numCategories += 1
			categoryID, _, _, _ := jsonparser.Get(value, "id")
			categoryName, _, _, _ := jsonparser.Get(value, "name")
			categoryRules, _, _, _ := jsonparser.Get(value, "rules")
			categoryNumPlayers, _ := jsonparser.GetInt(value, "players", "value")
			categoryType, _, _, _ := jsonparser.Get(value, "type")
			categoryOutputFile.WriteString(fmt.Sprintf("%s,%s,\"%s\",\"%s\",%s,%d\n", gameID, categoryID, categoryName, categoryRules, categoryType, categoryNumPlayers))
		}, "data", "categories", "data")
		if err != nil {
			return
		}

		// Step 2. Process each level for a game
		numLevels := 0
		_, err = jsonparser.ArrayEach(response, func(value []byte, dataType jsonparser.ValueType, offset int, err error) {
			numLevels += 1
			levelID, _, _, _ := jsonparser.Get(value, "id")
			levelName, _, _, _ := jsonparser.Get(value, "name")
			levelRules, _, _, _ := jsonparser.Get(value, "rules")
			levelOutputFile.WriteString(fmt.Sprintf("%s,%s,\"%s\",\"%s\"\n", gameID, levelID, levelName, levelRules))
		}, "data", "levels", "data")
		if err != nil {
			return
		}

		// Step 3. Process each variable/value for a game
		_, err = jsonparser.ArrayEach(response, func(value []byte, dataType jsonparser.ValueType, offset int, err error) {
			variableID, _, _, _ := jsonparser.Get(value, "id")
			variableName, _, _, _ := jsonparser.Get(value, "name")
			variableCategory, _, _, _ := jsonparser.Get(value, "category")
			variableScope, _, _, _ := jsonparser.Get(value, "scope", "type")
			variableIsSubcategory, _ := jsonparser.GetBoolean(value, "is-subcategory")
			variableDefault, _, _, _ := jsonparser.Get(value, "values", "default")
			variableOutputFile.WriteString(fmt.Sprintf("%s,%s,\"%s\",%s,%s,%t,%s\n", gameID, variableID, variableName, variableCategory, variableScope, variableIsSubcategory, variableDefault))

			err = jsonparser.ObjectEach(value, func(key []byte, value []byte, dataType jsonparser.ValueType, offset int) error {
				valueID := string(key)
				valueLabel, _, _, _ := jsonparser.Get(value, "label")
				valueRules, _, _, _ := jsonparser.Get(value, "rules")
				valueOutputFile.WriteString(fmt.Sprintf("%s,%s,%s,\"%s\",\"%s\"\n", gameID, variableID, valueID, valueLabel, valueRules))
				return nil
			}, "values", "values")
			if err != nil {
				return
			}
		}, "data", "variables", "data")
		if err != nil {
			return
		}

		// Step N. Process each game
		gameName, _, _, _ := jsonparser.Get(response, "data", "names", "international")
		gameURL, _, _, _ := jsonparser.Get(response, "data", "abbreviation")
		gameReleaseDate, _, _, _ := jsonparser.Get(response, "data", "release-date")
		gameCreatedDate, _, _, _ := jsonparser.Get(response, "data", "created")
		gameOuptutFile.WriteString(fmt.Sprintf("%s,\"%s\",%s,%s,%s,%d,%d\n", gameID, gameName, gameURL, gameReleaseDate, gameCreatedDate, numCategories, numLevels))
	}

	if err := scanner.Err(); err != nil {
		fmt.Println(err)
		return
	}
}

//nolint:errcheck// Not worth checking for an error for every file write.
func getGameDataV2() {
	inputFile, err := filesystem.OpenInputFile(allGameIDListV2)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer inputFile.Close()

	gameOuptutFile, err := filesystem.CreateOutputFile(gameOutputFilenameV2)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer gameOuptutFile.Close()
	gameOuptutFile.WriteString("#ID,name,URL,type,rules,releaseDate,addedDate,runCount,playerCount,numCategories,numLevels,emulator\n")

	categoryOutputFile, err := filesystem.CreateOutputFile(categoryOutputFilenameV2)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer categoryOutputFile.Close()
	categoryOutputFile.WriteString("#parentGameID,ID,name,rules,numPlayers\n")

	levelOutputFile, err := filesystem.CreateOutputFile(levelOutputFilenameV2)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer levelOutputFile.Close()
	levelOutputFile.WriteString("#parentGameID,ID,name,rules,numPlayers\n")

	variableOutputFile, err := filesystem.CreateOutputFile(variableOutputFileV2)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer variableOutputFile.Close()
	variableOutputFile.WriteString("#parentGameID,ID,name,category,scope\n")

	valueOutputFile, err := filesystem.CreateOutputFile(valueOutputFileV2)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer valueOutputFile.Close()
	valueOutputFile.WriteString("#parentGameID,variableID,ID,label,rules\n")

	// Scan the input file and get information for each of the game ID's in the
	// input file. We progress to the next line using scanner.Scan()
	scanner := bufio.NewScanner(inputFile)
	scanner.Scan()
	for scanner.Scan() {
		response, err := srcomv2.GetGameData(scanner.Text())
		if err != nil {
			continue
		}

		gameID, _, _, _ := jsonparser.Get(response, "game", "id")

		// Step 1. Process each category for a game
		numCategories := 0
		_, err = jsonparser.ArrayEach(response, func(value []byte, dataType jsonparser.ValueType, offset int, err error) {
			numCategories += 1
			categoryID, _, _, _ := jsonparser.Get(value, "id")
			categoryName, _, _, _ := jsonparser.Get(value, "name")
			categoryRules, _, _, _ := jsonparser.Get(value, "rules")
			categoryNumPlayers, _ := jsonparser.GetInt(value, "numPlayers")
			categoryOutputFile.WriteString(fmt.Sprintf("%s,%s,\"%s\",\"%s\",%d\n", gameID, categoryID, categoryName, categoryRules, categoryNumPlayers))
		}, "categories")
		if err != nil {
			return
		}

		// Step 2. Process each level for a game
		numLevels := 0
		_, err = jsonparser.ArrayEach(response, func(value []byte, dataType jsonparser.ValueType, offset int, err error) {
			numLevels += 1
			levelID, _, _, _ := jsonparser.Get(value, "id")
			levelName, _, _, _ := jsonparser.Get(value, "name")
			levelRules, _, _, _ := jsonparser.Get(value, "rules")
			levelNumPlayers, _ := jsonparser.GetInt(value, "numPlayers")
			levelOutputFile.WriteString(fmt.Sprintf("%s,%s,\"%s\",\"%s\",%d\n", gameID, levelID, levelName, levelRules, levelNumPlayers))
		}, "levels")
		if err != nil {
			return
		}

		// Step 3. Process each variable/value for a game
		_, err = jsonparser.ArrayEach(response, func(value []byte, dataType jsonparser.ValueType, offset int, err error) {
			variableID, _, _, _ := jsonparser.Get(value, "id")
			variableName, _, _, _ := jsonparser.Get(value, "name")
			variableCategory, _, _, _ := jsonparser.Get(value, "categordId")
			variableScope, _, _, _ := jsonparser.Get(value, "categoryScope")
			variableOutputFile.WriteString(fmt.Sprintf("%s,%s,\"%s\",%s,%s\n", gameID, variableID, variableName, variableCategory, variableScope))
		}, "variables")
		if err != nil {
			return
		}

		_, err = jsonparser.ArrayEach(response, func(value []byte, dataType jsonparser.ValueType, offset int, err error) {
			valueID, _, _, _ := jsonparser.Get(value, "id")
			variableID, _, _, _ := jsonparser.Get(value, "variableId")
			valueLabel, _, _, _ := jsonparser.Get(value, "name")
			valueRules, _, _, _ := jsonparser.Get(value, "rules")
			valueOutputFile.WriteString(fmt.Sprintf("%s,%s,%s,\"%s\",\"%s\"\n", gameID, variableID, valueID, valueLabel, valueRules))
		}, "values")
		if err != nil {
			return
		}

		// Step N. Process each game
		gameName, _, _, _ := jsonparser.Get(response, "game", "name")
		gameURL, _, _, _ := jsonparser.Get(response, "game", "url")
		gameType, _, _, _ := jsonparser.Get(response, "game", "type")
		gameEmulator, _ := jsonparser.GetInt(response, "game", "emulator")
		gameReleaseDate, _ := jsonparser.GetInt(response, "game", "releaseDate")
		gameAddedDate, _ := jsonparser.GetInt(response, "game", "addedDate")
		gameRunCount, _ := jsonparser.GetInt(response, "game", "runCount")
		gamePlayerCount, _ := jsonparser.GetInt(response, "game", "totalPlayerCount")
		gameRules, _, _, _ := jsonparser.Get(response, "game", "rules")
		gameOuptutFile.WriteString(fmt.Sprintf("%s,\"%s\",%s,%s,\"%s\",%d,%d,%d,%d,%d,%d,%d\n", gameID, gameName, gameURL, gameType, gameRules, gameReleaseDate, gameAddedDate, gameRunCount, gamePlayerCount, numCategories, numLevels, gameEmulator))
	}

	if err := scanner.Err(); err != nil {
		fmt.Println(err)
		return
	}
}
