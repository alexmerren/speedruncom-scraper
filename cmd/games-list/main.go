package main

import (
	"fmt"

	"github.com/alexmerren/speedruncom-scraper/internal/srcomv1"
	"github.com/alexmerren/speedruncom-scraper/internal/srcomv2"
)

func main() {

}

func getGameListV1() {
	initialRequest, err := srcomv1.GetGameList(1)
	fmt.Println(initialRequest, err)
}

func getGameListV2() {
	initialRequest, err := srcomv2.GetGameList(1)
	fmt.Println(initialRequest, err)
}
