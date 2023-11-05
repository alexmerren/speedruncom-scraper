package srcomv1

import (
	"bufio"
	"os"

	"github.com/alexmerren/speedruncom-scraper/pkg/srcomv1"
	"github.com/buger/jsonparser"
)

func ProcessLeaderboardsAndGamesData(
	gameListInputFile,
	gamesOutputFile,
	categoriesOutputFile,
	levelsOutputFile,
	variablesOutputFile,
	valuesOutputFile,
	leaderboardsOutputFile *os.File,
) error {
	gamesOutputFile.WriteString(gamesOutputFileHeader)
	categoriesOutputFile.WriteString(categoriesOutputFileHeader)
	levelsOutputFile.WriteString(levelsOutputFileHeader)
	variablesOutputFile.WriteString(variablesOutputFileHeader)
	valuesOutputFile.WriteString(valuesOutputFileHeader)
	leaderboardsOutputFile.WriteString(leaderboardsOutputFileHeader)

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

		_, err = jsonparser.ArrayEach(response, func(value []byte, dataType jsonparser.ValueType, offset int, _ error) {
			categoryID, _ := jsonparser.GetString(value, "id")
			categoryType, _ := jsonparser.GetString(value, "type")

			if string(categoryType) == "per-game" {
				leaderboardResponse, _ := srcomv1.GetGameCategoryLeaderboard(gameID, categoryID)
				processLeaderboard(leaderboardsOutputFile, leaderboardResponse)
			}

			// The levels are embedded so we can immediately iterate over each
			// of the levels to retrieve their respective leaderboard.
			if string(categoryType) == "per-level" {
				_, err = jsonparser.ArrayEach(response, func(value []byte, dataType jsonparser.ValueType, offset int, _ error) {
					levelID, _ := jsonparser.GetString(value, "id")
					leaderboardResponse, _ := srcomv1.GetGameCategoryLevelLeaderboard(gameID, categoryID, levelID)
					processLeaderboard(leaderboardsOutputFile, leaderboardResponse)
				}, "data", "levels", "data")
			}

		}, "data", "categories", "data")
	}

	if err := scanner.Err(); err != nil {
		return err
	}

	return nil
}
