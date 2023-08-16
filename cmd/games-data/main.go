package main

import (
	"bufio"
	"errors"
	"fmt"
	"os"

	"github.com/alexmerren/speedruncom-scraper/internal/srcomv2"
)

const (
	allGameIDListV1  = "./data/games-id-list-v1.csv"
	outputFilenameV1 = "./data/games-data-v1.csv"

	allGameIDListV2  = "./data/games-id-list-v2.csv"
	outputFilenameV2 = "./data/games-data-v2.csv"
)

func main() {
	gameID := "76r55vd8"
	response, _ := srcomv2.GetGameData(gameID)
	fmt.Println(string(response))
	// game.id
	// game.name
	// game.url
	// game.type
	// game.emulator
	// game.releaseDate
	// game.addedDate
	// game.runCount
	// game.totalPlayerCount
	// game.rules
	// game.platforms
	// game.categories (count)
	// game.categories.id (each)
	// game.categories.name (each)
	// game.categories.rules (each)
	// game.categories.numPlayers (each)
	// what should I do in terms of the breakdown of each number of runs to the game,category,variable,value combination?
	// What should I do for the values/variable for each game?

}

func getGameDataV1() {
	inputFile, err := openInputFile(allGameIDListV1)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer inputFile.Close()

	outputFile, err := createOutputFile(outputFilenameV1)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer outputFile.Close()

	scanner := bufio.NewScanner(inputFile)
	for scanner.Scan() {
		// TODO: Insert logic of getting data, formatting the data, and writing to output file.
		fmt.Println(scanner.Text())
	}

	if err := scanner.Err(); err != nil {
		fmt.Println(err)
		return
	}
}

func getGameDataV2() {
	inputFile, err := openInputFile(allGameIDListV2)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer inputFile.Close()

	outputFile, err := createOutputFile(outputFilenameV2)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer outputFile.Close()

	scanner := bufio.NewScanner(inputFile)
	for scanner.Scan() {
		// TODO: Insert logic of getting data, formatting the data, and writing to output file.
		fmt.Println(scanner.Text())
	}

	if err := scanner.Err(); err != nil {
		fmt.Println(err)
		return
	}
}

func openInputFile(filename string) (*os.File, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}

	return file, err
}

func createOutputFile(filename string) (*os.File, error) {
	if _, err := os.Stat(filename); errors.Is(err, os.ErrNotExist) {
		if _, err := os.Create(filename); err != nil {
			return nil, err
		}
	}

	outputFile, err := os.OpenFile(filename, os.O_RDWR, 0)
	if err != nil {
		return nil, err
	}
	return outputFile, nil
}
