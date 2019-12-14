package analyze

import (
	"github.com/spf13/viper"
)

type FlagKeys struct {
	File    string
	ColName string
}

var ViperPrefix = "analyze-"

func NewFlagKeys() *FlagKeys {
	return &FlagKeys{
		File:    "file",
		ColName: "col",
	}
}

type Config struct {
	FilePath string
	ColName  string
}

func NewConfigFromViper() *Config {
	flagKeys := NewFlagKeys()
	return &Config{
		FilePath: viper.GetString(ViperPrefix + flagKeys.File),
		ColName:  viper.GetString(ViperPrefix + flagKeys.ColName),
	}
}
