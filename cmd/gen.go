package cmd

import (
	"fmt"
	"os"
	"strings"

	"github.com/mpppk/iroha/lib"
	"github.com/spf13/cobra"
)

var genCmd = &cobra.Command{
	Use:   "gen",
	Short: "generate iroha",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		words, err := lib.ReadCSV("pokemon_list.csv")
		//words, err := lib.ReadCSV("pokemon_list_medium_mod.csv")
		//words, err := lib.ReadCSV("pokemon_list_small_mod.csv")
		if err != nil {
			panic(err)
		}

		iroha := lib.NewIroha(words)
		irohaWordsList := iroha.Search()
		for _, irohaWords := range irohaWordsList {
			if ok, _ := IsValidIroha(irohaWords); !ok {
				fmt.Fprintln(os.Stderr, "invalid result is returned", irohaWords)
				os.Exit(1)
			}
			fmt.Println(irohaWords)
		}
		fmt.Println(len(irohaWordsList))
	},
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
}
