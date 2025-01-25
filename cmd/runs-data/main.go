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
	usersIdListFile, usersIdListFileCloseFunc, err := repository.NewReadRepository(repository.UsersIdListFilename)
	if err != nil {
		return err
	}
	defer usersIdListFileCloseFunc()

	runsFile, runsFileCloseFunc, err := repository.NewWriteRepository(repository.RunsDataFilename)
	if err != nil {
		return err
	}
	defer runsFileCloseFunc()

	runsDataProcessor := &processor.RunsDataProcessor{
		UsersIdListFile: usersIdListFile,
		RunsFile:        runsFile,
		Client:          srcom_api.DefaultV1Client,
	}

	return runsDataProcessor.Process()
}
