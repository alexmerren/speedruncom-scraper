package main

import (
	"errors"
	"os"
)

func main() {
	getLeaderboardDataV2()
}

func getLeaderboardDataV1() {}

func getLeaderboardDataV2() {
	// Deal with this stuff in leaderboards-data! \/
	// what should I do in terms of the breakdown of each number of runs to the game,category,variable,value combination?
	// What should I do for the values/variable for each game?
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
