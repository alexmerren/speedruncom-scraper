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

type RunsDataProcessor struct {
	UsersIdListFile *repository.ReadRepository
	RunsFile        *repository.WriteRepository
	Client          *srcom_api.SrcomV1Client
}

func (p *RunsDataProcessor) Process() error {
	p.UsersIdListFile.Read()

	for {
		record, err := p.UsersIdListFile.Read()
		if err != nil && errors.Is(err, io.EOF) {
			break
		}
		userID := record[0]

		err = p.processRuns(userID)
		if err != nil {
			return err
		}
	}

	return nil
}

func (p *RunsDataProcessor) processRuns(userID string) error {
	currentPage := 0

	for {
		response, err := p.Client.GetRunsByUser(userID, currentPage)
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
			realTime, _ := jsonparser.GetFloat(value, "times", "realtime_t")
			realTimeNoLoads, _ := jsonparser.GetFloat(value, "times", "realtime_noloads_t")
			ingameTime, _ := jsonparser.GetFloat(value, "times", "ingame_t")
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

			p.RunsFile.Write([]string{
				runID,
				gameID,
				categoryID,
				levelID,
				date,
				strconv.FormatFloat(primaryTime, 'f', -1, 64),
				strconv.FormatFloat(realTime, 'f', -1, 64),
				strconv.FormatFloat(realTimeNoLoads, 'f', -1, 64),
				strconv.FormatFloat(ingameTime, 'f', -1, 64),
				platform,
				strconv.FormatBool(emulated),
				players,
				examiner,
				values,
				verifiedDate,
				status,
				strconv.Quote(statusReason),
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
