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
	usersFieldIndex = 8

	usersListOutputFileHeader          = "#userID\n"
	usersDataOutputFileHeader          = "#ID,name,signupDate,location,numPersonalBests\n"
	usersPersonalBestsOutputFileHeader = "#userID,runID,game,category,level,values,place\n"
)

func ProcessUsersList(gameListInputFile, usersListOutputFile *os.File) error {
	allUsers := make(map[string]struct{})
	reader := csv.NewReader(gameListInputFile)
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
	usersPersonalBestsOutputFile *os.File,
) error {
	usersDataOutputFile.WriteString(usersDataOutputFileHeader)
	usersPersonalBestsOutputFile.WriteString(usersPersonalBestsOutputFileHeader)

	scanner := bufio.NewScanner(usersListInputFile)
	scanner.Scan()
	for scanner.Scan() {
		userID := scanner.Text()
		userResponse, err := srcomv1.GetUser(userID)
		if err != nil {
			return err
		}

		numPersonalBests, err := processUserPersonalBests(userID, userResponse, usersPersonalBestsOutputFile)
		if err != nil {
			return err
		}

		err = processUser(numPersonalBests, userResponse, usersDataOutputFile)
		if err != nil {
			return err
		}
	}

	if err := scanner.Err(); err != nil {
		return err
	}

	return nil
}

func processUser(numPersonalBests int, response []byte, outputFile *os.File) error {
	userData, _, _, err := jsonparser.Get(response, "data", "[0]", "players", "data", "[0]")
	if err != nil {
		return err
	}

	userID, _ := jsonparser.GetString(userData, "id")
	userName, _ := jsonparser.GetString(userData, "names", "international")
	userSignup, _ := jsonparser.GetString(userData, "signup")
	userLocation, _ := jsonparser.GetString(userData, "location", "country", "code")

	outputFile.WriteString(fmt.Sprintf("%s,%q,%s,%s,%d\n", userID, userName, userSignup, userLocation, numPersonalBests))

	return nil
}

func processUserPersonalBests(userID string, response []byte, outputFile *os.File) (int, error) {
	numPersonalBests := 0
	_, err := jsonparser.ArrayEach(response, func(value []byte, dataType jsonparser.ValueType, offset int, err error) {
		numPersonalBests += 1
		runData, _, _, _ := jsonparser.Get(value, "run")
		runPlace, _ := jsonparser.GetInt(value, "place")
		runID, _ := jsonparser.GetString(runData, "id")
		runGame, _ := jsonparser.GetString(runData, "game")
		runCategory, _ := jsonparser.GetString(runData, "category")
		runLevel, _ := jsonparser.GetString(runData, "level")
		runValuesArray := []string{}
		jsonparser.ObjectEach(runData, func(key []byte, value []byte, dataType jsonparser.ValueType, offset int) error {
			runValuesArray = append(runValuesArray, fmt.Sprintf("%s=%s", string(key), string(value)))
			return nil
		}, "values")
		runValues := strings.Join(runValuesArray, ",")

		outputFile.WriteString(fmt.Sprintf("%s,%s,%s,%s,%s,%q,%d\n", userID, runID, runGame, runCategory, runLevel, runValues, runPlace))
	}, "data")
	if err != nil {
		return 0, err
	}

	return numPersonalBests, nil
}
