package srcomv1

import (
	"bufio"
	"encoding/csv"
	"fmt"
	"os"
	"strings"

	"github.com/alexmerren/speedruncom-scraper/pkg/srcomv1"
	"github.com/buger/jsonparser"
)

const (
	usersFieldIndex    = 9
	examinerFieldIndex = 10
	maxRunsPerPage     = 200

	usersListOutputFileHeader = "#userID\n"
	usersDataOutputFileHeader = "#ID,name,signupDate,location,numRuns\n"
	usersRunsOutputFileHeader = "#ID,gameID,categoryID,levelID,date,primaryTime,platform,emulated,players,examiner,values,status,statusReason,verifiedDate\n"
)

func ProcessUsersList(leaderboardInputFile, usersListOutputFile *os.File) error {
	usersListOutputFile.WriteString(usersListOutputFileHeader)

	allUsers := make(map[string]struct{})
	reader := csv.NewReader(leaderboardInputFile)
	reader.Read()
	records, err := reader.ReadAll()
	if err != nil {
		return err
	}

	for _, record := range records {
		if record[usersFieldIndex] == "" || record[usersFieldIndex] == "," {
			continue
		}

		users := strings.Split(record[usersFieldIndex], ",")
		for _, user := range users {
			allUsers[user] = struct{}{}
			allUsers[record[examinerFieldIndex]] = struct{}{}
		}
	}

	for userID := range allUsers {
		usersListOutputFile.WriteString(fmt.Sprintf("%s\n", userID))
	}

	return nil
}

func ProcessUsersData(
	usersListInputFile,
	usersDataOutputFile,
	usersRunsOutputFile *os.File,
) error {
	usersDataOutputFile.WriteString(usersDataOutputFileHeader)
	usersRunsOutputFile.WriteString(usersRunsOutputFileHeader)

	scanner := bufio.NewScanner(usersListInputFile)
	scanner.Scan()
	for scanner.Scan() {
		userID := scanner.Text()
		userResponse, err := srcomv1.GetUser(userID)
		if err != nil {
			continue
		}

		numRuns, err := processUserRuns(userID, usersRunsOutputFile)
		if err != nil {
			return err
		}

		err = processUser(userID, numRuns, userResponse, usersDataOutputFile)
		if err != nil {
			return err
		}
	}

	if err := scanner.Err(); err != nil {
		return err
	}

	return nil
}

func processUser(userID string, numRuns int, response []byte, outputFile *os.File) error {
	userData, _, _, err := jsonparser.Get(response, "data")
	if err != nil {
		return err
	}

	userName, _ := jsonparser.GetString(userData, "names", "international")
	userSignup, _ := jsonparser.GetString(userData, "signup")
	userLocation, _ := jsonparser.GetString(userData, "location", "country", "code")

	outputFile.WriteString(fmt.Sprintf("%s,%q,%s,%s,%d\n", userID, userName, userSignup, userLocation, numRuns))

	return nil
}

func processUserRuns(userID string, outputFile *os.File) (int, error) {
	numRuns := 0
	currentPage := 0

	for {
		response, err := srcomv1.GetUserRuns(userID, currentPage)
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
			}, "players")
			players := strings.Join(playerIDArray, ",")

			runValuesArray := []string{}
			jsonparser.ObjectEach(value, func(key []byte, value []byte, dataType jsonparser.ValueType, offset int) error {
				runValuesArray = append(runValuesArray, fmt.Sprintf("%s=%s", string(key), string(value)))
				return nil
			}, "values")
			values := strings.Join(runValuesArray, ",")

			outputFile.WriteString(fmt.Sprintf("%s,%s,%s,%s,%s,%0.2f,%s,%t,%q,%s,%q,%s,%q,%s\n", runID, gameID, categoryID, levelID, date, primaryTime, platform, emulated, players, examiner, values, status, statusReason, verifiedDate))
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
