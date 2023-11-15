package main

import (
	"fmt"
	"os"

	"github.com/alexmerren/speedruncom-scraper/internal/filesystem"
	"github.com/alexmerren/speedruncom-scraper/internal/srcomv1"
	"github.com/alexmerren/speedruncom-scraper/internal/srcomv2"
)

const (
	gamesListOutputFilenameV1 = "./data/v1/games-id-list.csv"
	gamesListOutputFilenameV2 = "./data/v2/games-id-list.csv"
)

func main() {
	if err := getGameListV1(); err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		os.Exit(1)
	}

	if err := getGameListV2(); err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		os.Exit(1)
	}
}

func getGameListV1() error {
	outputFile, err := filesystem.CreateOutputFile(gamesListOutputFilenameV1)
	if err != nil {
		return err
	}
	defer outputFile.Close()

	err = srcomv1.ProcessGamesList(outputFile)
	if err != nil {
		return err
	}

	return nil
}

func getGameListV2() error {
	outputFile, err := filesystem.CreateOutputFile(gamesListOutputFilenameV2)
	if err != nil {
		return err
	}
	defer outputFile.Close()

	err = srcomv2.ProcessGamesList(outputFile)
	if err != nil {
		return err
	}

	return nil
}
