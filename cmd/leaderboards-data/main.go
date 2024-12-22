package main

import (
	"fmt"
	"os"

	"github.com/alexmerren/speedruncom-scraper/internal/processor"
	"github.com/alexmerren/speedruncom-scraper/internal/repository"
	"github.com/alexmerren/speedruncom-scraper/internal/srcom_api"
)

func main() {
	if err := run(); err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		os.Exit(1)
	}
}

func run() error {
	gamesIdListFile, gamesIdListFileCloseFunc, err := repository.NewReadRepository(repository.GamesIdListFilename)
	if err != nil {
		return err
	}
	defer gamesIdListFileCloseFunc()

	leaderboardsFile, leaderboardsFileCloseFunc, err := repository.NewWriteRepository(repository.LeaderboardsDataFilename)
	if err != nil {
		return err
	}
	defer leaderboardsFileCloseFunc()

	leaderboardsDataProcessor := &processor.LeaderboardsDataProcessor{
		GamesIdListFile:  gamesIdListFile,
		LeaderboardsFile: leaderboardsFile,
		Client:           srcom_api.SrcomClient,
	}

	return leaderboardsDataProcessor.Process()
}
