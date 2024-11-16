# Bee Solver

Bee Solver is a command-line tool and Text User Interface (TUI) application designed to help solve the New York Times Spelling Bee puzzle. It takes a set of letters (one required center letter and six outer letters) and generates all possible valid words that can be made using these letters.

## Features

- Command-line interface for quick solving
- Interactive TUI for an enhanced user experience
- Validates words against Spelling Bee rules:
  - Words must be 4+ letters long
  - Words must include the center letter
  - Letters can be reused
  - Only valid dictionary words are accepted

## Installation

You'll need `go` 1.23 or later and run 
```sh
go install
```

## Usage

##### Solve a bee puzzle and ouput solution:

```sh
./bee solve -l "abcdefg"
```

##### Launch TUI for interactive exploration of solutions

```sh
bee
```

## Configuration

The application can be configured using the following flags

| flag               | description                                         |
| ------------------ | --------------------------------------------------- |
| `-w`, `--wordlist` | Path to the word list file (default: `data/en.txt`) |
| `--logfile`        | Path to the log file (default: `bee.log`)           |
| `-l`, `--letters`  | Letters to use, starting with the center letter     |
