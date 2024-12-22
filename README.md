# Speedrun.com API Scraper

<p align="center">
  <img src="docs/speedrun_com_logo.png" />
</p>

[![Go Report Card](https://goreportcard.com/badge/github.com/alexmerren/speedruncom-scraper)](https://goreportcard.com/report/github.com/alexmerren/speedruncom-scraper)
![Go Version](https://img.shields.io/badge/go%20version-%3E=1.20-61CFDD.svg?style=flat-square)

A series of executables to collect all data available from [speedrun.com](https://www.speedrun.com). A version of the collected data has been published [here](https://www.kaggle.com/datasets/alexmerren1/speedrun-com-data)!

## ⬇️  Installation


The repository can be installed easily, and binaries can be compiled with the following commands:

```bash
$ git clone git@github.com:alexmerren/speedruncom-scraper.git
...
$ cd speedruncom-scraper
$ make all
...
```

This project requires:

 * [Golang 1.20+](https://go.dev/dl/)
 * [`gcc` Compatible Compiler](https://gcc.gnu.org)

## 🚀 Usage

The compiled binaries can be executed to collect data from the [speedrun.com API](https://github.com/speedruncomorg/api). The following command collects data for all runs, all leaderboards, all games, and all users (who have contributed to leaderboards) on [speedrun.com](https://www.speedrun.com):

```bash
./dist/games-list && ./dist/games-data && ./dist/leaderboards-data && ./dist/users-list && ./dist/users-data && ./dist/runs-data
```

NOTE: For each executable there are repeated API calls. A local HTTP cache has been implemented to remove repeated API calls from the rate-limited API. This cache is saved locally under `httpcache.db`.

## 🏃 Executables

| Path                       | Description | Pre-requisite(s)           |
|----------------------------|-------------|----------------------------|
| `./dist/games-list`        | | None                       |
| `./dist/games-data`        | | `./dist/games-list`        |
| `./dist/leaderboards-data` | | `./dist/games-data`        |
| `./dist/users-list`        | | `./dist/leaderboards-data` |
| `./dst/users-data`         | | `./dist/users-list`        |
| `./dist/runs-data`         | | `./dist/users-list`        |

## 💭 Feedback and Contribution

Any improvements or requests can be raised via [GitHub Issues](https://github.com/alexmerren/speedruncom-scraper/issues). Any development conversations can be found on on [GitHub Discussions](https://github.com/alexmerren/speedruncom-scraper/discussions).
