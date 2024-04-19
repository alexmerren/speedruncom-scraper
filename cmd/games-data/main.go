package main

import (
	"encoding/csv"
	"errors"
	"fmt"
	"io"
	"os"
	"strconv"
	"sync"

	"github.com/alexmerren/speedruncom-scraper/internal"
	"github.com/buger/jsonparser"
)

func main() {
	var wg sync.WaitGroup

	wg.Add(1)
	go func() {
		defer wg.Done()
		if err := generateGamesDataV1(); err != nil {
			fmt.Fprintf(os.Stderr, "error: %v\n", err)
			os.Exit(1)
		}
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()
		if err := generateGamesDataV2(); err != nil {
			fmt.Fprintf(os.Stderr, "error: %v\n", err)
			os.Exit(1)
		}
	}()

	wg.Wait()
}

func generateGamesDataV1() error {
	client := internal.NewSrcomV1Client()

	gamesIdListFile, closeFunc, _ := internal.NewCsvReader(internal.GamesIdListFilenameV1)
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
		response, err := client.GetGame(gameId)
		if err != nil {
			continue
		}

		numCategories, err := processCategory(categoriesDataFile, response, gameId)
		if err != nil {
			return err
		}

		numLevels, err := processLevel(levelsDataFile, response, gameId)
		if err != nil {
			return err
		}

		err = processVariableValue(variablesDataFile, valuesDataFile, response, gameId)
		if err != nil {
			return err
		}

		err = processGame(gamesDataFile, response, numCategories, numLevels, gameId)
		if err != nil {
			return err
		}
	}

	return nil
}

func processCategory(categoriesDataFile *csv.Writer, gameResponse []byte, gameId string) (int, error) {
	numCategories := 0
	_, err := jsonparser.ArrayEach(gameResponse, func(value []byte, dataType jsonparser.ValueType, offset int, err error) {
		numCategories += 1
		categoryId, _ := jsonparser.GetString(value, "id")
		categoryName, _ := jsonparser.GetString(value, "name")
		categoryRules, _ := jsonparser.GetString(value, "rules")
		categoryNumPlayers, _ := jsonparser.GetInt(value, "players", "value")
		categoryType, _ := jsonparser.GetString(value, "type")

		categoriesDataFile.Write([]string{
			gameId,
			categoryId,
			categoryName,
			internal.FormatCsvString(categoryRules),
			categoryType,
			strconv.Itoa(int(categoryNumPlayers)),
		})
	}, "data", "categories", "data")

	return numCategories, err
}

func processLevel(levelsDataFile *csv.Writer, gameResponse []byte, gameId string) (int, error) {
	numLevels := 0
	_, err := jsonparser.ArrayEach(gameResponse, func(value []byte, dataType jsonparser.ValueType, offset int, err error) {
		numLevels += 1
		levelId, _ := jsonparser.GetString(value, "id")
		levelName, _ := jsonparser.GetString(value, "name")
		levelRules, _ := jsonparser.GetString(value, "rules")

		levelsDataFile.Write([]string{gameId, levelId, levelName, internal.FormatCsvString(levelRules)})
	}, "data", "levels", "data")

	return numLevels, err
}

func processVariableValue(variablesDataFile, valuesDataFile *csv.Writer, gameResponse []byte, gameId string) error {
	_, err := jsonparser.ArrayEach(gameResponse, func(value []byte, dataType jsonparser.ValueType, offset int, _ error) {
		variableId, _ := jsonparser.GetString(value, "id")
		variableName, _ := jsonparser.GetString(value, "name")
		variableCategory, _ := jsonparser.GetString(value, "category")
		variableScope, _ := jsonparser.GetString(value, "scope", "type")
		variableIsSubcategory, _ := jsonparser.GetBoolean(value, "is-subcategory")
		variableDefault, _ := jsonparser.GetString(value, "values", "default")

		variablesDataFile.Write([]string{
			gameId,
			variableId,
			variableName,
			variableCategory,
			variableScope,
			strconv.FormatBool(variableIsSubcategory),
			variableDefault,
		})

		jsonparser.ObjectEach(value, func(key []byte, value []byte, dataType jsonparser.ValueType, offset int) error {
			valueId := string(key)
			valueLabel, _ := jsonparser.GetString(value, "label")
			valueRules, _ := jsonparser.GetString(value, "rules")

			valuesDataFile.Write([]string{
				gameId,
				variableId,
				valueId,
				valueLabel,
				internal.FormatCsvString(valueRules),
			})
			return nil
		}, "values", "values")
	}, "data", "variables", "data")

	return err
}

func processGame(gamesDataFile *csv.Writer, gameResponse []byte, numCategories, numLevels int, gameId string) error {
	gameData, _, _, err := jsonparser.Get(gameResponse, "data")
	if err != nil {
		return err
	}

	gameName, _ := jsonparser.GetString(gameData, "names", "international")
	gameURL, _ := jsonparser.GetString(gameData, "abbreviation")
	gameReleaseDate, _ := jsonparser.GetString(gameData, "release-date")
	gameCreatedDate, _ := jsonparser.GetString(gameData, "created")

	gamesDataFile.Write([]string{
		gameId,
		gameName,
		gameURL,
		gameReleaseDate,
		gameCreatedDate,
		strconv.Itoa(numCategories),
		strconv.Itoa(numLevels),
	})

	return nil
}

func generateGamesDataV2() error {
	client := internal.NewSrcomV2Client()

	gamesDataFile, closeFunc, _ := internal.NewCsvWriter(internal.GamesDataFilenameV2)
	gamesDataFile.Write(internal.FileHeaders[internal.GamesDataFilenameV2])
	defer closeFunc()

	currentPage := 0
	request, _ := client.GetGameList(currentPage)
	lastPage, err := jsonparser.GetInt(request, "pagination", "pages")
	if err != nil {
		return err
	}

	for int64(currentPage) <= lastPage {
		_, err := jsonparser.ArrayEach(request, func(value []byte, dataType jsonparser.ValueType, offset int, err error) {
			id, _ := jsonparser.GetString(value, "id")
			name, _ := jsonparser.GetString(value, "name")
			url, _ := jsonparser.GetString(value, "url")
			gameType, _ := jsonparser.GetString(value, "type")
			releaseDate, _ := jsonparser.GetInt(value, "releaseDate")
			addedDate, _ := jsonparser.GetInt(value, "addedDate")
			runCount, _ := jsonparser.GetInt(value, "runCount")
			playerCount, _ := jsonparser.GetInt(value, "totalPlayerCount")
			rules, _ := jsonparser.GetString(value, "rules")

			gamesDataFile.Write([]string{
				id,
				name,
				url,
				gameType,
				strconv.Itoa(int(releaseDate)),
				strconv.Itoa(int(addedDate)),
				strconv.Itoa(int(runCount)),
				strconv.Itoa(int(playerCount)),
				internal.FormatCsvString(rules),
			})
		}, "gameList")
		if err != nil {
			return err
		}

		currentPage += 1

		request, err = client.GetGameList(currentPage)
		if err != nil {
			return err
		}
	}

	return nil
}
