package main

import (
	"encoding/csv"
	"errors"
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"

	"github.com/alexmerren/speedruncom-scraper/internal"
	"github.com/buger/jsonparser"
)

func main() {
	if err := getLeaderboardDataV1(); err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		os.Exit(1)
	}
}

func getLeaderboardDataV1() error {
	client := internal.NewSrcomV1Client()

	gamesIdListFile, closeFunc, _ := internal.NewCsvReader(internal.GamesIdListFilenameV1)
	defer closeFunc()

	leaderboardsDataFile, closeFunc, _ := internal.NewCsvWriter(internal.LeaderboardsDataFilenameV1)
	leaderboardsDataFile.Write(internal.FileHeaders[internal.LeaderboardsDataFilenameV1])
	defer closeFunc()

	gamesIdListFile.Read()

	for {
		record, err := gamesIdListFile.Read()
		if err != nil && errors.Is(err, io.EOF) {
			break
		}
		gameId := record[0]
		gameResponse, err := client.GetGame(gameId)
		if err != nil {
			continue
		}

		_, err = jsonparser.ArrayEach(gameResponse, func(value []byte, dataType jsonparser.ValueType, offset int, _ error) {
			categoryId, _ := jsonparser.GetString(value, "id")
			categoryType, _ := jsonparser.GetString(value, "type")

			// Per-game means that the category is only associated with a full-game run.
			if string(categoryType) == "per-game" {
				leaderboardResponse, _ := client.GetGameCategoryLeaderboard(gameId, categoryId)
				processLeaderboard(leaderboardsDataFile, leaderboardResponse)
			}

			// Per-level means a category can be applied to all levels of a game, and the full-game run.
			// We must retrieve the leaderboard for the given category for all levels.
			if string(categoryType) == "per-level" {
				_, err = jsonparser.ArrayEach(gameResponse, func(value []byte, dataType jsonparser.ValueType, offset int, _ error) {
					levelId, _ := jsonparser.GetString(value, "id")
					leaderboardResponse, _ := client.GetGameLevelCategoryLeaderboard(gameId, levelId, categoryId)
					processLeaderboard(leaderboardsDataFile, leaderboardResponse)
				}, "data", "levels", "data")
			}
		}, "data", "categories", "data")
	}

	return nil
}

func processLeaderboard(leaderboardsOutputFile *csv.Writer, leaderboardResponse []byte) error {
	_, err := jsonparser.ArrayEach(leaderboardResponse, func(value []byte, dataType jsonparser.ValueType, offset int, _ error) {
		place, _ := jsonparser.GetInt(value, "place")
		runData, _, _, _ := jsonparser.Get(value, "run")
		runId, _ := jsonparser.GetString(runData, "id")
		gameId, _ := jsonparser.GetString(runData, "game")
		categoryId, _ := jsonparser.GetString(runData, "category")
		levelId, _ := jsonparser.GetString(runData, "level")
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

		leaderboardsOutputFile.Write([]string{
			runId,
			gameId,
			categoryId,
			levelId,
			strconv.Itoa(int(place)),
			runDate,
			strconv.FormatFloat(runPrimaryTime, 'f', -1, 64),
			runPlatform,
			strconv.FormatBool(runEmulated),
			runPlayers,
			runExaminer,
			runVerifiedDate,
			runValues,
		})
	}, "data", "runs")

	return err
}
