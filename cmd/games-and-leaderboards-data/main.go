package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/alexmerren/speedruncom-scraper/internal/filesystem"
	"github.com/alexmerren/speedruncom-scraper/internal/srcomv1"
	"github.com/buger/jsonparser"
)

const (
	allGameIDListV1             = "./data/v1/games-id-list.csv"
	gameOutputFilenameV1        = "./data/v1/games-data.csv"
	categoryOutputFilenameV1    = "./data/v1/categories-data.csv"
	levelOutputFilenameV1       = "./data/v1/levels-data.csv"
	variableOutputFileV1        = "./data/v1/variables-data.csv"
	valueOutputFileV1           = "./data/v1/values-data.csv"
	leaderboardOutputFilenameV1 = "./data/v1/leaderboards-data.csv"
)

func main() {
	getGameAndLeaderboardDataV1()
}

func getGameAndLeaderboardDataV1() {
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

	leaderboardOutputFile, err := filesystem.CreateOutputFile(leaderboardOutputFilenameV1)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer leaderboardOutputFile.Close()
	leaderboardOutputFile.WriteString("#runID,gameID,categoryID,levelID,date,primaryTime,platform,emulated,players,examiner,verifiedDate,variablesAndValues\n")

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

		numCategories, err := processCategories(gameID, response, categoryOutputFile)
		if err != nil {
			return
		}

		// Step 2. Process each level for a game
		numLevels, err := processLevels(gameID, response, levelOutputFile)
		if err != nil {
			return
		}

		// Step 3. Process each variable/value for a game
		err = processVariablesAndValues(gameID, response, variableOutputFile, valueOutputFile)
		if err != nil {
			return
		}

		// Step 4. Process each game
		processGame(gameID, numCategories, numLevels, response, gameOutputFile)

		// Step 5. Process the leaderboard for each game.
		_, err = jsonparser.ArrayEach(response, func(value []byte, dataType jsonparser.ValueType, offset int, err error) {
			categoryID, _, _, _ := jsonparser.Get(value, "id")
			categoryType, _, _, _ := jsonparser.Get(value, "type")

			if string(categoryType) == "per-game" {
				leaderboardResponse, err := srcomv1.GetGameCategoryLeaderboard(gameID, string(categoryID))
				if err != nil {
					fmt.Println(err)
					return
				}

				err = processLeaderboard(leaderboardResponse, leaderboardOutputFile)
				if err != nil {
					fmt.Println(err)
					return
				}
			}

			// The levels are embedded so we can immediately iterate over each
			// of the levels to retrieve their respective leaderboard.
			if string(categoryType) == "per-level" {
				_, err = jsonparser.ArrayEach(response, func(value []byte, dataType jsonparser.ValueType, offset int, err error) {
					levelID, _, _, _ := jsonparser.Get(value, "id")
					leaderboardResponse, err := srcomv1.GetGameCategoryLevelLeaderboard(gameID, string(categoryID), string(levelID))
					if err != nil {
						fmt.Println(err)
						return
					}

					err = processLeaderboard(leaderboardResponse, leaderboardOutputFile)
					if err != nil {
						fmt.Println(err)
						return
					}
				}, "data", "levels", "data")
			}
		}, "data", "categories", "data")

	}

	if err := scanner.Err(); err != nil {
		fmt.Println(err)
		return
	}
}

func processCategories(gameID string, responseBody []byte, outputFile *os.File) (int, error) {
	numCategories := 0
	_, err := jsonparser.ArrayEach(responseBody, func(value []byte, dataType jsonparser.ValueType, offset int, err error) {
		numCategories += 1
		categoryID, _, _, _ := jsonparser.Get(value, "id")
		categoryName, _, _, _ := jsonparser.Get(value, "name")
		categoryRules, _, _, _ := jsonparser.Get(value, "rules")
		categoryNumPlayers, _ := jsonparser.GetInt(value, "players", "value")
		categoryType, _, _, _ := jsonparser.Get(value, "type")
		outputFile.WriteString(fmt.Sprintf("%s,%s,\"%s\",\"%s\",%s,%d\n", gameID, categoryID, categoryName, categoryRules, categoryType, categoryNumPlayers))
	}, "data", "categories", "data")
	return numCategories, err
}

