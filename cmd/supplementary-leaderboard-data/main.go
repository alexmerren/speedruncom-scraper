package main

import (
	"errors"
	"flag"
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
	gameIdFlag := flag.String("gameId", "", "Game ID i.e. y65797de")
	flag.Parse()

	if *gameIdFlag == "" {
		return errors.New("gameId must be provided")
	}

	leaderboardsFile, leaderboardsFileCloseFunc, err := repository.NewAppendRepository(repository.SupplementaryLeaderboardDataFilename)
	if err != nil {
		return err
	}
	defer leaderboardsFileCloseFunc()

	additionalLeaderboardsDataProcessor := &processor.SupplementaryLeaderboardDataProcessor{
		GameId:                       *gameIdFlag,
		SupplementaryLeaderboardFile: leaderboardsFile,
		Client:                       srcom_api.DefaultV1Client,
	}

	return additionalLeaderboardsDataProcessor.Process()
}
