package cli

import (
	"log"
	"os"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "magicsort",
	Short: "MagicSort - An intelligent file classification tool.",
	Run: func(cmd *cobra.Command, args []string) {
		log.Println("Welcome to the MagicSort CLI vehicle!")
	},
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		log.Println(err)
		os.Exit(1)
	}
}