func processLevels(gameID string, responseBody []byte, outputFile *os.File) (int, error) {
	numLevels := 0
	_, err := jsonparser.ArrayEach(responseBody, func(value []byte, dataType jsonparser.ValueType, offset int, err error) {
		numLevels += 1
		levelID, _, _, _ := jsonparser.Get(value, "id")
		levelName, _, _, _ := jsonparser.Get(value, "name")
		levelRules, _, _, _ := jsonparser.Get(value, "rules")
		outputFile.WriteString(fmt.Sprintf("%s,%s,\"%s\",\"%s\"\n", gameID, levelID, levelName, levelRules))
	}, "data", "levels", "data")
	return numLevels, err
}

func processVariablesAndValues(gameID string, responseBody []byte, variableOutputFile, valueOutputFile *os.File) error {
	_, err := jsonparser.ArrayEach(responseBody, func(value []byte, dataType jsonparser.ValueType, offset int, err error) {
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
	return err
}

func processGame(gameID string, numCategories, numLevels int, responseBody []byte, outputFile *os.File) {
	gameName, _, _, _ := jsonparser.Get(responseBody, "data", "names", "international")
	gameURL, _, _, _ := jsonparser.Get(responseBody, "data", "abbreviation")
	gameReleaseDate, _, _, _ := jsonparser.Get(responseBody, "data", "release-date")
	gameCreatedDate, _, _, _ := jsonparser.Get(responseBody, "data", "created")
	outputFile.WriteString(fmt.Sprintf("%s,\"%s\",%s,%s,%s,%d,%d\n", gameID, gameName, gameURL, gameReleaseDate, gameCreatedDate, numCategories, numLevels))
}

func processLeaderboard(responseBody []byte, outputFile *os.File) error {
	_, err := jsonparser.ArrayEach(responseBody, func(value []byte, dataType jsonparser.ValueType, offset int, err error) {
		runData, _, _, _ := jsonparser.Get(value, "run")
		runID, _, _, _ := jsonparser.Get(runData, "id")
		runGame, _, _, _ := jsonparser.Get(runData, "game")
		runCategory, _, _, _ := jsonparser.Get(runData, "category")
		runLevel, _, _, _ := jsonparser.Get(runData, "level")
		runDate, _, _, _ := jsonparser.Get(runData, "date")
		runPrimaryTime, _ := jsonparser.GetFloat(runData, "times", "primary_t")
		runPlatform, _, _, _ := jsonparser.Get(runData, "system", "platform")
		runEmulated, _ := jsonparser.GetBoolean(runData, "system", "emulated")
		runVerifiedDate, _, _, _ := jsonparser.Get(runData, "status", "verify-date")
		runExaminer, _, _, _ := jsonparser.Get(runData, "status", "examiner")

		playerIDArray := []string{}
		_, err = jsonparser.ArrayEach(runData, func(value []byte, dataType jsonparser.ValueType, offset int, err error) {
			playerID, _, _, _ := jsonparser.Get(value, "id")
			playerIDArray = append(playerIDArray, string(playerID))
		}, "players")
		runPlayers := strings.Join(playerIDArray, ",")

		runValuesArray := []string{}
		err = jsonparser.ObjectEach(runData, func(key []byte, value []byte, dataType jsonparser.ValueType, offset int) error {
			runValuesArray = append(runValuesArray, fmt.Sprintf("%s=%s", string(key), string(value)))
			return nil
		}, "values")
		runValues := strings.Join(runValuesArray, ",")

		outputFile.WriteString(fmt.Sprintf("%s,%s,%s,%s,%s,%0.2f,%s,%t,\"%s\",%s,%s,\"%s\"\n", runID, runGame, runCategory, runLevel, runDate, runPrimaryTime, runPlatform, runEmulated, runPlayers, runExaminer, runVerifiedDate, runValues))
	}, "data", "runs")
	if err != nil {
		return err
	}
	return nil
}
