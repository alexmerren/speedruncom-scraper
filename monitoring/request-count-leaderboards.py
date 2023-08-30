#!/usr/bin/env python3
import pandas as pd
from datetime import timedelta

CATEGORIES_FILE = "../data/v1/categories-data.csv"
GAMES_FILE = "../data/v1/games-data.csv"

def main():
    categories_df = pd.read_csv(CATEGORIES_FILE)
    c_df = categories_df[['#parentGameID', 'type']].groupby(['#parentGameID', 'type']).size()
    c_df = c_df.reset_index().rename(columns={0: 'count', '#parentGameID': '#ID'})
    levels_df = c_df[c_df['type'] == 'per-level']
    categories_df = c_df[c_df['type'] == 'per-game']

    games_df = pd.read_csv(GAMES_FILE)
    games_df = games_df[['#ID', 'numLevels']]
    games_df.loc[games_df['numLevels'] == 0, 'numLevels'] = 1
    games_df = games_df.merge(levels_df, how='inner', on='#ID')[['#ID', 'numLevels', 'count']].rename(columns={'count': 'numLevelCategories'})
    games_df = games_df.merge(categories_df, how='inner', on='#ID')[['#ID', 'numLevels', 'numLevelCategories', 'count']].rename(columns={'count': 'numGameCategories'})
    games_df['total'] = games_df['numLevelCategories'] * games_df['numLevels'] + games_df['numGameCategories']

    print(f"Number of leaderboards-data requests: {games_df.sum()['total']}")
    print(f"Time of all leaderboards-data requests: {str(timedelta(seconds=games_df.sum()['total']*0.7))}")

if __name__ == "__main__":
    main()
