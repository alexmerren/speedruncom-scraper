package main

import (
	"fmt"
	"os"

	"github.com/alexmerren/speedruncom-scraper/internal"
	"github.com/buger/jsonparser"
)

func main() {
	if err := generateGamesListV1(); err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		os.Exit(1)
	}
}

func generateGamesListV1() error {
	client := internal.NewSrcomV1Client()

	gamesIdListFile, closeFunc, _ := internal.NewCsvWriter(internal.GamesIdListFilenameV1)
	gamesIdListFile.Write(internal.FileHeaders[internal.GamesIdListFilenameV1])
	defer closeFunc()

	currentPage := 0

	for {
		response, err := client.GetGameList(currentPage)
		if err != nil {
			return err
		}

		_, err = jsonparser.ArrayEach(response, func(value []byte, dataType jsonparser.ValueType, offset int, err error) {
			gameID, _ := jsonparser.GetString(value, "id")
			gamesIdListFile.Write([]string{gameID})
		}, "data")
		if err != nil {
			return err
		}

		size, _ := jsonparser.GetInt(response, "pagination", "size")
		if size < 1000 {
			break
		}

		currentPage += 1
	}

	return nil
}
