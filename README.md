# Speedrun.com API Scraper

<p align="center">
  <img src="docs/speedrun_com_logo.png" />
</p>

[![Go Report Card](https://goreportcard.com/badge/github.com/alexmerren/speedruncom-scraper)](https://goreportcard.com/report/github.com/alexmerren/speedruncom-scraper)
![Go Version](https://img.shields.io/badge/go%20version-%3E=1.23-61CFDD.svg?style=flat-square)

A series of executables to collect all data available from [speedrun.com](https://www.speedrun.com). A version of the collected data has been published [here](https://www.kaggle.com/datasets/alexmerren1/speedrun-com-data)!

## ⬇️ Installation

The repository can be installed easily, and binaries can be compiled with the following commands:

```bash
$ git clone git@github.com:alexmerren/speedruncom-scraper.git
...
$ cd speedruncom-scraper
$ make all
...
```

This project requires:

- [Golang 1.23+](https://go.dev/dl/)
- [`gcc` Compatible Compiler](https://gcc.gnu.org)

## 🚀 Usage

The compiled binaries can be executed to scrape data from the [speedrun.com API](https://github.com/speedruncomorg/api). The following command retrieves data for all runs, all leaderboards, all games, and all users (whom have contributed to leaderboards) on [speedrun.com](https://www.speedrun.com):

```bash
./dist/games-list && ./dist/games-data && ./dist/leaderboards-data && ./dist/users-list && ./dist/users-data && ./dist/runs-data
```

Alternatively, there is a Makefile target to run all executables in order:

```bash
make run
```

NOTE: During the scraping process there may be repeated API calls. A local HTTP cache has been implemented to handle repeated API calls locally instead of via the rate-limited API. This cache is saved as `httpcache.db`.

## 🏃 Executables

| Path | Description | Pre-requisite(s) |
| ---- | ----------- | ---------------- |
| `./dist/games-list` | Retrieve all Game IDs and other data (i.e. total number of runs for a game) for verification in other executables. Retrieve other miscellaneous pieces of data such as platforms, developers, genres, etc. | None |
| `./dist/games-data` | Retrieve data on categories, levels, variables, and values, etc. for all game IDs retrieved in `games-list`. | `./dist/games-list` |
| `./dist/leaderboards-data` | Retrieve leaderboard(s) data for all games retrieved in `games-list`. Note: This can fail for games with a high number of runs, use `additional-leaderboards-data` in this case. | `./dist/games-data` |
| `./dist/supplementary-leaderboard-data` | Retrieve leaderboard data for all category/level/variable/value combinations of a game. This executable is tailored to retrieve data for games with an extremely high number of runs i.e. Subway Surfers. This will be extremely inefficient for games with a high count of unique category/level/variable/value combinations. | None |
| `./dist/users-list` | Compile a list of all unique users found on all leaderboards of all games— includes both submitters and verifiers. | `./dist/leaderboards-data` |
| `./dst/users-data` | Retrieve non-PII data for all unique users compiled in `users-list`. | `./dist/users-list` |
| `./dist/runs-data` | Retrieve all runs for all unique users compiled in `users-list`. This **should** be all runs on speedrun.com! | `./dist/users-list` |
| `./dist/world-record-data` | Retrieve world record data for all valid category/level/variable/value combinations of a game. This is experimental, and has a delay of 1s applied to every request to ensure the V2 API is not rate limited externally. | None |

## 📝 Documentation

All documentation can be found in the [docs](./docs/) directory.

## 💭 Feedback and Contribution

Any improvements or requests can be raised via [GitHub Issues](https://github.com/alexmerren/speedruncom-scraper/issues). Any development conversations can be found on on [GitHub Discussions](https://github.com/alexmerren/speedruncom-scraper/discussions).
