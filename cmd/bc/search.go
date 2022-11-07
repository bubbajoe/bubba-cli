package bc

import (
	"fmt"
	"runtime"

	"github.com/bubbajoe/bubba-cli/pkg/search"
	"github.com/spf13/cobra"
)

var searchCmd = &cobra.Command{
	Use:     "search",
	Aliases: []string{"S"},
	Short:   "search for a string in a file(s)",
	Args:    cobra.ExactArgs(2),
	RunE: func(cmd *cobra.Command, args []string) error {
		isRegex, err := cmd.Flags().GetBool("regex")
		if err != nil {
			return err
		}
		threads, err := cmd.Flags().GetInt("threads")
		if err != nil {
			return err
		}

		fmt.Println("searching for '", args[0], "' in filename:", args[1])

		srs, err := search.SearchLineMany([]*search.SearchParam{
			{
				IsRegex:  isRegex,
				Match:    args[1],
				Filename: args[0],
			},
		}, threads)
		if err != nil {
			return err
		}
		for _, sr := range srs {
			fmt.Println(sr)
		}
		return err
	},
}

func init() {
	searchCmd.Flags().BoolP("regex", "r",
		false, "regex search")
	searchCmd.Flags().IntP("threads", "t",
		runtime.NumCPU(), "max of threads to use")
	rootCmd.AddCommand(searchCmd)
}
