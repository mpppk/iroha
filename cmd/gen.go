package cmd

import (
	"fmt"
	"io"
	"log"
	"os"
	"strconv"
	"strings"

	"github.com/mpppk/iroha/lib"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var fileFlagKey = "file"
var colNameKey = "col"

var genCmd = &cobra.Command{
	Use:   "gen",
	Short: "Generate iroha-uta",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		filePath := viper.GetString(fileFlagKey)
		if filePath == "" {
			FprintlnOrPanic(os.Stderr, "CSV file path should be specified by --file flag")
			os.Exit(1)
		}
		colName := viper.GetString(colNameKey)

		records, err := lib.ReadCSV(filePath)
		if err != nil {
			panic(err)
		}
		headers := records[0]
		colIndex, ok := findTargetColIndex(colName, headers)
		if !ok {
			panic(fmt.Errorf("failed to find colName(%s) in headers: %s", colName, headers))
		}

		log.Print("start preprocessing...")
		words := make([]string, 0, len(records)-1)
		for _, record := range records[1:] {
			words = append(words, record[colIndex])
		}

		iroha := lib.NewIroha(words)
		log.Println("preprocessing is done")
		iroha.PrintWordCountMap()
		iroha.PrintWordByKatakanaMap()
		log.Print("start searching...")
		rowIndicesList, err := iroha.Search()
		if err != nil {
			panic(err)
		}

		for _, rowIndices := range rowIndicesList {
			for _, rowIndex := range rowIndices {
				// 行番号表示には、header+行番号が1始まりであることを考慮して2を足す
				// recordsにはheaderを考慮して1を足す
				fmt.Println(rowIndex+2, strings.Join(records[rowIndex+1], ","))
			}
			fmt.Print("\n----------\n\n")
		}
		FprintfOrPanic(os.Stderr, "%d iroha-uta were found!", len(rowIndicesList))
	},
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
	n := lib.NormalizeKatakanaWord(concatenatedWord)
	runes := []rune(n)

	if len(runes) != int(lib.KatakanaLen) {
		return false, n
	}

	if lib.HasDuplicatedRune(n) {
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
	rootCmd.AddCommand(genCmd)
	genCmd.Flags().StringP(fileFlagKey, "f", "", "CSV file path")
	if err := viper.BindPFlag(fileFlagKey, genCmd.Flags().Lookup(fileFlagKey)); err != nil {
		panic(err)
	}
	genCmd.Flags().StringP(colNameKey, "c", "0", "Target column name or index")
	if err := viper.BindPFlag(colNameKey, genCmd.Flags().Lookup(colNameKey)); err != nil {
		panic(err)
	}
}
