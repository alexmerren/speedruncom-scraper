package main

import (
	"fmt"
	"os"

	"github.com/alexmerren/speedruncom-scraper/internal/filesystem"
	"github.com/alexmerren/speedruncom-scraper/internal/srcomv1"
)

const (
	allGameIDListV1              = "./data/v1/games-id-list.csv"
	gamesOutputFilenameV1        = "./data/v1/games-data.csv"
	categoriesOutputFilenameV1   = "./data/v1/categories-data.csv"
	levelsOutputFilenameV1       = "./data/v1/levels-data.csv"
	variablesOutputFileV1        = "./data/v1/variables-data.csv"
	valuesOutputFileV1           = "./data/v1/values-data.csv"
	LeaderboardsOutputFilenameV1 = "./data/v1/leaderboards-data.csv"
)

func main() {
	if err := getGameAndLeaderboardDataV1(); err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		os.Exit(1)
	}
}

func getGameAndLeaderboardDataV1() error {
	inputFile, err := filesystem.OpenInputFile(allGameIDListV1)
	if err != nil {
		return err
	}
	defer inputFile.Close()

	gamesOutputFile, err := filesystem.CreateOutputFile(gamesOutputFilenameV1)
	if err != nil {
		return err
	}
	defer gamesOutputFile.Close()

	categoriesOutputFile, err := filesystem.CreateOutputFile(categoriesOutputFilenameV1)
	if err != nil {
		return err
	}
	defer categoriesOutputFile.Close()

	levelsOutputFile, err := filesystem.CreateOutputFile(levelsOutputFilenameV1)
	if err != nil {
		return err
	}
	defer levelsOutputFile.Close()

	variablesOutputFile, err := filesystem.CreateOutputFile(variablesOutputFileV1)
	if err != nil {
		return err
	}
	defer variablesOutputFile.Close()

	valuesOutputFile, err := filesystem.CreateOutputFile(valuesOutputFileV1)
	if err != nil {
		return err
	}
	defer valuesOutputFile.Close()

	LeaderboardsOutputFile, err := filesystem.CreateOutputFile(LeaderboardsOutputFilenameV1)
	if err != nil {
		return err
	}
	defer LeaderboardsOutputFile.Close()

	err = srcomv1.ProcessLeaderboardsAndGamesData(inputFile, gamesOutputFile, categoriesOutputFile, levelsOutputFile, variablesOutputFile, valuesOutputFile, LeaderboardsOutputFile)
	if err != nil {
		return err
	}

	return nil
}
