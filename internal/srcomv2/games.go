package srcomv2

import (
	"encoding/csv"
	"os"
	"strconv"

	"github.com/alexmerren/speedruncom-scraper/internal/filesystem"
	"github.com/alexmerren/speedruncom-scraper/pkg/srcomv2"
	"github.com/buger/jsonparser"
)

const (
	gameListOutputFileHeader = "#ID,name,URL,type,releaseDate,addedDate,runCount,playerCount,rules\n"
)

func ProcessGamesList(gameListOutputFile *os.File) error {
	gameListOutputFile.WriteString(gameListOutputFileHeader)
	gameListCsvWriter := csv.NewWriter(gameListOutputFile)
	defer gameListCsvWriter.Flush()

	currentPage := 0
	request, _ := srcomv2.GetGameList(currentPage)
	lastPage, err := jsonparser.GetInt(request, "pagination", "pages")
	if err != nil {
		return err
	}

	for int64(currentPage) <= lastPage {
		_, err := jsonparser.ArrayEach(request, func(value []byte, dataType jsonparser.ValueType, offset int, err error) {
			id, _ := jsonparser.GetString(value, "id")
			name, _ := jsonparser.GetString(value, "name")
			url, _ := jsonparser.GetString(value, "url")
			gameType, _ := jsonparser.GetString(value, "type")
			releaseDate, _ := jsonparser.GetInt(value, "releaseDate")
			addedDate, _ := jsonparser.GetInt(value, "addedDate")
			runCount, _ := jsonparser.GetInt(value, "runCount")
			playerCount, _ := jsonparser.GetInt(value, "totalPlayerCount")
			rules, _ := jsonparser.GetString(value, "rules")

			gameListCsvWriter.Write([]string{
				id,
				name,
				url,
				gameType,
				strconv.Itoa(int(releaseDate)),
				strconv.Itoa(int(addedDate)),
				strconv.Itoa(int(runCount)),
				strconv.Itoa(int(playerCount)),
				filesystem.FormatStringForCsv(rules),
			})
		}, "gameList")
		if err != nil {
			return err
		}

		currentPage += 1

		request, err = srcomv2.GetGameList(currentPage)
		if err != nil {
			return err
		}
	}

	return nil
}
