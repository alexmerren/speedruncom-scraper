package main

import (
	"bufio"
	"fmt"

	"github.com/alexmerren/speedruncom-scraper/internal/filesystem"
	"github.com/alexmerren/speedruncom-scraper/internal/srcomv1"
	"github.com/alexmerren/speedruncom-scraper/internal/srcomv2"
)

const (
	allGameIDListV1                         = "./data/v1/games-id-list.csv"
	leaderboardOutputFilenameV1             = "./data/v1/leaderboard-data.csv"
	leaderboardCombinationsOutputFilenameV1 = "./data/v1/leaderboard-combinations-data.csv"

	allGameIDListV2             = "./data/v2/games-id-list.csv"
	leaderboardOutputFilenameV2 = "./data/v2/leaderboard-data.csv"
)

func main() {
	//getLeaderboardDataV1()
	response, _ := srcomv1.GetGame("76r55vd8")
	fmt.Println(string(response))
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
		response, err := srcomv2.GetGameData(scanner.Text())
		if err != nil {
			continue
		}
		fmt.Printf(string(response))
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
