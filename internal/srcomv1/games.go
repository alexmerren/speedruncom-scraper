package srcomv1

import (
	"bufio"
	"fmt"
	"os"

	"github.com/alexmerren/speedruncom-scraper/pkg/srcomv1"
	"github.com/buger/jsonparser"
)

const (
	maxSizeAPIv1 = 1000

	gameListOutputFileHeader   = "#gameID\n"
	gamesOutputFileHeader      = "#ID,name,URL,releaseDate,createdDate,numCategories,numLevels\n"
	categoriesOutputFileHeader = "#parentGameID,ID,name,rules,type,numPlayers\n"
	levelsOutputFileHeader     = "#parentGameID,ID,name,rules\n"
	variablesOutputFileHeader  = "#parentGameID,ID,name,category,scope,isSubcategory,defaultValue\n"
	valuesOutputFileHeader     = "#parentGameID,variableID,ID,label,rules\n"
)

func ProcessGamesList(gameListOutputFile *os.File) error {
	gameListOutputFile.WriteString(gameListOutputFileHeader)
	currentPage := 0

	for {
		request, err := srcomv1.GetGameList(currentPage)
		if err != nil {
			return err
		}

		_, err = jsonparser.ArrayEach(request, func(value []byte, dataType jsonparser.ValueType, offset int, err error) {
			gameID, _ := jsonparser.GetString(value, "id")
			gameListOutputFile.WriteString(fmt.Sprintf("%s\n", gameID))
		}, "data")
		if err != nil {
			return err
		}

		// Exit condition.
		size, _ := jsonparser.GetInt(request, "pagination", "size")
		if size < maxSizeAPIv1 {
			return nil
		}

		currentPage += 1
	}
}

func ProcessGamesData(
	gameListInputFile,
	gamesOutputFile,
	categoriesOutputFile,
	levelsOutputFile,
	variablesOutputFile,
	valuesOutputFile *os.File,
) error {
	gamesOutputFile.WriteString(gamesOutputFileHeader)
	categoriesOutputFile.WriteString(categoriesOutputFileHeader)
	levelsOutputFile.WriteString(levelsOutputFileHeader)
	variablesOutputFile.WriteString(variablesOutputFileHeader)
	valuesOutputFile.WriteString(valuesOutputFileHeader)

	scanner := bufio.NewScanner(gameListInputFile)
	scanner.Scan()
	for scanner.Scan() {
		gameID := scanner.Text()
		response, err := srcomv1.GetGame(gameID)
		if err != nil {
			continue
		}

		numCategories, err := processCategory(categoriesOutputFile, response, gameID)
		if err != nil {
			return err
		}

		numLevels, err := processLevel(levelsOutputFile, response, gameID)
		if err != nil {
			return err
		}

		err = processVariableValue(variablesOutputFile, valuesOutputFile, response, gameID)
		if err != nil {
			return err
		}

		err = processGame(gamesOutputFile, response, gameID, numCategories, numLevels)
		if err != nil {
			return err
		}
	}

	if err := scanner.Err(); err != nil {
		return err
	}

	return nil
}

func processCategory(
	categoryOutputFile *os.File,
	gameResponse []byte,
	gameID string,
) (int, error) {
	numCategories := 0
	_, err := jsonparser.ArrayEach(gameResponse, func(value []byte, dataType jsonparser.ValueType, offset int, err error) {
		numCategories += 1
		categoryID, _ := jsonparser.GetString(value, "id")
		categoryName, _ := jsonparser.GetString(value, "name")
		categoryRules, _ := jsonparser.GetString(value, "rules")
		categoryNumPlayers, _ := jsonparser.GetInt(value, "players", "value")
		categoryType, _ := jsonparser.GetString(value, "type")

		categoryOutputFile.WriteString(fmt.Sprintf("%s,%s,%q,%q,%s,%d\n", gameID, categoryID, categoryName, categoryRules, categoryType, categoryNumPlayers))
	}, "data", "categories", "data")

	return numCategories, err
}

func processLevel(
	levelOutputFile *os.File,
	gameResponse []byte,
	gameID string,
) (int, error) {
	numLevels := 0
	_, err := jsonparser.ArrayEach(gameResponse, func(value []byte, dataType jsonparser.ValueType, offset int, err error) {
		numLevels += 1
		levelID, _ := jsonparser.GetString(value, "id")
		levelName, _ := jsonparser.GetString(value, "name")
		levelRules, _ := jsonparser.GetString(value, "rules")

		levelOutputFile.WriteString(fmt.Sprintf("%s,%s,%q,%q\n", gameID, levelID, levelName, levelRules))
	}, "data", "levels", "data")

	return numLevels, err
}

func processVariableValue(
	variableOutputFile,
	valueOutputFile *os.File,
	gameResponse []byte,
	gameID string,
) error {
	_, err := jsonparser.ArrayEach(gameResponse, func(value []byte, dataType jsonparser.ValueType, offset int, _ error) {
		variableID, _ := jsonparser.GetString(value, "id")
		variableName, _ := jsonparser.GetString(value, "name")
		variableCategory, _ := jsonparser.GetString(value, "category")
		variableScope, _ := jsonparser.GetString(value, "scope", "type")
		variableIsSubcategory, _ := jsonparser.GetBoolean(value, "is-subcategory")
		variableDefault, _ := jsonparser.GetString(value, "values", "default")

		variableOutputFile.WriteString(fmt.Sprintf("%s,%s,%q,%s,%s,%t,%s\n", gameID, variableID, variableName, variableCategory, variableScope, variableIsSubcategory, variableDefault))

		jsonparser.ObjectEach(value, func(key []byte, value []byte, dataType jsonparser.ValueType, offset int) error {
			valueID := string(key)
			valueLabel, _ := jsonparser.GetString(value, "label")
			valueRules, _ := jsonparser.GetString(value, "rules")

			valueOutputFile.WriteString(fmt.Sprintf("%s,%s,%s,%q,%q\n", gameID, variableID, valueID, valueLabel, valueRules))
			return nil
		}, "values", "values")
	}, "data", "variables", "data")

	return err
}

func processGame(
	gameOutputFile *os.File,
	gameResponse []byte,
	gameID string,
	numCategories,
	numLevels int,
) error {
	gameData, _, _, err := jsonparser.Get(gameResponse, "data")
	if err != nil {
		return err
	}

	gameName, _ := jsonparser.GetString(gameData, "names", "international")
	gameURL, _ := jsonparser.GetString(gameData, "abbreviation")
	gameReleaseDate, _ := jsonparser.GetString(gameData, "release-date")
	gameCreatedDate, _ := jsonparser.GetString(gameData, "created")

	gameOutputFile.WriteString(fmt.Sprintf("%s,%q,%s,%s,%s,%d,%d\n", gameID, gameName, gameURL, gameReleaseDate, gameCreatedDate, numCategories, numLevels))

	return nil
}
