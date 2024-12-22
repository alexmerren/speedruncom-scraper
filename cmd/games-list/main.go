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
	gamesIdListFile, gamesIdListFileCloseFunc, err := repository.NewWriteRepository(repository.GamesIdListFilename)
	if err != nil {
		return err
	}
	defer gamesIdListFileCloseFunc()

	gamesListProcessor := &processor.GamesListProcessor{
		GamesIdListFile: gamesIdListFile,
		Client:          srcom_api.SrcomClient,
	}

	return gamesListProcessor.Process()
}
