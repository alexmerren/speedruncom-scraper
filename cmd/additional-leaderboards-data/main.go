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
	if len(os.Args) < 2 {
		return fmt.Errorf("Must provide a game ID i.e. 76r55vd8")
	}

	gameId := os.Args[1]
	leaderboardsFile, leaderboardsFileCloseFunc, err := repository.NewWriteRepository(repository.AdditionalLeaderboardsDataFilename)
	if err != nil {
		return err
	}
	defer leaderboardsFileCloseFunc()

	additionalLeaderboardsDataProcessor := &processor.AdditionalLeaderboardsDataProcessor{
		GameId:                     gameId,
		AdditionalLeaderboardsFile: leaderboardsFile,
		Client:                     srcom_api.NewSrcomClient(),
	}

	return additionalLeaderboardsDataProcessor.Process()
}
