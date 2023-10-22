package main

import (
	"fmt"
	"os"

	"github.com/alexmerren/speedruncom-scraper/internal/filesystem"
	"github.com/alexmerren/speedruncom-scraper/internal/srcomv1"
)

const (
	allUserIDListV1                    = "./data/v1/users-id-list.csv"
	usersDataOutputFilenameV1          = "./data/v1/users-data.csv"
	usersPersonalBestsOutputFilenameV1 = "./data/v1/users-personal-bests-data.csv"
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

	usersPersonalBestsOutputFile, err := filesystem.CreateOutputFile(usersPersonalBestsOutputFilenameV1)
	if err != nil {
		return err
	}
	defer usersPersonalBestsOutputFile.Close()

	err = srcomv1.ProcessUsersData(inputFile, usersDataOutputFile, usersPersonalBestsOutputFile)
	if err != nil {
		return err
	}

	return nil
}
