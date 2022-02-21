#!/usr/bin/env python3

import argparse
from asyncore import write
import pathlib

EXCLUDED_DIRS = ["generated"]

FILE_HEADER = """// Copyright StackRox Authors
// SPDX-License-Identifier: Apache-2.0

"""


def find_files(path, fileglob):
    files_full = list(path.glob(fileglob))
    return files_full


def contains_header(filepath):
    with open(filepath, mode="r") as file:
        # Check for the content of the first comment line
        if FILE_HEADER in file.read():
            return True
        return False


def write_header(filepath):
    with open(filepath, "r+") as file:
        content = file.read()
        file.seek(0, 0)
        file.write(FILE_HEADER + content)


def main():
    parser = argparse.ArgumentParser()
    parser.add_argument("--path", type=pathlib.Path, help="Path to project basedir")
    arguments = parser.parse_args()

    print(f"Checking all files in {arguments.path} and its child folders")

    proto_files = find_files(arguments.path, "**/?*.proto")
    source_files = find_files(arguments.path, "**/*[!pb|!pb.gw].go")
    
    print(f"Pre-Filter length: {len(proto_files+source_files)}")

    # Filter out any files in excluded directories
    filtered_files = [x for x in (proto_files + source_files) if any(exclusion+"/" not in str(x) for exclusion in EXCLUDED_DIRS) ]

    print(f"Post-Filter length: {len(filtered_files)}")

    for file in filtered_files:
        if not contains_header(file):
            write_header(file)

    print("Please remember to commit any changes!")



if __name__ == "__main__":
    main()
