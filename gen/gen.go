package gen

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/mpppk/iroha/storage"

	"github.com/mpppk/iroha/lib"

	"github.com/spf13/viper"
)

type OutputMode string

var PrettyOutputMode OutputMode = "pretty"
var IndicesOutputMode OutputMode = "indices"
var NoneOutputMode OutputMode = "none"

type FlagKeys struct {
	File               string
	DBPath             string
	ColName            string
	MinParallelDepth   string
	MaxLogDepth        string
	MaxParallelDepth   string
	MaxStorageDepth    string
	Storage            string
	OutputMode         string
	ResetProgress      string
	GCPProjectId       string
	GCPCredentialsPath string
}

func NewFlagKeys() *FlagKeys {
	return &FlagKeys{
		File:               "file",
		DBPath:             "db-path",
		ColName:            "col",
		MinParallelDepth:   "min-p-depth",
		MaxLogDepth:        "max-log-depth",
		MaxParallelDepth:   "max-p-depth",
		MaxStorageDepth:    "max-s-depth",
		Storage:            "storage",
		OutputMode:         "output-mode",
		ResetProgress:      "reset-progress",
		GCPProjectId:       "gcp-project-id",
		GCPCredentialsPath: "gcp-credentials-path",
	}
}

type GCPConfig struct {
	ProjectId       string
	CredentialsPath string
}

type Config struct {
	FilePath      string
	DBPath        string
	ColName       string
	DepthOptions  *lib.DepthOptions
	OutputMode    OutputMode
	ResetProgress bool
	Storage       storage.Type
	GCP           *GCPConfig
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
		Storage:       storage.Type(viper.GetString(flagKeys.Storage)),
		ResetProgress: viper.GetBool(flagKeys.ResetProgress),
		GCP: &GCPConfig{
			ProjectId:       viper.GetString(flagKeys.GCPProjectId),
			CredentialsPath: viper.GetString(flagKeys.GCPCredentialsPath),
		},
	}
}

func (c *Config) IsValid() error {
	if !isValidOutputMode(c.OutputMode) {
		return fmt.Errorf("invalid output mode was given. provided string: %s", c.OutputMode)
	}
	if !isValidStorageType(c.Storage) {
		return fmt.Errorf("invalid storage type was given. provided string: %s", c.Storage)
	}
	return nil
}

func (c *Config) IsValidGCPConfig() error {
	// if storage type is GCP, and GCP config file or project is does not specified
	if c.Storage != storage.GCPType {
		// FIXME: print warning if gcp config is provided.
		return nil
	}
	if c.GCP.ProjectId == "" {
		return fmt.Errorf("even though storage type is set to 'GCP', GCP project ID is not provided")
	}
	if c.GCP.CredentialsPath == "" {
		return fmt.Errorf("even though storage type is set to 'GCP', GCP credentials path is not provided")
	}
	return nil
}

func isValidOutputMode(mode OutputMode) bool {
	return mode == PrettyOutputMode ||
		mode == IndicesOutputMode ||
		mode == NoneOutputMode
}

func isValidStorageType(storageType storage.Type) bool {
	for _, st := range storage.GetStorageTypes() {
		if st == storageType {
			return true
		}
	}
	return false
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
