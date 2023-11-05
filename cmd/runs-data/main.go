package main

import (
	"fmt"
	"os"

	"github.com/alexmerren/speedruncom-scraper/internal/filesystem"
	"github.com/alexmerren/speedruncom-scraper/internal/srcomv1"
)

const (
	allUserIDListV1              = "./data/v1/users-id-list.csv"
	runsDataOutputFilenameV1     = "./data/v1/runs-data.csv"
	userRunsDataOutputFilenameV1 = "./data/v1/users-runs-data.csv"
)

func main() {
	if err := getUsersRunsV1(); err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		os.Exit(1)
	}
}

func getUsersRunsV1() error {
	inputFile, err := filesystem.OpenInputFile(allUserIDListV1)
	if err != nil {
		return err
	}
	defer inputFile.Close()

	runsDataOutputFile, err := filesystem.CreateOutputFile(runsDataOutputFilenameV1)
	if err != nil {
		return err
	}
	defer runsDataOutputFile.Close()

	userRunsDataOutputFile, err := filesystem.CreateOutputFile(userRunsDataOutputFilenameV1)
	if err != nil {
		return err
	}
	defer userRunsDataOutputFile.Close()

	err = srcomv1.ProcessRunsData(inputFile, runsDataOutputFile, userRunsDataOutputFile)
	if err != nil {
		return err
	}

	return nil
}
