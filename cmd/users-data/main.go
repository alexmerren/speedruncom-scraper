package main

import (
	"errors"
	"fmt"
	"io"
	"os"

	"github.com/alexmerren/speedruncom-scraper/internal"
	"github.com/buger/jsonparser"
)

func main() {
	if err := generateUsersDataV1(); err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		os.Exit(1)
	}
}

func generateUsersDataV1() error {
	client := internal.NewSrcomV1Client()

	usersIdListFile, closeFunc, _ := internal.NewCsvReader(internal.UsersIdListFilenameV1)
	defer closeFunc()

	usersDataFile, closeFunc, _ := internal.NewCsvWriter(internal.UsersDataFilenameV1)
	usersDataFile.Write(internal.FileHeaders[internal.UsersDataFilenameV1])
	defer closeFunc()

	usersIdListFile.Read()

	for {
		record, err := usersIdListFile.Read()
		if err != nil && errors.Is(err, io.EOF) {
			break
		}

		userId := record[0]
		response, err := client.GetUser(userId)
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

		// We write 0 as the number of runs as this executable won't deal with runs data.
		usersDataFile.Write([]string{userId, userName, userSignup, userLocation, "0"})
	}

	return nil
}
