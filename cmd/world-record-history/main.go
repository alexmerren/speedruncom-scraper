package main

import (
	"bufio"
	"fmt"

	"github.com/alexmerren/speedruncom-scraper/internal/filesystem"
)

const (
	allGameIDListV2           = "./data/v2/games-id-list.csv"
	worldRecordOutputFilename = "./data/v2/world-record-history-data.csv"
)

func main() {
	getWorldRecordHistory()
	fmt.Println("Not implemented yet...")
}

func getWorldRecordHistory() {
	inputFile, err := filesystem.OpenInputFile(allGameIDListV2)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer inputFile.Close()

	wrOutputFile, err := filesystem.CreateOutputFile(worldRecordOutputFilename)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer wrOutputFile.Close()
	wrOutputFile.WriteString("#\n")

	scanner := bufio.NewScanner(inputFile)
	scanner.Scan()
	for scanner.Scan() {
		// response, err := srcomv2.GetGameCategoryWorldRecordHistory(scanner.Text())
		// if err != nil {
		// 	return
		// }
		fmt.Println(scanner.Text())
	}
}
