package main

import (
	"fmt"
	"os"
	"strings"

	"github.com/alexmerren/speedruncom-scraper/internal/srcomv1"
	"github.com/alexmerren/speedruncom-scraper/internal/srcomv2"
	"github.com/buger/jsonparser"
)

const (
	maxSizeAPIv1 = 1000
	outputFile   = "../data/games-id-list.csv"
)

func main() {
	getGameListV1()
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

	// TODO: If file does not exist, create it!
	file, err := os.Open(outputFile)
	if err != nil {
		fmt.Println(err)
		return
	}

	file.Write([]byte("Game ID\n"))

	for size == maxSizeAPIv1 {
		_, err := jsonparser.ArrayEach(request, func(value []byte, dataType jsonparser.ValueType, offset int, err error) {
			gameId, _ := jsonparser.GetString(value, "id")
			output := strings.Join([]string{gameId, "\n"}, "")
			file.Write([]byte(output))
		}, "data")
		if err != nil {
			fmt.Println(err)
			return
		}

		currentPage += 1
		request, _ = srcomv1.GetGameList(currentPage)
		size, _ = jsonparser.GetInt(request, "pagination", "size")
	}
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

	// TODO: Handle the gameIds. Write to file? Insert into database?
	fmt.Println(gameIds)
}
