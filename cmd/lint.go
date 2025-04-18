package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/tsukinoko-kun/reflex/internal/lint"
)

var (
	lintFix bool
	lintCmd = &cobra.Command{
		Use:   "lint",
		Short: "A brief description of your command",

		Run: func(cmd *cobra.Command, args []string) {
			if err := lint.Run(lintFix); err != nil {
				fmt.Println(err)
				os.Exit(1)
			}
		},
	}
)

func init() {
	rootCmd.AddCommand(lintCmd)
	lintCmd.Flags().BoolVar(&lintFix, "fix", lintFix, "Fix lint errors if possible")
}
