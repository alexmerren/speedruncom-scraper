# Speedrun.com API Scraper

[![Go Report Card](https://goreportcard.com/badge/github.com/alexmerren/speedruncom-scraper)](https://goreportcard.com/report/github.com/alexmerren/speedruncom-scraper)
![Go Version](https://img.shields.io/badge/go%20version-%3E=1.20-61CFDD.svg?style=flat-square)

A series of executables to collect all data available from speedrun.com. The data has been published [here](https://www.kaggle.com/datasets/alexmerren1/speedrun-com-data)!

## 🏃 Executables

// TODO(alex): Write this into a table

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

The compiled binaries can be executed to collect data from the speedrun.com API. The following command collects data for all runs, all leaderboards, all games, and all users (who have contributed to leaderboards) on speedrun.com:

```bash
$ ./dist/games-list && ./dist/games-and-leaderboards-data && ./dist/users-list && ./dist/users-and-runs-data
```

NOTE: For each executable there are repeated API calls. A local HTTP cache has been implemented to remove repeated API calls from the rate-limited API. This cache is saved locally under `httpcache.db`.

## 💭 Feedback and Contribution

Any improvements or requests can be raised via [GitHub Issues](https://github.com/alexmerren/speedruncom-scraper/issues). Any development conversations can be found on on [GitHub Discussions](https://github.com/alexmerren/speedruncom-scraper/discussions).
