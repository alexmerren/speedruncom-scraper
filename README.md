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

## 🚀 Usage

A set of executables can be compiled using `make build`. These can be executed in a specific order to collect (most) of the data available from speedrun.com.

```bash
$ cd speedruncom-scraper
$ make all
...
```

Since an order is required for each of the executables, please find an order required for each executable below:

1. *games-list*: No dependencies.
2. *games-data*: `games-list`.
3. *world-records-data*: `games-list`.
3. *leaderboards-data*: `games-list -> games-data`.
4. *users-list*: `games-list -> games-data`.
5. *users-data*: `games-list -> games-data -> users-list`.
6. *runs-data*: `games-list -> games-data -> users-list`.

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
