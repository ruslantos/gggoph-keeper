package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"goph-keeper/client/internal/tui"
)

var (
	apiURL = "http://localhost:8080"
)

var rootCmd = &cobra.Command{
	Use:   "gophkeeper",
	Short: "CLI для GophKeeper",
	Run: func(cmd *cobra.Command, args []string) {
		if err := tui.NewTUICommand().Execute(); err != nil {
			fmt.Println("Ошибка запуска программы:", err)
			os.Exit(1)
		}
	},
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}
