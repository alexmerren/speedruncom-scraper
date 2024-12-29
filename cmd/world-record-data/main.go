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

	worldRecordFile, worldRecordFileCloseFunc, err := repository.NewWriteRepository(repository.WorldRecordDataFilename)
	if err != nil {
		return err
	}
	defer worldRecordFileCloseFunc()

	worldRecordDataProcessor := &processor.WorldRecordDataProcessor{
		GameId:              *gameIdFlag,
		WorldRecordDataFile: worldRecordFile,
		ClientV1:            srcom_api.DefaultV1Client,
		ClientV2:            srcom_api.DefaultV2Client,
	}

	return worldRecordDataProcessor.Process()
}
