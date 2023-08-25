package main

import (
	"bufio"
	"fmt"

	"github.com/alexmerren/speedruncom-scraper/internal/filesystem"
)

const (
	allGameIDListV1             = "./data/v1/games-id-list.csv"
	leaderboardOutputFilenameV1 = "./data/v1/leaderboard-data.csv"

	allGameIDListV2             = "./data/v2/games-id-list.csv"
	leaderboardOutputFilenameV2 = "./data/v2/leaderboard-data.csv"
)

func main() {
	//getLeaderboardDataV1()
}

func getLeaderboardDataV1() {
	inputFile, err := filesystem.OpenInputFile(allGameIDListV1)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer inputFile.Close()

	leaderboardOuptutFile, err := filesystem.CreateOutputFile(leaderboardOutputFilenameV1)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer leaderboardOuptutFile.Close()
	leaderboardOuptutFile.WriteString("#ID,name,URL,type,rules,releaseDate,addedDate,runCount,playerCount,numCategories,numLevels,emulator\n")

	scanner := bufio.NewScanner(inputFile)
	scanner.Scan()
	for scanner.Scan() {
		fmt.Println(scanner.Text())
	}
}

func getLeaderboardDataV2() {
	inputFile, err := filesystem.OpenInputFile(allGameIDListV2)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer inputFile.Close()

	leaderboardOuptutFile, err := filesystem.CreateOutputFile(leaderboardOutputFilenameV2)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer leaderboardOuptutFile.Close()
	leaderboardOuptutFile.WriteString("")

	scanner := bufio.NewScanner(inputFile)
	scanner.Scan()
	for scanner.Scan() {
		fmt.Println(scanner.Text())
	}
}
