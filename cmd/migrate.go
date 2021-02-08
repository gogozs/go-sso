package cmd

import (
	"github.com/spf13/cobra"
	"go-sso/internal/repository/storage/mysql"
)

var migrateCmd = &cobra.Command{
	Use:   "migrate",
	Short: "migrate",
	Long:  `database migrate`,
	Run: func(cmd *cobra.Command, args []string) {
		mysql.Migrate()
	},
}

func init() {
	rootCmd.AddCommand(migrateCmd)
}
