## Number: 02
## Date: 2024-12-24
## Title: Spot Collection for Leaderboard Data

Collecting leaderboards data for games that have a high number of runs 
(i.e. [Subway Surfers](www.speedrun.com/subsurf)) leads to issues with pagination 
becoming too high. After 10,000 runs the pagination is too high and the requests 
fail. The answer to this is to collect leaderboard data for every variable and 
value combination. These are represented as additional filters when viewing a category 
on [speedrun.com](www.speedrun.com). 

The logic behind choosing appropriate variables and values is a little complicated, 
so below is a POC that shows how the [`additional-leaderboards-data`](../cmd/additional-leaderboards-data/main.go) 
executable decides to collect this data.

This executable is an exception, as it accepts a parameter of a single game ID.
In the below example, all current leaderboard runs of game `y65797de` will be persisted 
to [`additional-leaderboards-data.csv`](../data/v1/additional-leaderboards-data.csv)

```bash
./dist/additional-leaderboards-data y65797de
```

### Spot Collection POC

```python
import requests
import json
import itertools

GAME_ID = "y65797de"


def main():
    response = requests.get(
        f"https://www.speedrun.com/api/v1/games/{GAME_ID}?embed=levels,categories,variables"
    )
    data = response.json().get("data")
    combinations = generate_combinations(data)
    for combination in combinations:
        print(combination)


def generate_combinations(data: dict[str, any]):
    game_id = data["id"]
    categories = data["categories"]["data"]
    variables = data["variables"]["data"]
    levels = data["levels"]["data"]

    combinations = []

    for category in categories:
        category_id = category["id"]

        if category["type"] == "per-game":
            applicable_variables = {
                variable["id"]: list(variable["values"]["values"].keys())
                for variable in variables
                if variable_is_valid_for_category(variable, category_id)
            }

            for element in itertools.product(*applicable_variables.values()):
                combinations.append(
                    Combination(
                        game_id,
                        category_id,
                        None,
                        list(applicable_variables.keys()),
                        list(element),
                    )
                )

        if category["type"] == "per-level":
            for level in levels:
                level_id = level["id"]
                applicable_variables = {
                    variable["id"]: list(variable["values"]["values"].keys())
                    for variable in variables
                    if variable_is_valid_for_category_and_level(
                        variable, category_id, level_id
                    )
                }

                for element in itertools.product(*applicable_variables.values()):
                    combinations.append(
                        Combination(
                            game_id,
                            category_id,
                            level_id,
                            list(applicable_variables.keys()),
                            list(element),
                        )
                    )

    return combinations


class Combination:
    def __init__(
        self,
        game_id: str,
        category_id: str,
        level_id: str | None,
        variables: list[str],
        values: list[str],
    ):
        self.game_id = game_id
        self.category_id = category_id
        self.level_id = level_id
        self.variables = variables
        self.values = values

    def __repr__(self):
        return f'Combination(game_id="{self.game_id}",category_id="{self.category_id}",level_id="{self.level_id}",variables="{self.variables}",values="{self.values}")'

    def __str__(self):
        return f"game_id: {self.game_id}, category_id: {self.category_id}, level_id: {self.level_id}, variables: {self.variables}, values: {self.values}"


def variable_is_valid_for_category(variable, category_id):
    if not variable["is-subcategory"]:
        return False

    if variable["scope"]["type"] not in ("global", "full-game"):
        return False

    if variable["category"] != None and variable["category"] != category_id:
        return False

    return True


def variable_is_valid_for_category_and_level(variable, category_id, level_id):
    if not variable["is-subcategory"]:
        return False

    if variable["scope"]["type"] not in ("global", "all-levels", "single-level"):
        return False

    if (
        variable["scope"]["type"] == "single-level"
        and variable["scope"]["level"] != level_id
    ):
        return False

    if variable["category"] != None and variable["category"] != category_id:
        return False

    return True
```