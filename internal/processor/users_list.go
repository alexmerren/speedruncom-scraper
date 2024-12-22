package processor

import (
	"errors"
	"io"
	"strings"

	"github.com/alexmerren/speedruncom-scraper/internal/repository"
	"github.com/alexmerren/speedruncom-scraper/internal/srcom_api"
)

const (
	usersFieldIndex    = 9
	examinerFieldIndex = 10
)

type UsersListProcessor struct {
	LeaderboardsFile *repository.ReadRepository
	UsersIdListFile  *repository.WriteRepository
	Client           *srcom_api.SrcomV1Client
}

func (p *UsersListProcessor) Process() error {
	allUsers := make(map[string]struct{})
	p.LeaderboardsFile.Read()

	for {
		record, err := p.LeaderboardsFile.Read()
		if err != nil && errors.Is(err, io.EOF) {
			break
		}

		if record[usersFieldIndex] == "" || record[usersFieldIndex] == "," {
			continue
		}

		users := strings.Split(record[usersFieldIndex], ",")
		for _, user := range users {
			allUsers[user] = struct{}{}
			allUsers[record[examinerFieldIndex]] = struct{}{}
		}
	}

	for userID := range allUsers {
		p.UsersIdListFile.Write([]string{userID})
	}

	return nil
}
