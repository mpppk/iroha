package cmd

import (
	"fmt"
	"os"

	"github.com/mpppk/iroha/analyze"
	"github.com/mpppk/iroha/storage"

	"github.com/mpppk/iroha/lib"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var analyzeCmd = &cobra.Command{
	Use:   "analyze",
	Short: "Analyze words",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		config := analyze.NewConfigFromViper()
		if config.FilePath == "" {
			FprintlnOrPanic(os.Stderr, "CSV file path should be specified by --file flag")
			os.Exit(1)
		}

		records, err := lib.ReadCSV(config.FilePath)
		panicIfErrorExists(err)

		headers := records[0]
		colIndex, ok := findTargetColIndex(config.ColName, headers)
		if !ok {
			panic(fmt.Errorf("failed to find colName(%s) in headers: %s", config.ColName, headers))
		}

		words := make([]string, 0, len(records)-1)
		for _, record := range records[1:] {
			words = append(words, record[colIndex])
		}

		iroha := lib.NewIroha(words, &storage.None{}, nil)
		iroha.PrintWordCountMap()
		iroha.PrintWordByKatakanaMap()
	},
}

func init() {
	flagKeys := analyze.NewFlagKeys()
	rootCmd.AddCommand(analyzeCmd)
	analyzeCmd.Flags().StringP(flagKeys.File, "f", "", "CSV file path")
	if err := viper.BindPFlag(analyze.ViperPrefix+flagKeys.File, analyzeCmd.Flags().Lookup(flagKeys.File)); err != nil {
		panic(err)
	}
	analyzeCmd.Flags().StringP(flagKeys.ColName, "c", "0", "CSV file path")
	if err := viper.BindPFlag(analyze.ViperPrefix+flagKeys.ColName, analyzeCmd.Flags().Lookup(flagKeys.ColName)); err != nil {
		panic(err)
	}
}
