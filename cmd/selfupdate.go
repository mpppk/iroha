package cmd

import (
	"log"

	"github.com/mpppk/iroha/lib"
	"github.com/spf13/cobra"
)

var selfUpdateCmd = &cobra.Command{
	Use:   "selfupdate",
	Short: "Update iroha",
	//Long: `Update iroha`,
	Run: func(cmd *cobra.Command, args []string) {
		updated, err := lib.DoSelfUpdate()
		if err != nil {
			log.Println("Binary update failed:", err)
			return
		}
		if updated {
			log.Println("Current binary is the latest version", lib.Version)
		} else {
			log.Println("Successfully updated to version", lib.Version)
		}
	},
}

func init() {
	rootCmd.AddCommand(selfUpdateCmd)
}
