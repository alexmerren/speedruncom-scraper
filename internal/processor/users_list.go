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
	LeaderboardsFile             *repository.ReadRepository
	SupplementaryLeaderboardFile *repository.ReadRepository
	UsersIdListFile              *repository.WriteRepository
	Client                       *srcom_api.SrcomV1Client
}

func (p *UsersListProcessor) Process() error {
	leaderboardsUsers, err := getUsersFromFile(p.LeaderboardsFile)
	if err != nil {
		return err
	}

	for userID := range leaderboardsUsers {
		p.UsersIdListFile.Write([]string{userID})
	}

	supplementaryLeaderboardUsers, err := getUsersFromFile(p.SupplementaryLeaderboardFile)
	if err != nil {
		return err
	}

	for userID := range supplementaryLeaderboardUsers {
		p.UsersIdListFile.Write([]string{userID})
	}

	return nil
}

func getUsersFromFile(file *repository.ReadRepository) (map[string]struct{}, error) {
	allUsers := make(map[string]struct{})
	file.Read()

	for {
		record, err := file.Read()
		if err != nil && !errors.Is(err, io.EOF) {
			return nil, err
		}

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

	return allUsers, nil
}
