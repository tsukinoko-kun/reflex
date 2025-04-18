package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/tsukinoko-kun/reflex/internal/build"
	"github.com/tsukinoko-kun/reflex/internal/lint"
	"github.com/tsukinoko-kun/reflex/internal/new"
)

var initCmd = &cobra.Command{
	Use:   "init",
	Short: "initialize a new reflex application",

	Run: func(cmd *cobra.Command, args []string) {
		if err := new.New(); err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		_ = lint.Format()

		_ = build.Bundle()
	},
}

func init() {
	rootCmd.AddCommand(initCmd)
}
