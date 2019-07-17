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
	ColName          string
	MinParallelDepth string
	MaxLogDepth      string
	MaxParallelDepth string
	OutputMode       string
}

func NewFlagKeys() *FlagKeys {
	return &FlagKeys{
		File:             "file",
		ColName:          "col",
		MinParallelDepth: "min-p-depth",
		MaxLogDepth:      "max-log-depth",
		MaxParallelDepth: "max-p-depth",
		OutputMode:       "output-mode",
	}
}

type Config struct {
	FilePath     string
	ColName      string
	depthOptions *lib.DepthOptions
	OutputMode   OutputMode
}

func NewConfigFromViper() *Config {
	flagKeys := NewFlagKeys()
	return &Config{
		FilePath: viper.GetString(flagKeys.File),
		ColName:  viper.GetString(flagKeys.ColName),
		depthOptions: &lib.DepthOptions{
			MinParallel: viper.GetInt(flagKeys.MinParallelDepth),
			MaxParallel: viper.GetInt(flagKeys.MaxParallelDepth),
			MaxLog:      viper.GetInt(flagKeys.MaxLogDepth),
		},
		OutputMode: OutputMode(viper.GetString(flagKeys.OutputMode)),
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
