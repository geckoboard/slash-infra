package main

import (
	"log"

	"github.com/joho/godotenv"
	"github.com/spf13/cobra"

	"github.com/geckoboard/slash-infra/cmd/http"
)

var (
	rootCmd = &cobra.Command{
		Use:   "slash-infra",
		Short: "A slackbot for working with AWS infrastructure",
	}
)

func main() {
	log.SetFlags(log.Ldate | log.Ltime | log.Lmicroseconds | log.Llongfile)

	// In development it's easier to store environment variables in a .env folder
	godotenv.Load()

	rootCmd.AddCommand(http.Command())

	rootCmd.Execute()
}
