package main

import (
	"bufio"
	"fmt"

	"github.com/alexmerren/speedruncom-scraper/internal/filesystem"
	"github.com/alexmerren/speedruncom-scraper/internal/srcomv1"
)

const (
	allUserIDListV1      = "./data/v1/users-id-list.csv"
	userOutputFilenameV1 = "./data/v1/users-data.csv"
)

func main() {
	response, _ := srcomv1.GetUser("zxznzp0x")
	fmt.Println(string(response))
	// getUsersDataV1()
}

func getUsersDataV1() {
	inputFile, err := filesystem.OpenInputFile(allUserIDListV1)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer inputFile.Close()

	userOutputFile, err := filesystem.CreateOutputFile(userOutputFilenameV1)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer userOutputFile.Close()
	userOutputFile.WriteString("#")

	scanner := bufio.NewScanner(inputFile)
	scanner.Scan()
	for scanner.Scan() {
		fmt.Printf("%s\n", scanner.Text())
		userID := scanner.Text()

		runsResponse, err := srcomv1.GetUserRuns(userID, 0)
		if err != nil {
			fmt.Println(err)
			continue
		}

		fmt.Println(string(runsResponse))

		// loop through the runsResponse(s) until we get information about every run that a user has done.
		// we want the information about the status, when they did the run, what game, category, level, and variables.
		// also get the reviewer for the run. Sum the number of runs so we can include that in a metadata.

		userResponse, err := srcomv1.GetUser(userID)
		if err != nil {
			fmt.Println(err)
			continue
		}

		fmt.Println(string(userResponse))

		// process the user response to get user metadata, personal bests and their place.
		// sum the number of personal bests so we can include that in the metadata.
	}
}
