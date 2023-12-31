package main

import (
	"fmt"
	"os"

	"github.com/alexmerren/speedruncom-scraper/internal/filesystem"
	"github.com/alexmerren/speedruncom-scraper/internal/srcomv1"
)

const (
	allUserIDListV1           = "./data/v1/users-id-list.csv"
	usersDataOutputFilenameV1 = "./data/v1/users-data.csv"
)

func main() {
	if err := getUsersDataV1(); err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		os.Exit(1)
	}
}

func getUsersDataV1() error {
	inputFile, err := filesystem.OpenInputFile(allUserIDListV1)
	if err != nil {
		return err
	}
	defer inputFile.Close()

	usersDataOutputFile, err := filesystem.CreateOutputFile(usersDataOutputFilenameV1)
	if err != nil {
		return err
	}
	defer usersDataOutputFile.Close()

	err = srcomv1.ProcessUsersData(inputFile, usersDataOutputFile)
	if err != nil {
		return err
	}

	return nil
}
