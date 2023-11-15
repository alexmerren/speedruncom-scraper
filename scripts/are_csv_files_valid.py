#!/usr/bin/env python3

import pandas as pd
from glob import glob

def main():
    for file in glob('./data/**/*.csv', recursive=True):
        print(file)
        pd.read_csv(file)

if __name__ == "__main__":
    main()
