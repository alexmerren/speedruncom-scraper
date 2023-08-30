#!/usr/bin/env python3
import pandas as pd
import psutil as ps

from datetime import datetime
import time

LEADERBOARD_FILE = "../data/v1/leaderboards-data.csv"
PROCESS_NAME = "leaderboards-data"

LEFT_COLUMN_PAD = 15 
RIGHT_COLUMN_PAD = 10

def main():
    df = pd.read_csv(LEADERBOARD_FILE)
    num_games = len(df['gameID'].unique())

    process = None
    for proc in ps.process_iter():
        if PROCESS_NAME in proc.name():
            process = proc
            break

    elapsed_time = time.time() - process.create_time()    
    amount_left = 35000 - num_games

    print(f"{'numRuns: ':<{LEFT_COLUMN_PAD}}{len(df.index):>{RIGHT_COLUMN_PAD}}")
    print(f"{'numGames: ':<{LEFT_COLUMN_PAD}}{num_games:>{RIGHT_COLUMN_PAD}}")
    print(f"{'elapsedTime: ':<{LEFT_COLUMN_PAD}}{elapsed_time:>{RIGHT_COLUMN_PAD}.2f}s")

if __name__ == "__main__":
    main()
