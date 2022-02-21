#!/usr/bin/env python3

import argparse
from asyncore import write
import pathlib

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

    source_files = find_files(arguments.path, "*/(!generated)/*[!pb|!pb.gw].go")
    proto_files = find_files(arguments.path, "**/?*.proto")

    for file in proto_files+source_files:
        if not contains_header(file):
            write_header(file)

    print("Please remember to commit any changes!")



if __name__ == "__main__":
    main()
