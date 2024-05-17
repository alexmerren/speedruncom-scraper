# Speedrun.com API Scraper

[![Go Report Card](https://goreportcard.com/badge/github.com/alexmerren/speedruncom-scraper)](https://goreportcard.com/report/github.com/alexmerren/speedruncom-scraper)
![Go Version](https://img.shields.io/badge/go%20version-%3E=1.20-61CFDD.svg?style=flat-square)

A system to collect data from speedrun.com for machine learning and data science applications.

> Let us know what other data that needs to be collected from speedrun.com! Open a [GitHub Issue](https://github.com/alexmerren/speedruncom-scraper/issues) today.

Find the published dataset [here](https://www.kaggle.com/datasets/alexmerren1/speedrun-com-data)!

## 🌟 Highlights

 - Want to collect data from speedrun.com? speedruncom-scraper provides an accesible method of data collection.
 - Data from speedrun.com is easily-accessible, and is formatted for applications in data science and machine learning.

## ℹ️  Overview

My final project and dissertation at University of Exeter required data that focused on user behaviour and cumulative culture of online speedrunning communities. This project focuses on reproducing the data used in that study, and publishing tools to recreate that dataset.

Speedruncom-scraper is written in Golang. It can be compiled and deployed to collect data continuously and is formatted to publish after collection.

## 💨 Executables

 * [`games-list`](./cmd/games-list/main.go)

 * [`games-data`](./cmd/games-data/main.go)

 * [`leaderboards-data`](./cmd/leaderboards-data/main.go)

 * [`games-and-leaderboards-data`](./cmd/games-and-leaderboards-data/main.go)

 * [`users-list`](./cmd/users-list/main.go)

 * [`users-and-runs-data`](./cmd/users-and-runs-data/main.go)

 * [`world-record-history-data`](./cmd/world-record-history-data/main.go) (WIP)

## 🚀 Usage

A set of executables can be compiled using `make build`. These can be executed in a specific order to collect (most) of the data available from speedrun.com.

```bash
$ cd speedruncom-scraper
$ make all
...
```

A complete set of data from speedrun.com can be obtained via the commands:

```bash
$ ./dist/games-list && ./dist/games-and-leaderboards-data && ./dist/users-list && ./dist/users-and-runs-data
```

NOTE: For each executable (or, each piece of data) there is repeated API calls. A local HTTP cache has been implemented to remove repeated API calls from the rate-limited API. This cache is saved locally under `httpcache.db`.

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

## Previous Collections

The last full data collection occurred in November 2023, here are the number of lines generated for each executable:

```
  390286 ./data/v1/users-data.csv
  391092 ./data/v1/users-id-list.csv
 1847581 ./data/v1/leaderboards-data.csv
 3995740 ./data/v1/runs-data.csv
   37249 ./data/v1/games-data.csv
   56721 ./data/v1/variables-data.csv
  252642 ./data/v1/values-data.csv
  266380 ./data/v1/levels-data.csv
   37251 ./data/v1/games-id-list.csv
  147959 ./data/v1/categories-data.csv
   37870 ./data/v2/games-id-list.csv
```

## 💭 Feedback and Contributing

If you use this repository- great! If you could let me know any improvements or requests through [GitHub issues](https://github.com/alexmerren/speedruncom-scraper/issues), that would be great.

Furthermore, if you want to join discussions on the developement of `speedruncom-scraper`, find the conversations on [GitHub discussions](https://github.com/alexmerren/speedruncom-scraper/discussions).
