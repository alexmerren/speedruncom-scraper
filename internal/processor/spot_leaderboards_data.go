package processor

import (
	"fmt"

	"github.com/alexmerren/speedruncom-scraper/internal/repository"
	"github.com/alexmerren/speedruncom-scraper/internal/srcom_api"
)

type SpotLeaderboardsDataProcessor struct {
	LeaderboardsFile *repository.WriteRepository
	Client           *srcom_api.SrcomV1Client
}

func (p *SpotLeaderboardsDataProcessor) Process(gameId string) error {
	response, err := p.Client.GetGame(gameId)
	if err != nil {
		return err
	}

	combinations, err := generateCombinations(response)
	if err != nil {
		return err
	}

	for _, combination := range combinations {
		if isValid := combination.isValid(); !isValid {
			return fmt.Errorf("combination is invalid: %+v", combination)
		}
	}

	for _, combination := range combinations {
		fmt.Println(combination)
	}

	return nil
}

type combination struct {
	gameId      string
	categoryId  string
	levelId     *string
	variableIds []string
	valueIds    []string
}

func (c *combination) isValid() bool {
	return len(c.variableIds) == len(c.valueIds)
}

func generateCombinations(response []byte) ([]*combination, error) {
	return nil, nil
}

// func (p *LeaderboardsDataProcessor) processLeaderboard(leaderboardResponse []byte) error {
// 	_, err := jsonparser.ArrayEach(leaderboardResponse, func(value []byte, dataType jsonparser.ValueType, offset int, _ error) {
// 		place, _ := jsonparser.GetInt(value, "place")
// 		runData, _, _, _ := jsonparser.Get(value, "run")
// 		runId, _ := jsonparser.GetString(runData, "id")
// 		gameId, _ := jsonparser.GetString(runData, "game")
// 		categoryId, _ := jsonparser.GetString(runData, "category")
// 		levelId, _ := jsonparser.GetString(runData, "level")
// 		runDate, _ := jsonparser.GetString(runData, "date")
// 		runPrimaryTime, _ := jsonparser.GetFloat(runData, "times", "primary_t")
// 		runPlatform, _ := jsonparser.GetString(runData, "system", "platform")
// 		runEmulated, _ := jsonparser.GetBoolean(runData, "system", "emulated")
// 		runVerifiedDate, _ := jsonparser.GetString(runData, "status", "verify-date")
// 		runExaminer, _ := jsonparser.GetString(runData, "status", "examiner")

// 		playerIDArray := []string{}
// 		jsonparser.ArrayEach(runData, func(value []byte, dataType jsonparser.ValueType, offset int, err error) {
// 			playerID, _ := jsonparser.GetString(value, "id")
// 			playerIDArray = append(playerIDArray, string(playerID))
// 		}, "players")
// 		runPlayers := strings.Join(playerIDArray, ",")

// 		runValuesArray := []string{}
// 		jsonparser.ObjectEach(runData, func(key []byte, value []byte, dataType jsonparser.ValueType, offset int) error {
// 			runValuesArray = append(runValuesArray, fmt.Sprintf("%s=%s", string(key), string(value)))
// 			return nil
// 		}, "values")
// 		runValues := strings.Join(runValuesArray, ",")

// 		p.LeaderboardsFile.Write([]string{
// 			runId,
// 			gameId,
// 			categoryId,
// 			levelId,
// 			strconv.Itoa(int(place)),
// 			runDate,
// 			strconv.FormatFloat(runPrimaryTime, 'f', -1, 64),
// 			runPlatform,
// 			strconv.FormatBool(runEmulated),
// 			runPlayers,
// 			runExaminer,
// 			runVerifiedDate,
// 			runValues,
// 		})
// 	}, "data", "runs")

// 	return err
// }
