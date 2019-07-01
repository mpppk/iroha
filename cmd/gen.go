package cmd

import (
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/mpppk/iroha/lib"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var fileFlagKey = "file"

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
		words, err := lib.ReadCSV(filePath)
		if err != nil {
			panic(err)
		}

		iroha := lib.NewIroha(words)
		//iroha.PrintWordCountMap()
		//iroha.PrintWordByKatakanaMap()
		irohaWordsList := iroha.Search()
		for _, irohaWords := range irohaWordsList {
			if ok, _ := IsValidIroha(irohaWords); !ok {
				FprintlnOrPanic(os.Stderr, "invalid result is returned", irohaWords)
				os.Exit(1)
			}
			fmt.Println(strings.Join(irohaWords, ","))
		}
		FprintfOrPanic(os.Stderr, "%d iroha-uta were found!", len(irohaWordsList))
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

func init() {
	rootCmd.AddCommand(genCmd)
	genCmd.Flags().StringP(fileFlagKey, "f", "", "CSV file path")
	if err := viper.BindPFlag(fileFlagKey, genCmd.Flags().Lookup(fileFlagKey)); err != nil {
		panic(err)
	}
}
