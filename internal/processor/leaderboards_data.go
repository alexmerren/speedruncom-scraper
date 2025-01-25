package processor

import (
	"errors"
	"fmt"
	"io"
	"strconv"
	"strings"

	"github.com/alexmerren/speedruncom-scraper/internal/repository"
	"github.com/alexmerren/speedruncom-scraper/internal/srcom_api"
	"github.com/buger/jsonparser"
)

type LeaderboardsDataProcessor struct {
	GamesIdListFile  *repository.ReadRepository
	LeaderboardsFile *repository.WriteRepository
	Client           *srcom_api.SrcomV1Client
}

func (p *LeaderboardsDataProcessor) Process() error {
	p.GamesIdListFile.Read()

	for {
		record, err := p.GamesIdListFile.Read()
		if err != nil && errors.Is(err, io.EOF) {
			break
		}

		gameId := record[0]
		gameResponse, err := p.Client.GetGame(gameId)
		if err != nil {
			continue
		}

		_, err = jsonparser.ArrayEach(gameResponse, func(value []byte, dataType jsonparser.ValueType, offset int, _ error) {
			categoryId, _ := jsonparser.GetString(value, "id")
			categoryType, _ := jsonparser.GetString(value, "type")

			// Per-game means that the category is only associated with a full-game run.
			if string(categoryType) == "per-game" {
				leaderboardResponse, _ := p.Client.GetLeaderboardByGameCategory(gameId, categoryId)
				p.processLeaderboard(leaderboardResponse)
			}

			// Per-level means a category can be applied to all levels of a game, and the full-game run.
			// We must retrieve the leaderboard for the given category for all levels.
			if string(categoryType) == "per-level" {
				_, err = jsonparser.ArrayEach(gameResponse, func(value []byte, dataType jsonparser.ValueType, offset int, _ error) {
					levelId, _ := jsonparser.GetString(value, "id")
					leaderboardResponse, _ := p.Client.GetLeaderboardByGameLevelCategory(gameId, levelId, categoryId)
					p.processLeaderboard(leaderboardResponse)
				}, "data", "levels", "data")
			}
		}, "data", "categories", "data")
	}

	return nil
}

func (p *LeaderboardsDataProcessor) processLeaderboard(leaderboardResponse []byte) error {
	_, err := jsonparser.ArrayEach(leaderboardResponse, func(value []byte, dataType jsonparser.ValueType, offset int, _ error) {
		place, _ := jsonparser.GetInt(value, "place")
		runData, _, _, _ := jsonparser.Get(value, "run")
		runId, _ := jsonparser.GetString(runData, "id")
		gameId, _ := jsonparser.GetString(runData, "game")
		categoryId, _ := jsonparser.GetString(runData, "category")
		levelId, _ := jsonparser.GetString(runData, "level")
		runDate, _ := jsonparser.GetString(runData, "date")
		runPrimaryTime, _ := jsonparser.GetFloat(runData, "times", "primary_t")
		runRealTime, _ := jsonparser.GetFloat(runData, "times", "realtime_t")
		runRealTimeNoLoads, _ := jsonparser.GetFloat(runData, "times", "realtime_noloads_t")
		runIngameTime, _ := jsonparser.GetFloat(runData, "times", "ingame_t")
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

		p.LeaderboardsFile.Write([]string{
			runId,
			gameId,
			categoryId,
			levelId,
			strconv.Itoa(int(place)),
			runDate,
			strconv.FormatFloat(runPrimaryTime, 'f', -1, 64),
			strconv.FormatFloat(runRealTime, 'f', -1, 64),
			strconv.FormatFloat(runRealTimeNoLoads, 'f', -1, 64),
			strconv.FormatFloat(runIngameTime, 'f', -1, 64),
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
