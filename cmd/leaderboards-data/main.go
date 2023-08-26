package main

import (
	"bufio"
	"fmt"

	"github.com/alexmerren/speedruncom-scraper/internal/filesystem"
	"github.com/alexmerren/speedruncom-scraper/internal/srcomv1"
	"github.com/buger/jsonparser"
)

const (
	allGameIDListV1                         = "./data/v1/games-id-list.csv"
	leaderboardOutputFilenameV1             = "./data/v1/leaderboard-data.csv"
	leaderboardCombinationsOutputFilenameV1 = "./data/v1/leaderboard-combinations-data.csv"
)

func main() {
	getLeaderboardDataV1()
}

func getLeaderboardDataV1() {
	inputFile, err := filesystem.OpenInputFile(allGameIDListV1)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer inputFile.Close()

	leaderboardCombinationsOutputFile, err := filesystem.CreateOutputFile(leaderboardCombinationsOutputFilenameV1)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer leaderboardCombinationsOutputFile.Close()
	leaderboardCombinationsOutputFile.WriteString("")

	leaderboardOuptutFile, err := filesystem.CreateOutputFile(leaderboardOutputFilenameV1)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer leaderboardOuptutFile.Close()
	leaderboardOuptutFile.WriteString("")

	scanner := bufio.NewScanner(inputFile)
	scanner.Scan()
	for scanner.Scan() {
		response, err := srcomv1.GetGame(scanner.Text())
		if err != nil {
			fmt.Println(err)
			continue
		}

		// Step 1. Go through each 'per-game' category.
		_, err = jsonparser.ArrayEach(response, func(value []byte, dataType jsonparser.ValueType, offset int, err error) {
			// categoryID, _, _, _ := jsonparser.Get(value, "id")
			// categoryType, _, _, _ := jsonparser.Get(value, "type")
			// if categoryType == "per-game" {
			// }

			// if categoryType == "per-level" {
			// }
		}, "data", "categories", "data")
	}
}
