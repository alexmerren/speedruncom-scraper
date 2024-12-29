package processor

import (
	"errors"
	"io"
	"strconv"

	"github.com/alexmerren/speedruncom-scraper/internal/repository"
	"github.com/alexmerren/speedruncom-scraper/internal/srcom_api"
	"github.com/buger/jsonparser"
)

type GamesDataProcessor struct {
	GamesIdListFile *repository.ReadRepository
	GamesFile       *repository.WriteRepository
	CategoriesFile  *repository.WriteRepository
	LevelsFile      *repository.WriteRepository
	VariablesFile   *repository.WriteRepository
	ValuesFile      *repository.WriteRepository
	Client          *srcom_api.SrcomV1Client
}

func (p *GamesDataProcessor) Process() error {
	p.GamesIdListFile.Read()

	for {
		record, err := p.GamesIdListFile.Read()
		if err != nil && errors.Is(err, io.EOF) {
			break
		}

		gameId := record[0]
		response, err := p.Client.GetGame(gameId)
		if err != nil {
			continue
		}

		err = p.processCategory(response, gameId)
		if err != nil {
			return err
		}

		err = p.processLevel(response, gameId)
		if err != nil {
			return err
		}

		err = p.processVariableAndValue(response, gameId)
		if err != nil {
			return err
		}

		err = p.processGame(response, gameId)
		if err != nil {
			return err
		}
	}

	return nil
}

func (p *GamesDataProcessor) processCategory(gameResponse []byte, gameId string) error {
	_, err := jsonparser.ArrayEach(gameResponse, func(value []byte, dataType jsonparser.ValueType, offset int, err error) {
		categoryId, _ := jsonparser.GetString(value, "id")
		categoryName, _ := jsonparser.GetString(value, "name")
		categoryRules, _ := jsonparser.GetString(value, "rules")
		categoryNumPlayers, _ := jsonparser.GetInt(value, "players", "value")
		categoryType, _ := jsonparser.GetString(value, "type")

		err = p.CategoriesFile.Write([]string{
			gameId,
			categoryId,
			categoryName,
			strconv.Quote(categoryRules),
			categoryType,
			strconv.Itoa(int(categoryNumPlayers)),
		})
	}, "data", "categories", "data")

	return err
}

func (p *GamesDataProcessor) processLevel(gameResponse []byte, gameId string) error {
	_, err := jsonparser.ArrayEach(gameResponse, func(value []byte, dataType jsonparser.ValueType, offset int, err error) {
		levelId, _ := jsonparser.GetString(value, "id")
		levelName, _ := jsonparser.GetString(value, "name")
		levelRules, _ := jsonparser.GetString(value, "rules")

		err = p.LevelsFile.Write([]string{gameId, levelId, levelName, strconv.Quote(levelRules)})
	}, "data", "levels", "data")

	return err
}

func (p *GamesDataProcessor) processVariableAndValue(gameResponse []byte, gameId string) error {
	_, err := jsonparser.ArrayEach(gameResponse, func(value []byte, dataType jsonparser.ValueType, offset int, err error) {
		variableId, _ := jsonparser.GetString(value, "id")
		variableName, _ := jsonparser.GetString(value, "name")
		variableCategory, _ := jsonparser.GetString(value, "category")
		variableScope, _ := jsonparser.GetString(value, "scope", "type")
		variableScopeLevel, _ := jsonparser.GetString(value, "scope", "level")
		variableIsSubcategory, _ := jsonparser.GetBoolean(value, "is-subcategory")
		variableDefault, _ := jsonparser.GetString(value, "values", "default")

		err = p.VariablesFile.Write([]string{
			gameId,
			variableId,
			variableName,
			variableCategory,
			variableScope,
			variableScopeLevel,
			strconv.FormatBool(variableIsSubcategory),
			variableDefault,
		})

		err = jsonparser.ObjectEach(value, func(key []byte, value []byte, dataType jsonparser.ValueType, offset int) error {
			valueId := string(key)
			valueLabel, _ := jsonparser.GetString(value, "label")
			valueRules, _ := jsonparser.GetString(value, "rules")

			return p.ValuesFile.Write([]string{
				gameId,
				variableId,
				valueId,
				valueLabel,
				strconv.Quote(valueRules),
			})
		}, "values", "values")
	}, "data", "variables", "data")

	return err
}

func (p *GamesDataProcessor) processGame(gameResponse []byte, gameId string) error {
	gameData, _, _, err := jsonparser.Get(gameResponse, "data")
	if err != nil {
		return err
	}

	gameName, _ := jsonparser.GetString(gameData, "names", "international")
	gameURL, _ := jsonparser.GetString(gameData, "abbreviation")
	gameReleaseDate, _ := jsonparser.GetString(gameData, "release-date")
	gameCreatedDate, _ := jsonparser.GetString(gameData, "created")

	return p.GamesFile.Write([]string{
		gameId,
		gameName,
		gameURL,
		gameReleaseDate,
		gameCreatedDate,
	})
}
