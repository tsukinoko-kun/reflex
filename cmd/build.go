package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/tsukinoko-kun/reflex/internal/build"
)

var buildCmd = &cobra.Command{
	Use:   "build",
	Short: "A brief description of your command",

	Run: func(cmd *cobra.Command, args []string) {
		if err := build.Bundle(); err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
	},
}

func init() {
	rootCmd.AddCommand(buildCmd)
}
