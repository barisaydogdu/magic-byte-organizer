package cli

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "magicsort",
	Short: "MagicSort - An intelligent file classification tool.",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Welcome to the MagicSort CLI vehicle!")
	},
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
