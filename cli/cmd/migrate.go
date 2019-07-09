package cmd

import (
	"github.com/spf13/cobra"
	"go-weixin/service/models"
)

var migrateCmd = &cobra.Command{
	Use: "migrate",
	Short: "migrate",
	Long: `database migrate`,
	Run: func(cmd *cobra.Command, args []string) {
		models.Migrate()
	},
}

func init() {
	rootCmd.AddCommand(migrateCmd)
}

