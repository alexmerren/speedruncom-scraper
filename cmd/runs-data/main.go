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
	"github.com/alexmerren/speedruncom-scraper/pkg/srcom_api"
	"github.com/buger/jsonparser"
)

func main() {
	if err := generateRunsDataV1(); err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		os.Exit(1)
	}
}

func generateRunsDataV1() error {
	client := internal.NewSrcomV1Client()

	usersIdListFile, closeFunc, _ := internal.NewCsvReader(internal.UsersIdListFilenameV1)
	defer closeFunc()

	runsDataFile, closeFunc, _ := internal.NewCsvWriter(internal.RunsDataFilenameV1)
	runsDataFile.Write(internal.FileHeaders[internal.RunsDataFilenameV1])
	defer closeFunc()

	usersIdListFile.Read()

	for {
		record, err := usersIdListFile.Read()
		if err != nil && errors.Is(err, io.EOF) {
			break
		}

		userID := record[0]
		err = processRuns(runsDataFile, client, userID)
		if err != nil {
			return err
		}
	}

	return nil
}

func processRuns(runsDataFile *csv.Writer, client *srcom_api.SrcomV1Client, userID string) error {
	currentPage := 0

	for {
		response, err := client.GetRunsByUser(userID, currentPage)
		if err != nil {
			break
		}

		_, err = jsonparser.ArrayEach(response, func(value []byte, dataType jsonparser.ValueType, offset int, err error) {
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

			runsDataFile.Write([]string{
				runID,
				gameID,
				categoryID,
				levelID,
				date,
				strconv.FormatFloat(primaryTime, 'f', -1, 64),
				platform,
				strconv.FormatBool(emulated),
				players,
				examiner,
				values,
				status,
				internal.FormatCsvString(statusReason),
				verifiedDate,
			})
		}, "data")
		if err != nil {
			return err
		}

		size, _ := jsonparser.GetInt(response, "pagination", "size")
		if size < 200 {
			return nil
		}

		currentPage += 1
	}

	return nil
}
