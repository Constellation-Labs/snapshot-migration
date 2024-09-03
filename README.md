# Snapshot migration tool

## Overview

This script is designed to efficiently organize and move files from a source directory (incremental snapshots used by tessellation 2.x.x) into structured destination directories (incremental snapshots used by tessellation 3.x.x).

The script handles two types of files:

1. Hash-named files (files with hash names, 64 characters long, containing hexadecimal characters)
2. Ordinal-named files (files with numeric names)

## Directory structure

- Hash files: The script organizes hash-named files into a nested directory structure based on the first 6 characters of the file name.
    - *Example*: A file named `04c6b16eded6bf5e6393cbb8b94d9d5b77e6adcf88524bc6f3ca57820d486b8b` will be moved to: `hash/04c/6b1/04c6b16eded6bf5e6393cbb8b94d9d5b77e6adcf88524bc6f3ca57820d486b8b`
- Ordinal files: The script organizes ordinal-named files into directories based on ranges of 20,000.
    - *Example*: A file named `2199620` will be moved to `ordinal/2180000/2199620`.

## Requirements

### Non-nix usage

- Go programming language

### Nix usage

- Nix package manager
- `gomod2nix` for managing Go dependencies in a Nix environment

## Installation

### Non-nix setup

1. Clone the Repository:
```
git clone https://github.com/Constellation-Labs/snapshot-migration.git
cd snapshot-migration
```

2. Build the Script: Compile the Go script into an executable binary:
```
go build
```

### Nix setup

1. Clone the Repository:
```
git clone https://github.com/Constellation-Labs/snapshot-migration.git
cd snapshot-migration
```

2. Generate Nix dependencies:
```
gomod2nix
```

3. Build the project:
```
nix build
```

This will produce a binary in the `result/bin/` directory.

## Usage

### Command-Line flags

- `src`: (Required) The source directory containing incremental snapshots. The script will abort if this flag is not provided.

### Running the script

To execute the script, provide the required `src` flag:

```sh
$ ./snapshot-migration-tool -src ./data/incremental_snapshots
```

## Performance

The script uses concurrency to process files in parallel, taking advantage of multiple CPU cores for faster execution, especially when dealing with large number of files.
