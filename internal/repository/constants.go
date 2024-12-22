package repository

const (
	GamesIdListFilename      = "./data/v1/games-id-list.csv"
	GamesDataFilename        = "./data/v1/games-data.csv"
	CategoriesDataFilename   = "./data/v1/categories-data.csv"
	LevelsDataFilename       = "./data/v1/levels-data.csv"
	VariablesDataFilename    = "./data/v1/variables-data.csv"
	ValuesDataFilename       = "./data/v1/values-data.csv"
	LeaderboardsDataFilename = "./data/v1/leaderboards-data.csv"
	RunsDataFilename         = "./data/v1/runs-data.csv"
	UsersDataFilename        = "./data/v1/users-data.csv"
	UsersIdListFilename      = "./data/v1/users-id-list.csv"
)

var (
	FileColumnDefinitions = map[string][]string{
		GamesIdListFilename:      {"gameId"},
		UsersIdListFilename:      {"userId"},
		GamesDataFilename:        {"gameId", "gameName", "url", "releaseDate", "createdDate"},
		CategoriesDataFilename:   {"parentGameId", "categoryId", "categoryName", "rules", "type", "numPlayers"},
		LevelsDataFilename:       {"parentGameId", "levelId", "levelName", "rules"},
		VariablesDataFilename:    {"parentGameId", "variableId", "variableName", "category", "scope", "isSubcategory", "defaultValue"},
		ValuesDataFilename:       {"parentGameId", "variableId", "valueId", "label", "rules"},
		LeaderboardsDataFilename: {"runId", "gameId", "categoryId", "levelId", "place", "date", "primaryTime", "platform", "isEmulated", "players", "examiner", "verifiedDate", "variablesAndValues"},
		RunsDataFilename:         {"runId", "gameId", "categoryId", "levelId", "date", "primaryTime", "platform", "isEmulated", "players", "examiner", "verifiedDate", "variablesAndValues", "status", "statusReason"},
		UsersDataFilename:        {"userId", "username", "signupDate", "location"},
	}

	FileComments = map[string]string{
		GamesIdListFilename:      "#games-id-list.csv",
		UsersIdListFilename:      "#users-id-list.csv",
		GamesDataFilename:        "#games-data.csv",
		CategoriesDataFilename:   "#categories-data.csv",
		LevelsDataFilename:       "#levels-data.csv",
		VariablesDataFilename:    "#variables-data.csv",
		ValuesDataFilename:       "#values-data.csv",
		LeaderboardsDataFilename: "#leaderboads-data.csv",
		RunsDataFilename:         "#runs-data.csv",
		UsersDataFilename:        "#users-data.csv",
	}
)
