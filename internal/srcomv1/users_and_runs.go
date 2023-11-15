package srcomv1

import (
	"bufio"
	"encoding/csv"
	"os"

	"github.com/alexmerren/speedruncom-scraper/pkg/srcomv1"
)

func ProcessUsersAndRunsData(
	usersListInputFile,
	usersDataOutputFile,
	usersRunsOutputFile *os.File,
) error {
	usersDataOutputFile.WriteString(usersDataOutputFileHeader)
	usersDataCsvWriter := csv.NewWriter(usersDataOutputFile)
	defer usersDataCsvWriter.Flush()

	usersRunsOutputFile.WriteString(usersRunsOutputFileHeader)
	usersRunsCsvWriter := csv.NewWriter(usersRunsOutputFile)
	defer usersRunsCsvWriter.Flush()

	scanner := bufio.NewScanner(usersListInputFile)
	scanner.Scan()
	for scanner.Scan() {
		userID := scanner.Text()
		userResponse, err := srcomv1.GetUser(userID)
		if err != nil {
			continue
		}

		numRuns, err := processUserRuns(usersRunsCsvWriter, userID)
		if err != nil {
			return err
		}

		err = processUser(usersDataCsvWriter, userID, numRuns, userResponse)
		if err != nil {
			return err
		}
	}

	if err := scanner.Err(); err != nil {
		return err
	}

	return nil
}
