package repository

const (
	GamesIdListFilename                  = "./data/v1/games-ids.csv"
	GamesDataFilename                    = "./data/v1/games-data.csv"
	CategoriesDataFilename               = "./data/v1/categories-data.csv"
	LevelsDataFilename                   = "./data/v1/levels-data.csv"
	VariablesDataFilename                = "./data/v1/variables-data.csv"
	ValuesDataFilename                   = "./data/v1/values-data.csv"
	LeaderboardsDataFilename             = "./data/v1/leaderboards-data.csv"
	SupplementaryLeaderboardDataFilename = "./data/v1/supplementary-leaderboard-data.csv"
	RunsDataFilename                     = "./data/v1/runs-data.csv"
	UsersDataFilename                    = "./data/v1/users-data.csv"
	UsersIdListFilename                  = "./data/v1/users-ids.csv"
	DevelopersDataFilename               = "./data/v1/developers-data.csv"
	GenresDataFilename                   = "./data/v1/genres-data.csv"
	PlatformsDataFilename                = "./data/v1/platforms-data.csv"
	PublishersDataFilename               = "./data/v1/publishers-data.csv"
	EnginesDataFilename                  = "./data/v1/engines-data.csv"
	WorldRecordDataFilename              = "./data/v2/world-record-data.csv"
	GamesDataFilenameV2                  = "./data/v2/games-data.csv"
)

var (
	FileColumnDefinitions = map[string][]string{
		GamesIdListFilename:                  {"gameId"},
		UsersIdListFilename:                  {"userId"},
		GamesDataFilename:                    {"gameId", "gameName", "url", "releaseDate", "createdDate", "gameTypes", "platforms", "regions", "genres", "engines", "developers", "publishers", "runsRequireVerification", "runsRequireVideo", "runTimingOptions", "runDefaultTimingOption", "runsEmulatorsAllowed", "isRomhack"},
		GamesDataFilenameV2:                  {"gameId", "gameName", "url", "type", "releaseDate", "addedDate", "runCount", "playerCount", "rules"},
		CategoriesDataFilename:               {"parentGameId", "categoryId", "categoryName", "rules", "type", "numPlayers"},
		LevelsDataFilename:                   {"parentGameId", "levelId", "levelName", "rules"},
		VariablesDataFilename:                {"parentGameId", "variableId", "variableName", "category", "scopeType", "scopeLevel", "isSubcategory", "defaultValue"},
		ValuesDataFilename:                   {"parentGameId", "variableId", "valueId", "label", "rules"},
		LeaderboardsDataFilename:             {"runId", "gameId", "categoryId", "levelId", "place", "date", "primaryTime", "platform", "isEmulated", "players", "examiner", "verifiedDate", "variablesAndValues"},
		SupplementaryLeaderboardDataFilename: {"runId", "gameId", "categoryId", "levelId", "place", "date", "primaryTime", "platform", "isEmulated", "players", "examiner", "verifiedDate", "variablesAndValues"},
		RunsDataFilename:                     {"runId", "gameId", "categoryId", "levelId", "date", "primaryTime", "platform", "isEmulated", "players", "examiner", "verifiedDate", "variablesAndValues", "status", "statusReason"},
		UsersDataFilename:                    {"userId", "username", "signupDate", "location"},
		WorldRecordDataFilename:              {"runId", "gameId", "categoryId", "levelId", "variablesAndValues", "time", "platformId", "isEmulated", "players", "verifier", "runDate", "submittedDate", "verifiedDate", "comment"},
		DevelopersDataFilename:               {"developerId", "name"},
		GenresDataFilename:                   {"genreId", "name"},
		PlatformsDataFilename:                {"platformId", "name", "releaseYear"},
		PublishersDataFilename:               {"publisherId", "name"},
		EnginesDataFilename:                  {"engineId", "name"},
	}
)
