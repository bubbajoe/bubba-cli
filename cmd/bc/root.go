package bc

import (
	"fmt"
	"os"

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
		color.Set(color.FgRed)
		fmt.Fprintf(color.Output,
			"Thanks for using bubba-cli (%s) %s\n\n", currentVersion, `
▄▄▄▄· ▄• ▄▌▄▄▄▄· ▄▄▄▄·  ▄▄▄·
▐█ ▀█▪█▪██▌▐█ ▀█▪▐█ ▀█▪▐█ ▀█
▐█▀▀█▄█▌▐█▌▐█▀▀█▄▐█▀▀█▄▄█▀▀█
██▄▪▐█▐█▄█▌██▄▪▐███▄▪▐█▐█ ▪▐▌
·▀▀▀▀  ▀▀▀ ·▀▀▀▀ ·▀▀▀▀  ▀  ▀ `)
		color.Set(color.Reset)
	},
}

func Execute(version string) {
	currentVersion = version
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintf(color.Output, "Bubba CLI - Error: '%s'", err)
		os.Exit(1)
	}
}
