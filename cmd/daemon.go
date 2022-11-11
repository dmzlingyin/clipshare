package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

// daemonCmd represents the daemon command
var daemonCmd = &cobra.Command{
	Use:   "daemon",
	Short: "make the clipboard as a daemon",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("daemon start...")
	},
}

func init() {
	rootCmd.AddCommand(daemonCmd)
}
