package processor

import (
	"github.com/alexmerren/speedruncom-scraper/internal/repository"
	"github.com/alexmerren/speedruncom-scraper/internal/srcom_api"
	"github.com/buger/jsonparser"
)

type GamesListProcessor struct {
	GamesIdListFile *repository.WriteRepository
	Client          *srcom_api.SrcomV1Client
}

func (p *GamesListProcessor) Process() error {
	currentPage := 0

	for {
		response, err := p.Client.GetGameList(currentPage)
		if err != nil {
			return err
		}

		_, err = jsonparser.ArrayEach(response, func(value []byte, dataType jsonparser.ValueType, offset int, _ error) {
			gameID, _ := jsonparser.GetString(value, "id")
			err = p.GamesIdListFile.Write([]string{gameID})
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
