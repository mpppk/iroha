package cmd

import (
	"fmt"

	"github.com/mpppk/iroha/lib"
	"github.com/spf13/cobra"
)

var genCmd = &cobra.Command{
	Use:   "gen",
	Short: "generate iroha",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		words, err := lib.ReadCSV("pokemon_list_medium.csv")
		//words, err := lib.ReadCSV("pokemon_list_small.csv")
		if err != nil {
			panic(err)
		}
		normalizedWords := lib.NormalizeKatakanaWords(words)
		results := lib.CanUseIroha(normalizedWords, lib.NewIroha(), lib.IrohaCache{})
		fmt.Println(results)
	},
}

func init() {
	rootCmd.AddCommand(genCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// genCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// genCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
