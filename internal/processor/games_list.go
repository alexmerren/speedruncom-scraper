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
	DevelopersFile  *repository.WriteRepository
	GenresFile      *repository.WriteRepository
	PlatformsFile   *repository.WriteRepository
	PublishersFile  *repository.WriteRepository
	EnginesFile     *repository.WriteRepository
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

	err = p.processGenres()
	if err != nil {
		return err
	}

	err = p.processPlatforms()
	if err != nil {
		return err
	}

	err = p.processPublishers()
	if err != nil {
		return err
	}

	err = p.processDevelopers()
	if err != nil {
		return err
	}

	err = p.processEngines()
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
				strconv.Quote(rules),
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

func (p *GamesListProcessor) processDevelopers() error {
	currentPage := 0

	for {
		response, err := p.Client.GetDeveloperList(currentPage)
		if err != nil {
			return err
		}

		_, err = jsonparser.ArrayEach(response, func(value []byte, dataType jsonparser.ValueType, offset int, _ error) {
			developerId, _ := jsonparser.GetString(value, "id")
			name, _ := jsonparser.GetString(value, "name")
			err = p.DevelopersFile.Write([]string{developerId, name})
		}, "data")
		if err != nil {
			return err
		}

		size, _ := jsonparser.GetInt(response, "pagination", "size")
		if size < 200 {
			break
		}

		currentPage += 1
	}

	return nil
}

func (p *GamesListProcessor) processPublishers() error {
	currentPage := 0

	for {
		response, err := p.Client.GetPublisherList(currentPage)
		if err != nil {
			return err
		}

		_, err = jsonparser.ArrayEach(response, func(value []byte, dataType jsonparser.ValueType, offset int, _ error) {
			publisherId, _ := jsonparser.GetString(value, "id")
			name, _ := jsonparser.GetString(value, "name")
			err = p.PublishersFile.Write([]string{publisherId, name})
		}, "data")
		if err != nil {
			return err
		}

		size, _ := jsonparser.GetInt(response, "pagination", "size")
		if size < 200 {
			break
		}

		currentPage += 1
	}

	return nil
}

func (p *GamesListProcessor) processPlatforms() error {
	currentPage := 0

	for {
		response, err := p.Client.GetPlatformList(currentPage)
		if err != nil {
			return err
		}

		_, err = jsonparser.ArrayEach(response, func(value []byte, dataType jsonparser.ValueType, offset int, _ error) {
			platformId, _ := jsonparser.GetString(value, "id")
			name, _ := jsonparser.GetString(value, "name")
			releaseYear, _ := jsonparser.GetInt(value, "released")
			err = p.PlatformsFile.Write([]string{platformId, name, strconv.Itoa(int(releaseYear))})
		}, "data")
		if err != nil {
			return err
		}

		size, _ := jsonparser.GetInt(response, "pagination", "size")
		if size < 200 {
			break
		}

		currentPage += 1
	}

	return nil
}

func (p *GamesListProcessor) processGenres() error {
	currentPage := 0

	for {
		response, err := p.Client.GetGenreList(currentPage)
		if err != nil {
			return err
		}

		_, err = jsonparser.ArrayEach(response, func(value []byte, dataType jsonparser.ValueType, offset int, _ error) {
			genreId, _ := jsonparser.GetString(value, "id")
			name, _ := jsonparser.GetString(value, "name")
			err = p.GenresFile.Write([]string{genreId, name})
		}, "data")
		if err != nil {
			return err
		}

		size, _ := jsonparser.GetInt(response, "pagination", "size")
		if size < 200 {
			break
		}

		currentPage += 1
	}

	return nil
}

func (p *GamesListProcessor) processEngines() error {
	currentPage := 0

	for {
		response, err := p.Client.GetEngineList(currentPage)
		if err != nil {
			return err
		}

		_, err = jsonparser.ArrayEach(response, func(value []byte, dataType jsonparser.ValueType, offset int, _ error) {
			engineId, _ := jsonparser.GetString(value, "id")
			name, _ := jsonparser.GetString(value, "name")
			err = p.EnginesFile.Write([]string{engineId, name})
		}, "data")
		if err != nil {
			return err
		}

		size, _ := jsonparser.GetInt(response, "pagination", "size")
		if size < 200 {
			break
		}

		currentPage += 1
	}

	return nil
}
