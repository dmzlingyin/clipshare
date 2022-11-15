package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

const (
	version = `         _   _                   _                           
        | | (_)                 | |                          
   ___  | |  _   _ __      ___  | |__     __ _   _ __    ___ 
  / __| | | | | | '_ \    / __| | '_ \   / _' | | '__|  / _ \
 | (__  | | | | | |_) |   \__ \ | | | | | (_| | | |    |  __/
  \___| |_| |_| | .__/    |___/ |_| |_|  \__,_| |_|     \___|
                | |                                          
                |_|   v0.1.0                              `
)

// versionCmd represents the version command
var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "print current version of clipshare",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println(version)
	},
}

func init() {
	rootCmd.AddCommand(versionCmd)
}
