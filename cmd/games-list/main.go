package main

import (
	"fmt"

	"github.com/alexmerren/speedruncom-scraper/internal/srcomv1"
	"github.com/alexmerren/speedruncom-scraper/internal/srcomv2"
	"github.com/buger/jsonparser"
)

func main() {
	getGameListV2()
}

func getGameListV1() {
	currentPage := 1
	request, err := srcomv1.GetGameList(currentPage)
	if err != nil {
		fmt.Println(err)
		return
	}

	size, err := jsonparser.GetInt(request, "pagination", "size")
	if err != nil {
		fmt.Println(err)
		return
	}

	gameIds := make([]string, 0)

	for size == 1000 {
		_, err := jsonparser.ArrayEach(request, func(value []byte, dataType jsonparser.ValueType, offset int, err error) {
			gameId, _ := jsonparser.GetString(value, "id")
			gameIds = append(gameIds, gameId)
		}, "data")
		if err != nil {
			fmt.Println(err)
			return
		}

		currentPage += 1
		request, _ = srcomv1.GetGameList(currentPage)
		size, _ = jsonparser.GetInt(request, "pagination", "size")
	}

	fmt.Println(gameIds)
}

func getGameListV2() {
	currentPage := 1
	request, _ := srcomv2.GetGameList(currentPage)
	lastPage, err := jsonparser.GetInt(request, "pagination", "pages")
	if err != nil {
		fmt.Println(err)
		return
	}

	gameIds := make([]string, 0)

	for int64(currentPage) < lastPage {
		_, err := jsonparser.ArrayEach(request, func(value []byte, dataType jsonparser.ValueType, offset int, err error) {
			gameId, _ := jsonparser.GetString(value, "id")
			gameIds = append(gameIds, gameId)
		}, "gameList")
		if err != nil {
			fmt.Println(err)
			return
		}
		currentPage += 1
		request, _ = srcomv2.GetGameList(currentPage)
	}

	fmt.Println(gameIds)
}