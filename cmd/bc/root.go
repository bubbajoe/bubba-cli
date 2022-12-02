package bc

import (
	"fmt"
	"os"

	"github.com/bubbajoe/bubba-cli/pkg/i9e"
	"github.com/fatih/color"
	"github.com/spf13/cobra"
)

var currentVersion string
var rootCmd = &cobra.Command{
	Use:   "bc",
	Short: "bc - utils for days",
	Long: `
Bubba-CLI is a utiliy CLT made for, and by the one and only BubbaJoe!

this CLT currently contains utliies searching, encode/decode and networking`,
	Run: func(cmd *cobra.Command, args []string) {
		c := color.New(color.FgRed)
		version, err := cmd.Flags().GetBool("version")
		if err != nil {
			c.Printf("Bubba CLI - Erro"+
				"r: '%s'\n", err)
			return
		}
		if version {
			c.Printf("Bubba CLI - Version: '%s'\n", currentVersion)
			return
		}
		// LOGO :O
		c.Printf("Thanks for using bubba-cli (%s)\n%s\n\n", currentVersion,
			"▄▄▄▄· ▄• ▄▌▄▄▄▄· ▄▄▄▄·  ▄▄▄·\n▐█ ▀█▪█▪██▌▐█ ▀█▪▐█ ▀█▪▐█ ▀█\n"+
				"▐█▀▀█▄█▌▐█▌▐█▀▀█▄▐█▀▀█▄▄█▀▀█\n██▄▪▐█▐█▄█▌██▄▪▐███▄▪▐█▐█ ▪▐▌\n"+
				"·▀▀▀▀  ▀▀▀ ·▀▀▀▀ ·▀▀▀▀  ▀  ▀")
		i9e.StartInteractivePrompt()
	},
}

func Execute(version string) {
	currentVersion = version
	rootCmd.Flags().BoolP("version", "v",
		false, "print the version")
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintf(color.Output, "Bubba CLI - Error: '%s'", err)
		os.Exit(1)
	}
}
