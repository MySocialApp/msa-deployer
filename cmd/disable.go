package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

// disableCmd represents the disable command
var disableCmd = &cobra.Command{
	Use:   "disable",
	Short: "Disable Ingress client (not implemented)",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("disable called")
	},
}

func init() {
	rootCmd.AddCommand(disableCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// disableCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// disableCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
