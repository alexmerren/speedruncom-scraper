package srcomv1

import (
	"bufio"
	"encoding/csv"
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/alexmerren/speedruncom-scraper/internal/filesystem"
	"github.com/alexmerren/speedruncom-scraper/pkg/srcomv1"
	"github.com/buger/jsonparser"
)

const (
	usersRunsOutputFileHeader = "ID,gameID,categoryID,levelID,date,primaryTime,platform,emulated,players,examiner,values,status,statusReason,verifiedDate\n"

	maxRunsPerPage = 200
)

func ProcessRunsData(
	usersListInputFile,
	usersRunsOutputFile *os.File,
) error {
	usersRunsOutputFile.WriteString(usersRunsOutputFileHeader)
	usersRunsCsvWriter := csv.NewWriter(usersRunsOutputFile)
	defer usersRunsCsvWriter.Flush()

	scanner := bufio.NewScanner(usersListInputFile)
	scanner.Scan()
	for scanner.Scan() {
		userID := scanner.Text()
		_, err := processUserRuns(usersRunsCsvWriter, userID)
		if err != nil {
			return err
		}
	}

	if err := scanner.Err(); err != nil {
		return err
	}

	return nil
}

func processUserRuns(outputFile *csv.Writer, userID string) (int, error) {
	numRuns := 0
	currentPage := 0

	for {
		response, err := srcomv1.GetUserRuns(userID, currentPage)
		if err != nil {
			break
		}

		_, err = jsonparser.ArrayEach(response, func(value []byte, dataType jsonparser.ValueType, offset int, err error) {
			numRuns += 1
			runID, _ := jsonparser.GetString(value, "id")
			gameID, _ := jsonparser.GetString(value, "game")
			categoryID, _ := jsonparser.GetString(value, "category")
			levelID, _ := jsonparser.GetString(value, "level")
			date, _ := jsonparser.GetString(value, "date")
			primaryTime, _ := jsonparser.GetFloat(value, "times", "primary_t")
			platform, _ := jsonparser.GetString(value, "system", "platform")
			emulated, _ := jsonparser.GetBoolean(value, "system", "emulated")

			statusData, _, _, _ := jsonparser.Get(value, "status")
			status, _ := jsonparser.GetString(statusData, "status")
			examiner, _ := jsonparser.GetString(statusData, "examiner")
			statusReason, _ := jsonparser.GetString(statusData, "reason")
			verifiedDate, _ := jsonparser.GetString(statusData, "verify-date")

			playerIDArray := []string{}
			jsonparser.ArrayEach(value, func(value []byte, dataType jsonparser.ValueType, offset int, err error) {
				playerID, _ := jsonparser.GetString(value, "id")
				playerIDArray = append(playerIDArray, string(playerID))
			}, "players", "data")
			players := strings.Join(playerIDArray, ",")

			runValuesArray := []string{}
			jsonparser.ObjectEach(value, func(key []byte, value []byte, dataType jsonparser.ValueType, offset int) error {
				runValuesArray = append(runValuesArray, fmt.Sprintf("%s=%s", string(key), string(value)))
				return nil
			}, "values")
			values := strings.Join(runValuesArray, ",")

			outputFile.Write([]string{runID, gameID, categoryID, levelID, date, strconv.FormatFloat(primaryTime, 'f', -1, 64), platform, strconv.FormatBool(emulated), players, examiner, values, status, filesystem.FormatStringForCsv(statusReason), verifiedDate})
		}, "data")
		if err != nil {
			return 0, err
		}

		size, _ := jsonparser.GetInt(response, "pagination", "size")
		if size < maxRunsPerPage {
			return numRuns, nil
		}

		currentPage += 1
	}

	return numRuns, nil
}
