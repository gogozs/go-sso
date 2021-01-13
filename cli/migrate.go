package cli

import (
	"github.com/spf13/cobra"
	"go-sso/storage/mysql/model"
)

var migrateCmd = &cobra.Command{
	Use:   "migrate",
	Short: "migrate",
	Long:  `database migrate`,
	Run: func(cmd *cobra.Command, args []string) {
		model.Migrate()
	},
}

func init() {
	rootCmd.AddCommand(migrateCmd)
}
