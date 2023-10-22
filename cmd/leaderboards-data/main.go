package main

import (
	"fmt"
	"os"

	"github.com/alexmerren/speedruncom-scraper/internal/filesystem"
	"github.com/alexmerren/speedruncom-scraper/internal/srcomv1"
)

const (
	allGameIDListV1              = "./data/v1/games-id-list.csv"
	leaderboardsOutputFilenameV1 = "./data/v1/leaderboards-data.csv"
)

func main() {
	if err := getLeaderboardDataV1(); err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		os.Exit(1)
	}
}

func getLeaderboardDataV1() error {
	inputFile, err := filesystem.OpenInputFile(allGameIDListV1)
	if err != nil {
		return err
	}
	defer inputFile.Close()

	leaderboardsOutputFile, err := filesystem.CreateOutputFile(leaderboardsOutputFilenameV1)
	if err != nil {
		return err
	}
	defer leaderboardsOutputFile.Close()

	err = srcomv1.ProcessLeaderboardsData(inputFile, leaderboardsOutputFile)
	if err != nil {
		return err
	}

	return nil
}
