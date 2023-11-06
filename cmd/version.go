/*
Copyright © 2023 Jérémie Veillet <jeremie.veillet@gmail.com>
*/
package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

// versionCmd represents the version command
var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "print Opengraph version",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("opengraph version 1.2.0")
	},
}

func init() {
	rootCmd.AddCommand(versionCmd)
}
