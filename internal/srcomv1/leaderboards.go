package srcomv1

import (
	"bufio"
	"encoding/csv"
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/alexmerren/speedruncom-scraper/pkg/srcomv1"
	"github.com/buger/jsonparser"
)

const leaderboardsOutputFileHeader = "#runID,gameID,categoryID,levelID,date,primaryTime,place,platform,emulated,players,examiner,verifiedDate,variablesAndValues\n"

func ProcessLeaderboardsData(
	gameListInputFile,
	leaderboardsOutputFile *os.File,
) error {
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

func processLeaderboard(leaderboardsOutputFile *csv.Writer, leaderboardResponse []byte) error {
	_, err := jsonparser.ArrayEach(leaderboardResponse, func(value []byte, dataType jsonparser.ValueType, offset int, _ error) {
		runPlace, _ := jsonparser.GetInt(value, "place")
		runData, _, _, _ := jsonparser.Get(value, "run")
		runID, _ := jsonparser.GetString(runData, "id")
		runGame, _ := jsonparser.GetString(runData, "game")
		runCategory, _ := jsonparser.GetString(runData, "category")
		runLevel, _ := jsonparser.GetString(runData, "level")
		runDate, _ := jsonparser.GetString(runData, "date")
		runPrimaryTime, _ := jsonparser.GetFloat(runData, "times", "primary_t")
		runPlatform, _ := jsonparser.GetString(runData, "system", "platform")
		runEmulated, _ := jsonparser.GetBoolean(runData, "system", "emulated")
		runVerifiedDate, _ := jsonparser.GetString(runData, "status", "verify-date")
		runExaminer, _ := jsonparser.GetString(runData, "status", "examiner")

		playerIDArray := []string{}
		jsonparser.ArrayEach(runData, func(value []byte, dataType jsonparser.ValueType, offset int, err error) {
			playerID, _ := jsonparser.GetString(value, "id")
			playerIDArray = append(playerIDArray, string(playerID))
		}, "players")
		runPlayers := strings.Join(playerIDArray, ",")

		runValuesArray := []string{}
		jsonparser.ObjectEach(runData, func(key []byte, value []byte, dataType jsonparser.ValueType, offset int) error {
			runValuesArray = append(runValuesArray, fmt.Sprintf("%s=%s", string(key), string(value)))
			return nil
		}, "values")
		runValues := strings.Join(runValuesArray, ",")

		leaderboardsOutputFile.Write([]string{runID, runGame, runCategory, runLevel, runDate, strconv.FormatFloat(runPrimaryTime, 'g', 2, 64), strconv.Itoa(int(runPlace)), runPlatform, strconv.FormatBool(runEmulated), runPlayers, runExaminer, runVerifiedDate, runValues})

	}, "data", "runs")

	return err
}
