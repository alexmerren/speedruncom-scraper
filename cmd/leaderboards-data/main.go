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
	leaderboardOutputFilenameV1 = "./data/v1/leaderboards-data.csv"
)

func main() {
	getLeaderboardDataV1()
}

//nolint:errcheck // Don't need to check for errors.
func getLeaderboardDataV1() {
	inputFile, err := filesystem.OpenInputFile(allGameIDListV1)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer inputFile.Close()

	leaderboardOutputFile, err := filesystem.CreateOutputFile(leaderboardOutputFilenameV1)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer leaderboardOutputFile.Close()
	leaderboardOutputFile.WriteString("#runID,gameID,categoryID,levelID,date,primaryTime,place,platform,emulated,players,examiner,verifiedDate,variablesAndValues\n")

	scanner := bufio.NewScanner(inputFile)
	scanner.Scan()
	for scanner.Scan() {
		gameID := scanner.Text()
		response, err := srcomv1.GetGame(gameID)
		if err != nil {
			fmt.Println(err)
			continue
		}

		// Iterate through all the categories of a game. Retrieve 'per-game'
		// categories' leaderboards normally, then 'per-level' categories
		// leaderboards can be retrieved via each level.
		_, err = jsonparser.ArrayEach(response, func(value []byte, dataType jsonparser.ValueType, offset int, _ error) {
			categoryID, _ := jsonparser.GetString(value, "id")
			categoryType, _ := jsonparser.GetString(value, "type")

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
				_, err = jsonparser.ArrayEach(response, func(value []byte, dataType jsonparser.ValueType, offset int, _ error) {
					levelID, _ := jsonparser.GetString(value, "id")
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
				if err != nil {
					fmt.Println(err)
					return
				}
			}
		}, "data", "categories", "data")
		if err != nil {
			fmt.Println(err)
			return
		}
	}
}

//nolint:errcheck // Don't need to check for errors.
func processLeaderboard(responseBody []byte, outputFile *os.File) error {
	_, err := jsonparser.ArrayEach(responseBody, func(value []byte, dataType jsonparser.ValueType, offset int, err error) {
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

		outputFile.WriteString(fmt.Sprintf("%s,%s,%s,%s,%s,%0.2f,%d,%s,%t,\"%s\",%s,%s,\"%s\"\n", runID, runGame, runCategory, runLevel, runDate, runPrimaryTime, runPlace, runPlatform, runEmulated, runPlayers, runExaminer, runVerifiedDate, runValues))
	}, "data", "runs")
	if err != nil {
		return err
	}
	return nil
}
