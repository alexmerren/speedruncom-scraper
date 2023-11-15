package srcomv1

import (
	"bufio"
	"encoding/csv"
	"os"

	"github.com/alexmerren/speedruncom-scraper/pkg/srcomv1"
	"github.com/buger/jsonparser"
)

func ProcessGamesAndLeaderboardsData(
	gameListInputFile,
	gamesOutputFile,
	categoriesOutputFile,
	levelsOutputFile,
	variablesOutputFile,
	valuesOutputFile,
	leaderboardsOutputFile *os.File,
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

	leaderboardsOutputFile.WriteString(leaderboardsOutputFileHeader)
	leaderboardsCsvWriter := csv.NewWriter(leaderboardsOutputFile)
	defer leaderboardsCsvWriter.Flush()

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

		_, err = jsonparser.ArrayEach(response, func(value []byte, dataType jsonparser.ValueType, offset int, _ error) {
			categoryID, _ := jsonparser.GetString(value, "id")
			categoryType, _ := jsonparser.GetString(value, "type")

			if string(categoryType) == "per-game" {
				leaderboardResponse, _ := srcomv1.GetGameCategoryLeaderboard(gameID, categoryID)
				processLeaderboard(leaderboardsCsvWriter, leaderboardResponse)
			}

			// The levels are embedded so we can immediately iterate over each
			// of the levels to retrieve their respective leaderboard.
			if string(categoryType) == "per-level" {
				_, err = jsonparser.ArrayEach(response, func(value []byte, dataType jsonparser.ValueType, offset int, _ error) {
					levelID, _ := jsonparser.GetString(value, "id")
					leaderboardResponse, _ := srcomv1.GetGameCategoryLevelLeaderboard(gameID, categoryID, levelID)
					processLeaderboard(leaderboardsCsvWriter, leaderboardResponse)
				}, "data", "levels", "data")
			}

		}, "data", "categories", "data")
	}

	if err := scanner.Err(); err != nil {
		return err
	}

	return nil
}
