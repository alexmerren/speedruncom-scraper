package processor

import (
	"errors"
	"io"

	"github.com/alexmerren/speedruncom-scraper/internal/repository"
	"github.com/alexmerren/speedruncom-scraper/internal/srcom_api"
	"github.com/buger/jsonparser"
)

type UsersDataProcessor struct {
	UsersIdListFile *repository.ReadRepository
	UsersFile       *repository.WriteRepository
	Client          *srcom_api.SrcomV1Client
}

func (p *UsersDataProcessor) Process() error {
	p.UsersIdListFile.Read()

	for {
		record, err := p.UsersIdListFile.Read()
		if err != nil && errors.Is(err, io.EOF) {
			break
		}

		userId := record[0]
		response, err := p.Client.GetUser(userId)
		if err != nil {
			continue
		}

		userData, _, _, err := jsonparser.Get(response, "data")
		if err != nil {
			return err
		}
		userName, _ := jsonparser.GetString(userData, "names", "international")
		userSignup, _ := jsonparser.GetString(userData, "signup")
		userLocation, _ := jsonparser.GetString(userData, "location", "country", "code")

		err = p.UsersFile.Write([]string{userId, userName, userSignup, userLocation})
		if err != nil {
			return err
		}
	}

	return nil
}
