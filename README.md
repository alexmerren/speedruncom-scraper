# Speedrun.com API Scraper

[![Go Report Card](https://goreportcard.com/badge/github.com/alexmerren/speedruncom-scraper)](https://goreportcard.com/report/github.com/alexmerren/speedruncom-scraper)
![Go Version](https://img.shields.io/badge/go%20version-%3E=1.20-61CFDD.svg?style=flat-square)

An easily deployable method of collecting various data from speedrun.com for data science and machine learning applications.

> Let us know what other data that needs to be collected from speedrun.com! Open a [GitHub Issue](https://github.com/alexmerren/speedruncom-scraper/issues) today.

## 🌟 Highlights

 - Want to collect data from speedrun.com? speedruncom-scraper provides an accesible method of data collection.
 - Data from speedrun.com is easily-accessible, and is formatted for applications in data science and machine learning.

## ℹ️  Overview

My final project and dissertation at University of Exeter required data that focused on user behaviour and cumulative culture of online speedrunning communities. This project focuses on reproducing the data used in that study, and publishing tools to recreate that dataset.

Speedruncom-scraper is written in Golang. It can be compiled and deployed to collect data continuously and is formatted to publish after collection.

## 💨 Executables

 1. `dist/games-list`

    * **Reason**: List of all games available via the speedrun.com API. This only collects the internal ID of each game, further information is collected in subsequent functions.
    * **Requirements**: None.
    * **Number of Requests**: Approximately 35,000.

 2. `dist/games-data`

    * **Reason**: Collecting available information for each game using their internal ID. Metadata is collected on the games themselves. Furthermore, the categories, levels, variables, and values are collected and stored.
    * **Requirements**: `dist/games-list`.
    * **Number of Requests**: Approximately 35,000.

 3. `dist/leaderboards-data`

    * **Reason**: Retrieves all leaderboards for every combination of game, category, and level. Each run that conitrbutes to the leaderboards is recorded, along with each player that contributed to the run (amongst other metadata).
    * **Requirements**: `dist/games-list -> dist/games-data`.
    * **Number of Requests**: Approximately 640,000.

 4. `dist/users-list`

    * **Reason**: Creates a list of unique users that appear in the output of the `leaderboards-data` binary.
    * **Requirements**: `dist/games-list -> dist/games-data`.
    * **Number of Requests**: 0.

 5. `dist/users-data`

    * **Reason**: Collect metadata and run data for each user that has contributed to any given leaderboard on speedrun.com.
    * **Requirements**: `dist/games-list -> dist/games-data -> dist/users-list`.
    * **Number of Requests**: Approximately 350,000.

## 🚀 Usage

A set of executables can be compiled using `make build`. These can be executed in a specific order to collect (most) of the data available from speedrun.com.

```bash
$ cd speedruncom-scraper
$ make all
...
```

A complete set of data from speedrun.com can be obtained via the commands:

```bash
$ ./dist/games-list && ./dist/games-data && ./dist/leaderboards-data && ./dist/users-list && ./dist/users-data 
```

NOTE: For each executable (or, each piece of data) there is repeated API calls. A local HTTP cache has been implemented to remove repeated API calls from the rate-limited API. This cache is saved locally under `data/httpcache.db`.

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

## 💭 Feedback and Contributing

If you use this repository- great! If you could let me know any improvements or requests through [GitHub issues](https://github.com/alexmerren/speedruncom-scraper/issues), that would be great.

Furthermore, if you want to join discussions on the developement of `speedruncom-scraper`, find the conversations on [GitHub discussions](https://github.com/alexmerren/speedruncom-scraper/discussions).
