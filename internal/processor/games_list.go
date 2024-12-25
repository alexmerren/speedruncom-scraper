package processor

import (
	"strconv"

	"github.com/alexmerren/speedruncom-scraper/internal/repository"
	"github.com/alexmerren/speedruncom-scraper/internal/srcom_api"
	"github.com/buger/jsonparser"
)

type GamesListProcessor struct {
	GamesIdListFile *repository.WriteRepository
	GamesFileV2     *repository.WriteRepository
	Client          *srcom_api.SrcomV1Client
	ClientV2        *srcom_api.SrcomV2Client
}

func (p *GamesListProcessor) Process() error {
	err := p.processV1()
	if err != nil {
		return err
	}

	err = p.processV2()
	if err != nil {
		return err
	}

	return nil
}

func (p *GamesListProcessor) processV1() error {
	currentPage := 0

	for {
		responseV1, err := p.Client.GetGameList(currentPage)
		if err != nil {
			return err
		}

		_, err = jsonparser.ArrayEach(responseV1, func(value []byte, dataType jsonparser.ValueType, offset int, _ error) {
			gameID, _ := jsonparser.GetString(value, "id")
			err = p.GamesIdListFile.Write([]string{gameID})
		}, "data")
		if err != nil {
			return err
		}

		size, _ := jsonparser.GetInt(responseV1, "pagination", "size")
		if size < 1000 {
			break
		}

		currentPage += 1
	}

	return nil
}

func (p *GamesListProcessor) processV2() error {
	currentPage := 0
	request, _ := p.ClientV2.GetGameList(currentPage)
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

			p.GamesFileV2.Write([]string{
				id,
				name,
				url,
				gameType,
				strconv.Itoa(int(releaseDate)),
				strconv.Itoa(int(addedDate)),
				strconv.Itoa(int(runCount)),
				strconv.Itoa(int(playerCount)),
				repository.FormatCsvString(rules),
			})
		}, "gameList")
		if err != nil {
			return err
		}

		currentPage += 1

		request, err = p.ClientV2.GetGameList(currentPage)
		if err != nil {
			return err
		}
	}

	return nil
}
