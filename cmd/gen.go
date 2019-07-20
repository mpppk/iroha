package cmd

import (
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"

	"github.com/mpppk/iroha/ktkn"

	"github.com/mpppk/iroha/gen"
	"github.com/mpppk/iroha/storage"

	"github.com/mpppk/iroha/lib"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var genCmd = &cobra.Command{
	Use:   "gen",
	Short: "Generate iroha-uta",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		config := gen.NewConfigFromViper()
		if !config.IsValid() {
			FprintlnOrPanic(os.Stderr, "invalid config:", *config)
			os.Exit(1)
		}

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

		boltStorage, err := storage.NewBolt(config.DBPath)
		panicIfErrorExists(err)
		memoryStorage := storage.NewMemoryWithOtherStorage(boltStorage)

		iroha := lib.NewIroha(words, memoryStorage, config.DepthOptions)
		iroha.PrintWordCountMap()
		iroha.PrintWordByKatakanaMap()
		rowIndicesList, err := iroha.Search()
		panicIfErrorExists(err)

		switch config.OutputMode {
		case gen.PrettyOutputMode:
			gen.PrintAsPretty(records, rowIndicesList)
		case gen.IndicesOutputMode:
			gen.PrintIndices(rowIndicesList)
		}
		FprintfOrPanic(os.Stderr, "%d iroha-uta were found!", len(rowIndicesList))
	},
}

func panicIfErrorExists(err error) {
	if err != nil {
		panic(err)
	}
}

func FprintlnOrPanic(w io.Writer, a ...interface{}) {
	if _, err := fmt.Fprintln(w, a...); err != nil {
		panic(err)
	}
}

func FprintfOrPanic(w io.Writer, format string, a ...interface{}) {
	if _, err := fmt.Fprintf(w, format, a...); err != nil {
		panic(err)
	}
}

func IsValidIroha(words []string) (bool, string) {
	concatenatedWord := strings.Join(words, "")
	n := ktkn.NormalizeKatakanaWord(concatenatedWord)
	runes := []rune(n)

	if len(runes) != int(ktkn.KatakanaLen) {
		return false, n
	}

	if ktkn.HasDuplicatedRune(n) {
		return false, n
	}

	return true, n
}

func findTargetColIndex(colName string, headers []string) (int, bool) {
	if index, err := strconv.Atoi(colName); err == nil {
		return index, true
	}
	for i, header := range headers {
		if header == colName {
			return i, true
		}
	}
	return 0, false
}

func init() {
	flagKeys := gen.NewFlagKeys()
	rootCmd.AddCommand(genCmd)
	genCmd.Flags().StringP(flagKeys.File, "f", "", "CSV file path")
	if err := viper.BindPFlag(flagKeys.File, genCmd.Flags().Lookup(flagKeys.File)); err != nil {
		panic(err)
	}
	genCmd.Flags().StringP(flagKeys.DBPath, "d", "", "DB file path")
	if err := viper.BindPFlag(flagKeys.DBPath, genCmd.Flags().Lookup(flagKeys.DBPath)); err != nil {
		panic(err)
	}
	genCmd.Flags().StringP(flagKeys.ColName, "c", "0", "Target column name or index")
	if err := viper.BindPFlag(flagKeys.ColName, genCmd.Flags().Lookup(flagKeys.ColName)); err != nil {
		panic(err)
	}
	genCmd.Flags().String(flagKeys.OutputMode, "pretty", "Output mode(pretty,indices,none)")
	if err := viper.BindPFlag(flagKeys.OutputMode, genCmd.Flags().Lookup(flagKeys.OutputMode)); err != nil {
		panic(err)
	}
	if err := registerIntToFlags(genCmd, flagKeys.MinParallelDepth, -1, "min depth"); err != nil {
		panic(err)
	}
	if err := registerIntToFlags(genCmd, flagKeys.MaxParallelDepth, -1, "max depth"); err != nil {
		panic(err)
	}
	if err := registerIntToFlags(genCmd, flagKeys.MaxLogDepth, 0, "max log depth"); err != nil {
		panic(err)
	}
	if err := registerIntToFlags(genCmd, flagKeys.MaxStorageDepth, -1, "max storage depth"); err != nil {
		panic(err)
	}
}

func registerIntToFlags(cmd *cobra.Command, name string, value int, usage string) error {
	cmd.Flags().Int(name, value, usage)
	return viper.BindPFlag(name, cmd.Flags().Lookup(name))
}
