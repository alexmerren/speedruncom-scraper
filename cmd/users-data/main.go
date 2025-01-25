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

	usersFile, usersFileCloseFunc, err := repository.NewWriteRepository(repository.UsersDataFilename)
	if err != nil {
		return err
	}
	defer usersFileCloseFunc()

	usersListProcessor := &processor.UsersDataProcessor{
		UsersIdListFile: usersIdListFile,
		UsersFile:       usersFile,
		Client:          srcom_api.DefaultV1Client,
	}

	return usersListProcessor.Process()
}
