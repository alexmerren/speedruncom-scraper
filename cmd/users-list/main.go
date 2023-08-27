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
	usersListOutputFile.WriteString("#userID")

	reader := csv.NewReader(inputFile)
	reader.Read()
	for {
		record, err := reader.Read()
		if err != nil {
			return
		}

		playersString := strings.ReplaceAll(record[8], ",", "\n")
		usersListOutputFile.WriteString(fmt.Sprintf("%s\n", playersString))
	}
}
