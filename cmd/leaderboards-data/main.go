package main

import (
	"fmt"

	"github.com/alexmerren/speedruncom-scraper/internal/srcomv2"
)

func main() {
	gameID := "76r55vd8"
	response, _ := srcomv2.GetGameData(gameID)
	fmt.Println(string(response))
	// Deal with this stuff in leaderboards-data! \/
	// what should I do in terms of the breakdown of each number of runs to the game,category,variable,value combination?
	// What should I do for the values/variable for each game?
}
