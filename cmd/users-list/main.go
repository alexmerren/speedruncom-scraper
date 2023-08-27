package main

import (
	"encoding/csv"
	"fmt"
	"strings"

	"github.com/alexmerren/speedruncom-scraper/internal/filesystem"
)

const (
	leaderboardDataFilenameV1 = "./data/v1/leaderboards-data.csv"
	usersListOutputFilenameV1 = "./data/v1/users-id-list.csv"
	usersFieldIndex           = 8
)

func main() {
	getUsersListV1()
}

func getUsersListV1() {
	inputFile, err := filesystem.OpenInputFile(leaderboardDataFilenameV1)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer inputFile.Close()

	usersListOutputFile, err := filesystem.CreateOutputFile(usersListOutputFilenameV1)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer usersListOutputFile.Close()
	usersListOutputFile.WriteString("#userID\n")

	// We define void as to reduce the amount of memory to store all the user IDs.
	type void struct{}
	allUsers := make(map[string]void)

	// Call reader.Read() to not read the header line into the records variable.
	reader := csv.NewReader(inputFile)
	reader.Read()
	records, err := reader.ReadAll()
	if err != nil {
		fmt.Println(err)
		return
	}

	for _, record := range records {
		if record[usersFieldIndex] == "" || record[usersFieldIndex] == "," {
			continue
		}

		users := strings.Split(record[usersFieldIndex], ",")
		for _, user := range users {
			allUsers[user] = void{}
		}
	}

	for userID := range allUsers {
		usersListOutputFile.WriteString(fmt.Sprintf("%s\n", userID))
	}
}
