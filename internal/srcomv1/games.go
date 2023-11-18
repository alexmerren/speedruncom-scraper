package srcomv1

import (
	"bufio"
	"encoding/csv"
	"os"
	"strconv"

	"github.com/alexmerren/speedruncom-scraper/internal/filesystem"
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
	gameListCsvWriter := csv.NewWriter(gameListOutputFile)
	defer gameListCsvWriter.Flush()

	currentPage := 0

	for {
		response, err := srcomv1.GetGameList(currentPage)
		if err != nil {
			return err
		}

		_, err = jsonparser.ArrayEach(response, func(value []byte, dataType jsonparser.ValueType, offset int, err error) {
			gameID, _ := jsonparser.GetString(value, "id")
			gameListCsvWriter.Write([]string{gameID})
		}, "data")
		if err != nil {
			return err
		}

		// Exit condition.
		size, _ := jsonparser.GetInt(response, "pagination", "size")
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
	gamesCsvWriter := csv.NewWriter(gamesOutputFile)
	defer gamesCsvWriter.Flush()

	categoriesOutputFile.WriteString(categoriesOutputFileHeader)
	categoriesCsvWriter := csv.NewWriter(categoriesOutputFile)
	defer categoriesCsvWriter.Flush()

	levelsOutputFile.WriteString(levelsOutputFileHeader)
	levelsCsvWriter := csv.NewWriter(levelsOutputFile)
	defer levelsCsvWriter.Flush()

	variablesOutputFile.WriteString(variablesOutputFileHeader)
	variablesCsvWriter := csv.NewWriter(variablesOutputFile)
	defer variablesCsvWriter.Flush()

	valuesOutputFile.WriteString(valuesOutputFileHeader)
	valuesCsvWriter := csv.NewWriter(valuesOutputFile)
	defer valuesCsvWriter.Flush()

	scanner := bufio.NewScanner(gameListInputFile)
	scanner.Scan()
	for scanner.Scan() {
		gameID := scanner.Text()
		response, err := srcomv1.GetGame(gameID)
		if err != nil {
			continue
		}

		numCategories, err := processCategory(categoriesCsvWriter, response, gameID)
		if err != nil {
			return err
		}

		numLevels, err := processLevel(levelsCsvWriter, response, gameID)
		if err != nil {
			return err
		}

		err = processVariableValue(variablesCsvWriter, valuesCsvWriter, response, gameID)
		if err != nil {
			return err
		}

		err = processGame(gamesCsvWriter, response, gameID, numCategories, numLevels)
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
	categoryOutputFile *csv.Writer,
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

		categoryOutputFile.Write([]string{gameID, categoryID, categoryName, filesystem.FormatStringForCsv(categoryRules), categoryType, strconv.Itoa(int(categoryNumPlayers))})
	}, "data", "categories", "data")

	return numCategories, err
}

func processLevel(
	levelOutputFile *csv.Writer,
	gameResponse []byte,
	gameID string,
) (int, error) {
	numLevels := 0
	_, err := jsonparser.ArrayEach(gameResponse, func(value []byte, dataType jsonparser.ValueType, offset int, err error) {
		numLevels += 1
		levelID, _ := jsonparser.GetString(value, "id")
		levelName, _ := jsonparser.GetString(value, "name")
		levelRules, _ := jsonparser.GetString(value, "rules")

		levelOutputFile.Write([]string{gameID, levelID, levelName, filesystem.FormatStringForCsv(levelRules)})
	}, "data", "levels", "data")

	return numLevels, err
}

func processVariableValue(
	variableOutputFile,
	valueOutputFile *csv.Writer,
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

		variableOutputFile.Write([]string{gameID, variableID, variableName, variableCategory, variableScope, strconv.FormatBool(variableIsSubcategory), variableDefault})

		jsonparser.ObjectEach(value, func(key []byte, value []byte, dataType jsonparser.ValueType, offset int) error {
			valueID := string(key)
			valueLabel, _ := jsonparser.GetString(value, "label")
			valueRules, _ := jsonparser.GetString(value, "rules")

			valueOutputFile.Write([]string{gameID, variableID, valueID, valueLabel, filesystem.FormatStringForCsv(valueRules)})
			return nil
		}, "values", "values")
	}, "data", "variables", "data")

	return err
}

func processGame(
	gameOutputFile *csv.Writer,
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

	gameOutputFile.Write([]string{gameID, gameName, gameURL, gameReleaseDate, gameCreatedDate, strconv.Itoa(numCategories), strconv.Itoa(numLevels)})

	return nil
}
