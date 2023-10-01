package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/alexmerren/speedruncom-scraper/internal/filesystem"
	"github.com/alexmerren/speedruncom-scraper/internal/srcomv1"
	"github.com/buger/jsonparser"
)

const (
	allUserIDListV1                   = "./data/v1/users-id-list.csv"
	userOutputFilenameV1              = "./data/v1/users-data.csv"
	userPersonalBestsOutputFilenameV1 = "./data/v1/users-personal-bests-data.csv"
)

func main() {
	getUsersDataV1()
}

//nolint:errcheck // Don't need to check for errors.
func getUsersDataV1() {
	inputFile, err := filesystem.OpenInputFile(allUserIDListV1)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer inputFile.Close()

	userOutputFile, err := filesystem.CreateOutputFile(userOutputFilenameV1)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer userOutputFile.Close()
	userOutputFile.WriteString("#ID,name,signupDate,location,numPersonalBests\n")

	userPersonalBestsOutputFile, err := filesystem.CreateOutputFile(userPersonalBestsOutputFilenameV1)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer userPersonalBestsOutputFile.Close()
	userPersonalBestsOutputFile.WriteString("#userID,runID,game,category,level,values,place\n")

	scanner := bufio.NewScanner(inputFile)
	scanner.Scan()
	for scanner.Scan() {
		userID := scanner.Text()
		userResponse, err := srcomv1.GetUser(userID)
		if err != nil {
			fmt.Println(err)
			return
		}

		numPersonalBests, err := processUserPersonalBests(userID, userResponse, userPersonalBestsOutputFile)
		if err != nil {
			fmt.Println(err)
			return
		}

		err = processUser(numPersonalBests, userResponse, userOutputFile)
		if err != nil {
			fmt.Println(err)
			return
		}
	}
}

//nolint:errcheck // Don't need to check for errors.
func processUser(numPersonalBests int, response []byte, outputFile *os.File) error {
	userData, _, _, _ := jsonparser.Get(response, "data", "[0]", "players", "data", "[0]")
	userID, _ := jsonparser.GetString(userData, "id")
	userName, _ := jsonparser.GetString(userData, "names", "international")
	userSignup, _ := jsonparser.GetString(userData, "signup")
	userLocation, _ := jsonparser.GetString(userData, "location", "country", "code")
	outputFile.WriteString(fmt.Sprintf("%s,\"%s\",%s,%s,%d\n", userID, userName, userSignup, userLocation, numPersonalBests))
	return nil
}

//nolint:errcheck // Don't need to check for errors.
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

		outputFile.WriteString(fmt.Sprintf("%s,%s,%s,%s,%s,\"%s\",%d\n", userID, runID, runGame, runCategory, runLevel, runValues, runPlace))
	}, "data")
	if err != nil {
		return 0, err
	}
	return numPersonalBests, nil
}
