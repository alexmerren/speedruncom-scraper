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

	gamesFileV2, gamesFileV2CloseFunc, err := repository.NewWriteRepository(repository.GamesDataFilenameV2)
	if err != nil {
		return err
	}
	defer gamesFileV2CloseFunc()

	gamesListProcessor := &processor.GamesListProcessor{
		GamesIdListFile: gamesIdListFile,
		GamesFileV2:     gamesFileV2,
		Client:          srcom_api.DefaultV1Client,
		ClientV2:        srcom_api.DefaultV2Client,
	}

	return gamesListProcessor.Process()
}
