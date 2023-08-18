package main

import (
	"bufio"
	"fmt"

	"github.com/alexmerren/speedruncom-scraper/internal/filesystem"
	"github.com/alexmerren/speedruncom-scraper/internal/srcomv1"
	"github.com/alexmerren/speedruncom-scraper/internal/srcomv2"
	"github.com/buger/jsonparser"
)

const (
	allGameIDListV1          = "./data/v1/games-id-list.csv"
	gameOutputFilenameV1     = "./data/v1/games-data.csv"
	categoryOutputFilenameV1 = "./data/v1/categories-data.csv"
	levelOutputFilenameV1    = "./data/v1/level-data.csv"

	allGameIDListV2          = "./data/v2/games-id-list.csv"
	gameOutputFilenameV2     = "./data/v2/games-data.csv"
	categoryOutputFilenameV2 = "./data/v2/categories-data.csv"
	levelOutputFilenameV2    = "./data/v2/level-data.csv"
)

func main() {
	getGameDataV2()
}

//nolint:errcheck// Not worth checking for an error for every file write.
func getGameDataV2() {
	inputFile, err := filesystem.OpenInputFile(allGameIDListV2)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer inputFile.Close()

	gameOuptutFile, err := filesystem.CreateOutputFile(gameOutputFilenameV2)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer gameOuptutFile.Close()
	gameOuptutFile.WriteString("#ID,name,URL,type,rules,releaseDate,addedDate,runCount,playerCount,numCategories,numLevels,emulator\n")

	categoryOuptutFile, err := filesystem.CreateOutputFile(categoryOutputFilenameV2)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer categoryOuptutFile.Close()
	categoryOuptutFile.WriteString("#parentGameID,ID,name,rules,numPlayers\n")

	levelOutputFile, err := filesystem.CreateOutputFile(levelOutputFilenameV2)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer levelOutputFile.Close()
	levelOutputFile.WriteString("#parentGameID,ID,name,rules,numPlayers\n")

	// Scan the input file and get information for each of the game ID's in the
	// input file. We progress to the next line using scanner.Scan()
	scanner := bufio.NewScanner(inputFile)
	scanner.Scan()
	for scanner.Scan() {
		response, err := srcomv2.GetGameData(scanner.Text())
		if err != nil {
			return
		}

		gameID, _, _, _ := jsonparser.Get(response, "game", "id")

		// Step 1. Process each category for a game
		numCategories := 0
		_, err = jsonparser.ArrayEach(response, func(value []byte, dataType jsonparser.ValueType, offset int, err error) {
			numCategories += 1
			categoryID, _, _, _ := jsonparser.Get(value, "id")
			categoryName, _, _, _ := jsonparser.Get(value, "name")
			categoryRules, _, _, _ := jsonparser.Get(value, "rules")
			categoryNumPlayers, _ := jsonparser.GetInt(value, "numPlayers")
			categoryOuptutFile.WriteString(fmt.Sprintf("%s,%s,\"%s\",\"%s\",%d\n", gameID, categoryID, categoryName, categoryRules, categoryNumPlayers))
		}, "categories")
		if err != nil {
			return
		}

		// Step 2. Process each level for a game
		numLevels := 0
		_, err = jsonparser.ArrayEach(response, func(value []byte, dataType jsonparser.ValueType, offset int, err error) {
			numLevels += 1
			levelID, _, _, _ := jsonparser.Get(value, "id")
			levelName, _, _, _ := jsonparser.Get(value, "name")
			levelRules, _, _, _ := jsonparser.Get(value, "rules")
			levelNumPlayers, _ := jsonparser.GetInt(value, "numPlayers")
			levelOutputFile.WriteString(fmt.Sprintf("%s,%s,\"%s\",\"%s\",%d\n", gameID, levelID, levelName, levelRules, levelNumPlayers))
		}, "levels")
		if err != nil {
			return
		}

		// Step 3. Process each game
		gameName, _, _, _ := jsonparser.Get(response, "game", "name")
		gameURL, _, _, _ := jsonparser.Get(response, "game", "url")
		gameType, _, _, _ := jsonparser.Get(response, "game", "type")
		gameEmulator, _ := jsonparser.GetInt(response, "game", "emulator")
		gameReleaseDate, _ := jsonparser.GetInt(response, "game", "releaseDate")
		gameAddedDate, _ := jsonparser.GetInt(response, "game", "addedDate")
		gameRunCount, _ := jsonparser.GetInt(response, "game", "runCount")
		gamePlayerCount, _ := jsonparser.GetInt(response, "game", "totalPlayerCount")
		gameRules, _, _, _ := jsonparser.Get(response, "game", "rules")
		gameOuptutFile.WriteString(fmt.Sprintf("%s,\"%s\",%s,%s,\"%s\",%d,%d,%d,%d,%d,%d,%d\n", gameID, gameName, gameURL, gameType, gameRules, gameReleaseDate, gameAddedDate, gameRunCount, gamePlayerCount, numCategories, numLevels, gameEmulator))
	}

	if err := scanner.Err(); err != nil {
		fmt.Println(err)
		return
	}
}

//nolint:errcheck// Not worth checking for an error for every file write.
func getGameDataV1() {
	inputFile, err := filesystem.OpenInputFile(allGameIDListV1)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer inputFile.Close()

	gameOuptutFile, err := filesystem.CreateOutputFile(gameOutputFilenameV1)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer gameOuptutFile.Close()
	gameOuptutFile.WriteString("#\n")

	categoryOuptutFile, err := filesystem.CreateOutputFile(categoryOutputFilenameV1)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer categoryOuptutFile.Close()
	categoryOuptutFile.WriteString("#\n")

	levelOutputFile, err := filesystem.CreateOutputFile(levelOutputFilenameV1)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer levelOutputFile.Close()
	levelOutputFile.WriteString("#\n")

	// Scan the input file and get information for each of the game ID's in the
	// input file. We progress to the next line using scanner.Scan()
	scanner := bufio.NewScanner(inputFile)
	scanner.Scan()
	for scanner.Scan() {
		_, err := srcomv1.GetGame(scanner.Text())
		if err != nil {
			return
		}

		//gameName, _, _, _ := jsonparser.Get(response, "data", "names", "international")
		//gameURL, _, _, _ := jsonparser.Get(response, "game", "url")
		//gameType, _, _, _ := jsonparser.Get(response, "game", "type")
		//gameEmulator, _ := jsonparser.GetInt(response, "game", "emulator")
		//gameReleaseDate, _ := jsonparser.GetInt(response, "game", "releaseDate")
		//gameAddedDate, _ := jsonparser.GetInt(response, "game", "addedDate")
		//gameRunCount, _ := jsonparser.GetInt(response, "game", "runCount")
		//gamePlayerCount, _ := jsonparser.GetInt(response, "game", "totalPlayerCount")
		//gameRules, _, _, _ := jsonparser.Get(response, "game", "rules")
		//gameOuptutFile.WriteString(fmt.Sprintf("%s,\"%s\",%s,%s,\"%s\",%d,%d,%d,%d,%d,%d,%d\n", gameID, gameName, gameURL, gameType, gameRules, gameReleaseDate, gameAddedDate, gameRunCount, gamePlayerCount, numCategories, numLevels, gameEmulator))

		// Step 1. Process each category for a game
		// Step 2. Process each level for a game
		// Step 3. Process each game
	}

	if err := scanner.Err(); err != nil {
		fmt.Println(err)
		return
	}
}
