package bc

import (
	"fmt"
	"runtime"

	"github.com/bubbajoe/bubba-cli/pkg/search"
	"github.com/bubbajoe/bubba-cli/pkg/util"
	"github.com/spf13/cobra"
)

var searchCmd = &cobra.Command{
	Use:     "search",
	Aliases: []string{"S"},
	Short:   "search for a string in a file(s)",
	Args:    cobra.MinimumNArgs(2),
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

		srs, err := search.SearchLineMany(
			util.SliceMap(args[:len(args)-1], func(fp string) *search.SearchParam {
				return &search.SearchParam{
					IsRegex:  isRegex,
					Match:    args[0],
					FilePath: fp,
				}
			}), threads)
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
