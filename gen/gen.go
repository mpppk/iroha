package gen

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/mpppk/iroha/lib"

	"github.com/spf13/viper"
)

var PrettyOutputMode OutputMode = "pretty"
var IndicesOutputMode OutputMode = "indices"
var NoneOutputMode OutputMode = "none"

type OutputMode string
type FlagKeys struct {
	File             string
	DBPath           string
	ColName          string
	MinParallelDepth string
	MaxLogDepth      string
	MaxParallelDepth string
	MaxStorageDepth  string
	OutputMode       string
	ResetProgress    string
}

func NewFlagKeys() *FlagKeys {
	return &FlagKeys{
		File:             "file",
		DBPath:           "db-path",
		ColName:          "col",
		MinParallelDepth: "min-p-depth",
		MaxLogDepth:      "max-log-depth",
		MaxParallelDepth: "max-p-depth",
		MaxStorageDepth:  "max-s-depth",
		OutputMode:       "output-mode",
		ResetProgress:    "reset-progress",
	}
}

type Config struct {
	FilePath      string
	DBPath        string
	ColName       string
	DepthOptions  *lib.DepthOptions
	OutputMode    OutputMode
	ResetProgress bool
}

func NewConfigFromViper() *Config {
	flagKeys := NewFlagKeys()
	return &Config{
		FilePath: viper.GetString(flagKeys.File),
		DBPath:   viper.GetString(flagKeys.DBPath),
		ColName:  viper.GetString(flagKeys.ColName),
		DepthOptions: &lib.DepthOptions{
			MinParallel: viper.GetInt(flagKeys.MinParallelDepth),
			MaxParallel: viper.GetInt(flagKeys.MaxParallelDepth),
			MaxLog:      viper.GetInt(flagKeys.MaxLogDepth),
			MaxStorage:  viper.GetInt(flagKeys.MaxStorageDepth),
		},
		OutputMode:    OutputMode(viper.GetString(flagKeys.OutputMode)),
		ResetProgress: viper.GetBool(flagKeys.ResetProgress),
	}
}

func (c *Config) IsValid() bool {
	return isValidOutputMode(c.OutputMode)
}

func isValidOutputMode(mode OutputMode) bool {
	return mode == PrettyOutputMode ||
		mode == IndicesOutputMode ||
		mode == NoneOutputMode
}

func PrintAsPretty(records [][]string, rowIndicesList [][]int) {
	for _, rowIndices := range rowIndicesList {
		for _, rowIndex := range rowIndices {
			// 行番号表示には、header+行番号が1始まりであることを考慮して2を足す
			// recordsにはheaderを考慮して1を足す
			fmt.Println(rowIndex+2, strings.Join(records[rowIndex+1], ","))
		}
		fmt.Print("\n----------\n\n")
	}
}

func PrintIndices(rowIndicesList [][]int) {
	for _, rowIndices := range rowIndicesList {
		fmt.Println(strings.Join(intSliceToStringSlice(rowIndices), ","))
	}
}

func intSliceToStringSlice(slice []int) (strSlice []string) {
	for _, v := range slice {
		strSlice = append(strSlice, strconv.Itoa(v))
	}
	return
}
