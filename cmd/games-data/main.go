package main

import (
	"bufio"
	"fmt"
	"strings"
	"sync"

	"github.com/alexmerren/speedruncom-scraper/internal/filesystem"
	"github.com/alexmerren/speedruncom-scraper/internal/srcomv1"
	"github.com/alexmerren/speedruncom-scraper/internal/srcomv2"
	"github.com/buger/jsonparser"
)

const (
	allGameIDListV1          = "./data/v1/games-id-list.csv"
	gameOutputFilenameV1     = "./data/v1/games-data.csv"
	categoryOutputFilenameV1 = "./data/v1/categories-data.csv"
	levelOutputFilenameV1    = "./data/v1/levels-data.csv"
	variableOutputFileV1     = "./data/v1/variables-data.csv"
	valueOutputFileV1        = "./data/v1/values-data.csv"

	allGameIDListV2          = "./data/v2/games-id-list.csv"
	gameOutputFilenameV2     = "./data/v2/games-data.csv"
	categoryOutputFilenameV2 = "./data/v2/categories-data.csv"
	levelOutputFilenameV2    = "./data/v2/levels-data.csv"
	variableOutputFileV2     = "./data/v2/variables-data.csv"
	valueOutputFileV2        = "./data/v2/values-data.csv"
)

func main() {
	wg := &sync.WaitGroup{}

	wg.Add(2)
	go func() {
		defer wg.Done()
		getGameDataV1()
	}()

	go func() {
		defer wg.Done()
		getGameDataV2()
	}()

	wg.Wait()
}

//nolint:errcheck // Don't need to check for errors.
func getGameDataV1() {
	inputFile, err := filesystem.OpenInputFile(allGameIDListV1)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer inputFile.Close()

	gameOutputFile, err := filesystem.CreateOutputFile(gameOutputFilenameV1)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer gameOutputFile.Close()
	gameOutputFile.WriteString("#ID,name,URL,releaseDate,createdDate,numCategories,numLevels\n")

	categoryOutputFile, err := filesystem.CreateOutputFile(categoryOutputFilenameV1)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer categoryOutputFile.Close()
	categoryOutputFile.WriteString("#parentGameID,ID,name,rules,type,numPlayers\n")

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
			categoryID, _ := jsonparser.GetString(value, "id")
			categoryName, _ := jsonparser.GetString(value, "name")
			categoryRules, _ := jsonparser.GetString(value, "rules")
			categoryNumPlayers, _ := jsonparser.GetInt(value, "players", "value")
			categoryType, _ := jsonparser.GetString(value, "type")
			categoryOutputFile.WriteString(fmt.Sprintf("%s,%s,\"%s\",\"%s\",%s,%d\n", gameID, categoryID, categoryName, categoryRules, categoryType, categoryNumPlayers))
		}, "data", "categories", "data")
		if err != nil {
			return
		}

		// Step 2. Process each level for a game
		numLevels := 0
		_, err = jsonparser.ArrayEach(response, func(value []byte, dataType jsonparser.ValueType, offset int, err error) {
			numLevels += 1
			levelID, _ := jsonparser.GetString(value, "id")
			levelName, _ := jsonparser.GetString(value, "name")
			levelRules, _ := jsonparser.GetString(value, "rules")
			levelOutputFile.WriteString(fmt.Sprintf("%s,%s,\"%s\",\"%s\"\n", gameID, levelID, levelName, levelRules))
		}, "data", "levels", "data")
		if err != nil {
			return
		}

		// Step 3. Process each variable/value for a game
		_, err = jsonparser.ArrayEach(response, func(value []byte, dataType jsonparser.ValueType, offset int, _ error) {
			variableID, _ := jsonparser.GetString(value, "id")
			variableName, _ := jsonparser.GetString(value, "name")
			variableCategory, _ := jsonparser.GetString(value, "category")
			variableScope, _ := jsonparser.GetString(value, "scope", "type")
			variableIsSubcategory, _ := jsonparser.GetBoolean(value, "is-subcategory")
			variableDefault, _ := jsonparser.GetString(value, "values", "default")
			variableOutputFile.WriteString(fmt.Sprintf("%s,%s,\"%s\",%s,%s,%t,%s\n", gameID, variableID, variableName, variableCategory, variableScope, variableIsSubcategory, variableDefault))

			err = jsonparser.ObjectEach(value, func(key []byte, value []byte, dataType jsonparser.ValueType, offset int) error {
				valueID := string(key)
				valueLabel, _ := jsonparser.GetString(value, "label")
				valueRules, _ := jsonparser.GetString(value, "rules")
				valueOutputFile.WriteString(fmt.Sprintf("%s,%s,%s,\"%s\",\"%s\"\n", gameID, variableID, valueID, valueLabel, valueRules))
				return nil
			}, "values", "values")
			if err != nil {
				fmt.Println(err)
				return
			}
		}, "data", "variables", "data")
		if err != nil {
			fmt.Println(err)
			return
		}

		// Step N. Process each game
		gameName, _ := jsonparser.GetString(response, "data", "names", "international")
		gameURL, _ := jsonparser.GetString(response, "data", "abbreviation")
		gameReleaseDate, _ := jsonparser.GetString(response, "data", "release-date")
		gameCreatedDate, _ := jsonparser.GetString(response, "data", "created")
		gameOutputFile.WriteString(fmt.Sprintf("%s,\"%s\",%s,%s,%s,%d,%d\n", gameID, gameName, gameURL, gameReleaseDate, gameCreatedDate, numCategories, numLevels))
	}

	if err := scanner.Err(); err != nil {
		fmt.Println(err)
		return
	}
}

