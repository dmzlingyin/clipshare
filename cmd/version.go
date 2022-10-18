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
                |_|   v0.1                              `
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

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// versionCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// versionCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
