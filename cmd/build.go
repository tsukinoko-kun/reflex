package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/tsukinoko-kun/reflex/internal/build"
)

var buildCmd = &cobra.Command{
	Use:   "build",
	Short: "A brief description of your command",

	Run: func(cmd *cobra.Command, args []string) {
		if err := build.Bundle(); err != nil {
			fmt.Println(err)
		}
	},
}

func init() {
	rootCmd.AddCommand(buildCmd)
}
