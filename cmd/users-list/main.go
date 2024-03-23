package main

import (
	"errors"
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/alexmerren/speedruncom-scraper/internal"
)

const (
	usersFieldIndex    = 9
	examinerFieldIndex = 10
)

func main() {
	if err := generateUsersListV1(); err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		os.Exit(1)
	}
}

func generateUsersListV1() error {
	leaderboardsDataFile, closeFunc, _ := internal.NewCsvReader(internal.LeaderboardsDataFilenameV1)
	defer closeFunc()

	usersIdListFile, closeFunc, _ := internal.NewCsvWriter(internal.UsersIdListFilenameV1)
	usersIdListFile.Write(internal.FileHeaders[internal.UsersIdListFilenameV1])
	defer closeFunc()

	allUsers := make(map[string]struct{})
	leaderboardsDataFile.Read()

	for {
		record, err := leaderboardsDataFile.Read()
		if err != nil && errors.Is(err, io.EOF) {
			break
		}

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
		usersIdListFile.Write([]string{userID})
	}

	return nil
}
