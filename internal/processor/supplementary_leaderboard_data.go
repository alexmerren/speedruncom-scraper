package processor

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/alexmerren/speedruncom-scraper/internal/processor/combinations"
	"github.com/alexmerren/speedruncom-scraper/internal/repository"
	"github.com/alexmerren/speedruncom-scraper/internal/srcom_api"
	"github.com/buger/jsonparser"
)

type SupplementaryLeaderboardDataProcessor struct {
	GameId                       string
	SupplementaryLeaderboardFile *repository.WriteRepository
	Client                       *srcom_api.SrcomV1Client
}

func (p *SupplementaryLeaderboardDataProcessor) Process() error {
	response, err := p.Client.GetGame(p.GameId)
	if err != nil {
		return err
	}

	leaderboardCombinations, err := combinations.GenerateLeaderboardCombinations(response)
	if err != nil {
		return err
	}

	for _, combination := range leaderboardCombinations {

		leaderboardResponse, err := p.Client.GetLeaderboardByVariables(
			combination.GameId,
			combination.CategoryId,
			combination.LevelId,
			combination.VariableIds,
			combination.ValueIds,
		)
		if err != nil {
			return err
		}

		err = p.processLeaderboard(leaderboardResponse)
		if err != nil {
			return err
		}
	}

	return nil
}

func (p *SupplementaryLeaderboardDataProcessor) processLeaderboard(leaderboardResponse []byte) error {
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

		p.SupplementaryLeaderboardFile.Write([]string{
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
