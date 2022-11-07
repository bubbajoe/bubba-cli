package bc

import (
	"github.com/fatih/color"
	"github.com/spf13/cobra"
)

var onlyDigits bool
var inspectCmd = &cobra.Command{
	Use:     "init",
	Aliases: []string{"init", "I"},
	// Short:   "interactive",
	Args: cobra.NoArgs,
	Run: func(cmd *cobra.Command, args []string) {
		c := color.New(color.FgRed)
		c.Println("Interative search started")
	},
}

func init() {
	inspectCmd.Flags().BoolVarP(&onlyDigits, "digits", "d", false, "Count only digits")
	rootCmd.AddCommand(inspectCmd)
}
