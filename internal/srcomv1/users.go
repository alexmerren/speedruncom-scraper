package srcomv1

import (
	"bufio"
	"encoding/csv"
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/alexmerren/speedruncom-scraper/pkg/srcomv1"
	"github.com/buger/jsonparser"
)

const (
	usersFieldIndex    = 9
	examinerFieldIndex = 10

	usersListOutputFileHeader = "userID\n"
	usersDataOutputFileHeader = "ID,name,signupDate,location,numRuns\n"
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
	usersDataOutputFile *os.File,
) error {
	usersDataOutputFile.WriteString(usersDataOutputFileHeader)
	usersDataCsvWriter := csv.NewWriter(usersDataOutputFile)
	defer usersDataCsvWriter.Flush()

	scanner := bufio.NewScanner(usersListInputFile)
	scanner.Scan()
	for scanner.Scan() {
		userID := scanner.Text()
		userResponse, err := srcomv1.GetUser(userID)
		if err != nil {
			continue
		}

		err = processUser(usersDataCsvWriter, userID, -1, userResponse)
		if err != nil {
			return err
		}
	}

	if err := scanner.Err(); err != nil {
		return err
	}

	return nil
}

func processUser(outputFile *csv.Writer, userID string, numRuns int, response []byte) error {
	userData, _, _, err := jsonparser.Get(response, "data")
	if err != nil {
		return err
	}

	userName, _ := jsonparser.GetString(userData, "names", "international")
	userSignup, _ := jsonparser.GetString(userData, "signup")
	userLocation, _ := jsonparser.GetString(userData, "location", "country", "code")

	if numRuns < 0 {
		outputFile.Write([]string{userID, userName, userSignup, userLocation})
	} else {
		outputFile.Write([]string{userID, userName, userSignup, userLocation, strconv.Itoa(numRuns)})
	}

	return nil
}
