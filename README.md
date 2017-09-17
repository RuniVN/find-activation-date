# find-activation-date
Scan the csv file and find real activation date of unique numbers


## Prerequisite

Go installed: https://golang.org/doc/install

# Usage

Run all tests

```
make test
```

Build the project

```
make build
```

Run project

```
make run GOFLAGS="-input=test.csv -workers=8"
```

Flags includes:

- `input` which is location of csv file
- `workers` which is number of workers concurrently runs

The result will be written in `result.csv`
