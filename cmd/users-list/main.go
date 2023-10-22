package main

import (
	"fmt"
	"os"

	"github.com/alexmerren/speedruncom-scraper/internal/filesystem"
	"github.com/alexmerren/speedruncom-scraper/internal/srcomv1"
)

const (
	leaderboardDataFilenameV1 = "./data/v1/leaderboards-data.csv"
	usersListOutputFilenameV1 = "./data/v1/users-id-list.csv"
)

func main() {
	if err := getUsersListV1(); err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		os.Exit(1)
	}
}

func getUsersListV1() error {
	inputFile, err := filesystem.OpenInputFile(leaderboardDataFilenameV1)
	if err != nil {
		return err
	}
	defer inputFile.Close()

	usersListOutputFile, err := filesystem.CreateOutputFile(usersListOutputFilenameV1)
	if err != nil {
		return err
	}
	defer usersListOutputFile.Close()

	err = srcomv1.ProcessUsersList(inputFile, usersListOutputFile)
	if err != nil {
		return err
	}

	return nil
}
