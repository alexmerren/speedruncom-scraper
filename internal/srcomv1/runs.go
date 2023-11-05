package srcomv1

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/alexmerren/speedruncom-scraper/pkg/srcomv1"
	"github.com/buger/jsonparser"
)

const (
	maxRunsSize = 200

	runsDataOutputFileHeader     = "#runID,gameID,categoryID,levelID,userIDs,status,submittedDate,comment,primaryTime,realTime,platform,emulated\n"
	userRunsDataOutputFileHeader = "#userID,numRuns\n"
)

func ProcessRunsData(usersListInputFile, runsDataOutputFile, userRunsDataOutputFile *os.File) error {
	runsDataOutputFile.WriteString(runsDataOutputFileHeader)
	scanner := bufio.NewScanner(usersListInputFile)
	scanner.Scan()
	for scanner.Scan() {
		userID := scanner.Text()
		userNumRuns := 0
		currentPage := 0

		for {
			runsResponse, err := srcomv1.GetUserRuns(userID, currentPage)
			if err != nil {
				return err
			}

			numRuns, err := processRuns(runsResponse, runsDataOutputFile)
			if err != nil {
				return err
			}

			// Exit condition.
			size, _ := jsonparser.GetInt(runsResponse, "pagination", "size")
			if size < maxRunsSize {
				break
			}

			currentPage += 1
			userNumRuns += numRuns
		}

		userRunsDataOutputFile.WriteString(fmt.Sprintf("%s,%d", userID, userNumRuns))
	}

	if err := scanner.Err(); err != nil {
		return err
	}

	return nil
}

func processRuns(response []byte, outputFile *os.File) (int, error) {
	numRuns := 0
	_, err := jsonparser.ArrayEach(response, func(value []byte, dataType jsonparser.ValueType, offset int, err error) {
		numRuns += 1
		runID, _ := jsonparser.GetString(value, "id")
		gameID, _ := jsonparser.GetString(value, "game")
		categoryID, _ := jsonparser.GetString(value, "category")
		levelID, _ := jsonparser.GetString(value, "level")
		comment, _ := jsonparser.GetString(value, "comment")
		status, _ := jsonparser.GetString(value, "status", "status")
		submittedDate, _ := jsonparser.GetString(value, "submitted")
		timesData, _, _, _ := jsonparser.Get(value, "times")
		primaryTime, _ := jsonparser.GetFloat(timesData, "primary_t")
		realTime, _ := jsonparser.GetFloat(timesData, "realtime_t")
		systemData, _, _, _ := jsonparser.Get(value, "system")
		platform, _ := jsonparser.GetString(systemData, "platform")
		emulated, _ := jsonparser.GetBoolean(systemData, "emulated")

		userIDArray := []string{}
		jsonparser.ArrayEach(value, func(value []byte, dataType jsonparser.ValueType, offset int, err error) {
			playerID, _ := jsonparser.GetString(value, "id")
			userIDArray = append(userIDArray, string(playerID))
		}, "players")
		runUsers := strings.Join(userIDArray, ",")

		outputFile.WriteString(fmt.Sprintf("%s,%s,%s,%s,%s,%s,%s,%q,%f,%f,%s,%t\n", runID, gameID, categoryID, levelID, runUsers, status, submittedDate, comment, primaryTime, realTime, platform, emulated))
	}, "data")
	if err != nil {
		return 0, err
	}

	return numRuns, nil
}
