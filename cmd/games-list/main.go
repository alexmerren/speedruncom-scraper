package main

import (
	"errors"
	"fmt"
	"os"
	"strings"

	"github.com/alexmerren/speedruncom-scraper/internal/srcomv1"
	"github.com/alexmerren/speedruncom-scraper/internal/srcomv2"
	"github.com/buger/jsonparser"
)

const (
	maxSizeAPIv1 = 1000

	outputFilenameV1 = "./data/games-id-list-v1.csv"
	outputFilenameV2 = "./data/games-id-list-v2.csv"
)

func main() {
	getGameListV1()
	getGameListV2()
}

//nolint:errcheck// Not worth checking for an error for every file write -- that's the whole point of the file.
func getGameListV1() {
	outputFile, err := createOutputFile(outputFilenameV1)
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
			gameId, _ := jsonparser.GetString(value, "id")
			output := strings.Join([]string{gameId, "\n"}, "")
			outputFile.WriteString(output)
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

	outputFile, err := createOutputFile(outputFilenameV2)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer outputFile.Close()
	outputFile.WriteString("#gameID\n")

	for int64(currentPage) <= lastPage {
		_, err := jsonparser.ArrayEach(request, func(value []byte, dataType jsonparser.ValueType, offset int, err error) {
			gameId, _ := jsonparser.GetString(value, "id")
			output := strings.Join([]string{gameId, "\n"}, "")
			outputFile.WriteString(output)
		}, "gameList")
		if err != nil {
			fmt.Println(err)
			return
		}
		currentPage += 1
		request, _ = srcomv2.GetGameList(currentPage)
	}
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
