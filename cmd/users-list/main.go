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
	leaderboardsFile, leaderboardsFileCloseFunc, err := repository.NewReadRepository(repository.LeaderboardsDataFilename)
	if err != nil {
		return err
	}
	defer leaderboardsFileCloseFunc()

	usersIdListFile, usersIdListFileCloseFunc, err := repository.NewWriteRepository(repository.UsersIdListFilename)
	if err != nil {
		return err
	}
	defer usersIdListFileCloseFunc()

	usersListProcessor := &processor.UsersListProcessor{
		LeaderboardsFile: leaderboardsFile,
		UsersIdListFile:  usersIdListFile,
		Client:           srcom_api.DefaultV1Client,
	}

	return usersListProcessor.Process()
}
