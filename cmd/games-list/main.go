package main

import (
	"fmt"
	"sync"

	"github.com/alexmerren/speedruncom-scraper/internal/filesystem"
	"github.com/alexmerren/speedruncom-scraper/internal/srcomv1"
	"github.com/alexmerren/speedruncom-scraper/internal/srcomv2"
	"github.com/buger/jsonparser"
)

const (
	maxSizeAPIv1 = 1000

	outputFilenameV1 = "./data/v1/games-id-list.csv"
	outputFilenameV2 = "./data/v2/games-id-list.csv"
)

func main() {
	wg := &sync.WaitGroup{}

	wg.Add(1)
	go func() {
		defer wg.Done()
		getGameListV1()
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()
		getGameListV2()
	}()

	wg.Wait()
}

//nolint:errcheck// Not worth checking for an error for every file write -- that's the whole point of the file.
func getGameListV1() {
	outputFile, err := filesystem.CreateOutputFile(outputFilenameV1)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer outputFile.Close()
	outputFile.WriteString("#gameID\n")
	currentPage := 1

	for {
		request, _ := srcomv1.GetGameList(currentPage)
		_, err := jsonparser.ArrayEach(request, func(value []byte, dataType jsonparser.ValueType, offset int, err error) {
			gameID, _ := jsonparser.GetString(value, "id")
			outputFile.WriteString(fmt.Sprintf("%s\n", gameID))
		}, "data")
		if err != nil {
			fmt.Println(err)
			return
		}

		// The pagination size should always be at 1000, unless we get to last page then
		// it will be some random integer such that: 0 <= x <= 1000.
		size, _ := jsonparser.GetInt(request, "pagination", "size")
		if size < maxSizeAPIv1 {
			return
		}

		currentPage += 1
	}
}

//nolint:errcheck// Not worth checking for an error for every file write -- that's the whole point of the file.
func getGameListV2() {
	currentPage := 1
	request, _ := srcomv2.GetGameList(currentPage)
	lastPage, err := jsonparser.GetInt(request, "pagination", "pages")
	if err != nil {
		fmt.Println(err)
		return
	}

	outputFile, err := filesystem.CreateOutputFile(outputFilenameV2)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer outputFile.Close()
	outputFile.WriteString("#gameID\n")

	for int64(currentPage) <= lastPage {
		_, err := jsonparser.ArrayEach(request, func(value []byte, dataType jsonparser.ValueType, offset int, err error) {
			gameID, _ := jsonparser.GetString(value, "id")
			outputFile.WriteString(fmt.Sprintf("%s\n", gameID))
		}, "gameList")
		if err != nil {
			fmt.Println(err)
			return
		}
		currentPage += 1
		request, _ = srcomv2.GetGameList(currentPage)
	}
}
