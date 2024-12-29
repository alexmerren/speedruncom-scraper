package processor

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/alexmerren/speedruncom-scraper/internal/processor/combinations"
	"github.com/alexmerren/speedruncom-scraper/internal/repository"
	"github.com/alexmerren/speedruncom-scraper/internal/srcom_api"
	"github.com/buger/jsonparser"
)

type WorldRecordDataProcessor struct {
	GameId              string
	WorldRecordDataFile *repository.WriteRepository
	ClientV1            *srcom_api.SrcomV1Client
	ClientV2            *srcom_api.SrcomV2Client
}

func (p *WorldRecordDataProcessor) Process() error {
	response, err := p.ClientV1.GetGame(p.GameId)
	if err != nil {
		return err
	}

	leaderboardCombinations, err := combinations.GenerateLeaderboardCombinations(response)
	if err != nil {
		return err
	}

	for _, combination := range leaderboardCombinations {
		// There is an aggresive rate-limit on the V2 API
		time.Sleep(1_000)

		response, err := p.ClientV2.GetGameRecordHistory(
			combination.GameId,
			combination.CategoryId,
			combination.LevelId,
			combination.VariableIds,
			combination.ValueIds,
		)
		if err != nil {
			return err
		}

		err = p.processWorldRecordHistory(response, combination)
		if err != nil {
			return err
		}
	}

	return nil
}

func (p *WorldRecordDataProcessor) processWorldRecordHistory(data []byte, combination *combinations.Combination) error {
	variablesAndValues := make([]string, len(combination.VariableIds))
	for index := range len(combination.VariableIds) {
		variablesAndValues[index] = fmt.Sprintf("%s=%s", combination.VariableIds[index], combination.ValueIds[index])
	}
	variablesAndValuesString := strings.Join(variablesAndValues, ",")

	_, err := jsonparser.ArrayEach(data, func(runData []byte, dataType jsonparser.ValueType, offset int, err error) {
		runId, _ := jsonparser.GetString(runData, "id")
		gameId, _ := jsonparser.GetString(runData, "gameId")
		categoryId, _ := jsonparser.GetString(runData, "categoryId")
		runTime, _ := jsonparser.GetFloat(runData, "time")
		platformId, _ := jsonparser.GetString(runData, "platformId")
		isEmulated, _ := jsonparser.GetBoolean(runData, "emulator")
		comment, _ := jsonparser.GetString(runData, "comment")
		verifier, _ := jsonparser.GetString(runData, "verifiedById")
		runDate, _ := jsonparser.GetInt(runData, "date")
		submittedDate, _ := jsonparser.GetInt(runData, "dateSubmitted")
		verifiedDate, _ := jsonparser.GetInt(runData, "dateVerified")

		var levelId string
		if combination.LevelId != nil {
			levelId = *combination.LevelId
		}

		playerIdArray := []string{}
		jsonparser.ArrayEach(runData, func(playerId []byte, dataType jsonparser.ValueType, offset int, err error) {
			playerIdArray = append(playerIdArray, string(playerId))
		}, "playerIds")
		runPlayers := strings.Join(playerIdArray, ",")

		p.WorldRecordDataFile.Write([]string{
			runId,
			gameId,
			categoryId,
			levelId,
			variablesAndValuesString,
			strconv.FormatFloat(runTime, 'f', -1, 64),
			platformId,
			strconv.FormatBool(isEmulated),
			runPlayers,
			verifier,
			formatUnixTimestamp(runDate),
			formatUnixTimestamp(submittedDate),
			formatUnixTimestamp(verifiedDate),
			strconv.Quote(comment),
		})
	}, "runList")

	return err
}

func formatUnixTimestamp(timestamp int64) string {
	return time.Unix(timestamp, 0).Format(time.RFC3339)
}
