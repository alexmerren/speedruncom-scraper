package internal

const (
	GamesIdListFilenameV1      = "./data/v1/games-id-list.csv"
	GamesDataFilenameV1        = "./data/v1/games-data.csv"
	GamesDataFilenameV2        = "./data/v2/games-data.csv"
	CategoriesDataFilenameV1   = "./data/v1/categories-data.csv"
	LevelsDataFilenameV1       = "./data/v1/levels-data.csv"
	VariablesDataFilenameV1    = "./data/v1/variables-data.csv"
	ValuesDataFilenameV1       = "./data/v1/values-data.csv"
	LeaderboardsDataFilenameV1 = "./data/v1/leaderboards-data.csv"
	RunsDataFilenameV1         = "./data/v1/runs-data.csv"
	UsersDataFilenameV1        = "./data/v1/users-data.csv"
	UsersIdListFilenameV1      = "./data/v1/users-id-list.csv"
	WorldRecordDataFilenameV2  = "./data/v2/world-record-data.csv"
)

var (
	FileHeaders = map[string][]string{
		GamesIdListFilenameV1:      {"gameId"},
		UsersIdListFilenameV1:      {"userId"},
		GamesDataFilenameV1:        {"gameId", "gameName", "url", "releaseDate", "createdDate", "numCategories", "numLevels"},
		GamesDataFilenameV2:        {"gameId", "gameName", "url", "type", "releaseDate", "addedDate", "runCount", "playerCount", "rules"},
		CategoriesDataFilenameV1:   {"parentGameId", "categoryId", "categoryName", "rules", "type", "numPlayers"},
		LevelsDataFilenameV1:       {"parentGameId", "levelId", "levelName", "rules"},
		VariablesDataFilenameV1:    {"parentGameId", "variableId", "variableName", "category", "scope", "isSubcategory", "defaultValue"},
		ValuesDataFilenameV1:       {"parentGameId", "variableId", "valueId", "label", "rules"},
		LeaderboardsDataFilenameV1: {"runId", "gameId", "categoryId", "levelId", "place", "date", "primaryTime", "platform", "isEmulated", "players", "examiner", "verifiedDate", "variablesAndValues"},
		RunsDataFilenameV1:         {"runId", "gameId", "categoryId", "levelId", "date", "primaryTime", "platform", "isEmulated", "players", "examiner", "verifiedDate", "variablesAndValues", "status", "statusReason"},
		UsersDataFilenameV1:        {"userId", "username", "signupDate", "location", "numRuns"},
		WorldRecordDataFilenameV2:  {}, // TODO Implement world record data
	}
)
