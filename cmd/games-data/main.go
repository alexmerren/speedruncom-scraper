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

	gamesFile, gamesFileCloseFunc, err := repository.NewWriteRepository(repository.GamesDataFilename)
	if err != nil {
		return err
	}
	defer gamesFileCloseFunc()

	categoriesFile, categoriesFileCloseFunc, err := repository.NewWriteRepository(repository.CategoriesDataFilename)
	if err != nil {
		return err
	}
	defer categoriesFileCloseFunc()

	levelsFile, levelsFileCloseFunc, err := repository.NewWriteRepository(repository.LevelsDataFilename)
	if err != nil {
		return err
	}
	defer levelsFileCloseFunc()

	variablesFile, variablesFileCloseFunc, err := repository.NewWriteRepository(repository.VariablesDataFilename)
	if err != nil {
		return err
	}
	defer variablesFileCloseFunc()

	valuesFile, valuesFileCloseFunc, err := repository.NewWriteRepository(repository.ValuesDataFilename)
	if err != nil {
		return err
	}
	defer valuesFileCloseFunc()

	gamesDataProcessor := &processor.GamesDataProcessor{
		GamesIdListFile: gamesIdListFile,
		GamesFile:       gamesFile,
		CategoriesFile:  categoriesFile,
		LevelsFile:      levelsFile,
		VariablesFile:   variablesFile,
		ValuesFile:      valuesFile,
		Client:          srcom_api.NewSrcomClient(),
	}

	return gamesDataProcessor.Process()
}
