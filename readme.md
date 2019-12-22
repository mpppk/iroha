# iroha
CLI tool for generate Iroha Uta

# Installation

```shell
$ go get github.com/mpppk/iroha
```

# Usage

```shell
$ iroha gen --file word_list.csv
word1,word2, ... ,wordN
.
.
.
M iroha-uta were found!
```

## help
```shell
$ iroha gen --help
Generate iroha-uta

Usage:
  iroha gen [flags]

Flags:
  -c, --col string                    Target column name or index (default "0")
  -d, --db-path string                DB file path
  -f, --file string                   CSV file path
      --gcp-credentials-path string   GCP credentials file path
      --gcp-project-id string         GCP Project ID
  -h, --help                          help for gen
      --max-log-depth int             max log depth
      --max-p-depth int               max depth (default -1)
      --max-s-depth int               max storage depth (default -1)
      --min-p-depth int               min depth (default -1)
      --output-mode string            Output mode(pretty,indices,none) (default "pretty")
      --reset-progress                Ignore PROCESSING progress
      --storage string                Storage type (memory, bolt, gcp) (default "memory")

Global Flags:
      --config string   config file (default is $HOME/.iroha.yaml)
```