//nolint:errcheck // Don't need to check for errors.
func getGameDataV2() {
	gameOutputFile, err := filesystem.CreateOutputFile(gameOutputFilenameV2)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer gameOutputFile.Close()
	gameOutputFile.WriteString("#ID,name,URL,type,rules,releaseDate,addedDate,runCount,playerCount\n")

	currentPage := 0
	request, _ := srcomv2.GetGameList(currentPage)
	lastPage, err := jsonparser.GetInt(request, "pagination", "pages")
	if err != nil {
		fmt.Println(err)
		return
	}

	for int64(currentPage) < lastPage {
		_, err := jsonparser.ArrayEach(request, func(value []byte, dataType jsonparser.ValueType, offset int, _ error) {
			gameID, _ := jsonparser.GetString(value, "id")
			gameName, _ := jsonparser.GetString(value, "name")
			gameNameFormatted := strings.ReplaceAll(gameName, "\"", "'")
			gameURL, _ := jsonparser.GetString(value, "url")
			gameType, _ := jsonparser.GetString(value, "type")
			gameReleaseDate, _ := jsonparser.GetInt(value, "releaseDate")
			gameAddedDate, _ := jsonparser.GetInt(value, "addedDate")
			gameRunCount, _ := jsonparser.GetInt(value, "runCount")
			gamePlayerCount, _ := jsonparser.GetInt(value, "totalPlayerCount")
			gameRules, _ := jsonparser.GetString(value, "rules")
			gameRulesFormatted := strings.ReplaceAll(gameRules, "\"", "'")
			gameOutputFile.WriteString(fmt.Sprintf("%s,%q,%s,%s,%q,%d,%d,%d,%d\n", gameID, gameNameFormatted, gameURL, gameType, gameRulesFormatted, gameReleaseDate, gameAddedDate, gameRunCount, gamePlayerCount))
		}, "gameList")
		if err != nil {
			fmt.Println(err)
			return
		}

		currentPage += 1

		request, err = srcomv2.GetGameList(currentPage)
		if err != nil {
			fmt.Println(err)
			return
		}
	}
}
