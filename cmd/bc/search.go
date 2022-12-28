package bc

import (
	"context"
	"fmt"
	"runtime"

	"github.com/bubbajoe/bubba-cli/pkg/search"
	"github.com/spf13/cobra"
	"go.uber.org/zap"
)

var (
	log, _ = zap.NewDevelopment()
)

var searchCmd = &cobra.Command{
	Use:     "search",
	Aliases: []string{"S"},
	Short:   "search for a string in a file(s)",
	Args:    cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		ctx, _ := context.WithTimeout(context.Background(), 10_000)
		isRegex, err := cmd.Flags().GetBool("regex")
		if err != nil {
			return err
		}
		threads, err := cmd.Flags().GetInt("threads")
		if err != nil {
			return err
		}

		path, err := cmd.Flags().GetString("path")
		if err != nil {
			return err
		}

		q := args[0]
		// check if path is a directory or file
		isDir, err := search.IsDir(path)
		if err != nil {
			return err
		}

		var sr_ch <-chan *search.SearchResult
		var err_ch <-chan error

		if isDir {
			var sp_ch chan *search.SearchParam
			// search directory
			isRecursive, err := cmd.Flags().GetBool("recursive")
			if err != nil {
				return err
			}
			if isRecursive {
				log.Debug("searching recursively",
					zap.String("p", path),
					zap.String("q", q),
					zap.Bool("regex", isRegex),
				)
				sp_ch = search.RecursiveParams(path, q, isRegex)
			} else {
				log.Debug("searching non-recursively",
					zap.String("p", path),
					zap.String("q", q),
					zap.Bool("regex", isRegex),
				)
				sp_ch = search.DirectoryParams(path, q, isRegex)
			}
			sr_ch, err_ch = search.SearchLineManyChan(ctx, sp_ch, threads)
		} else {
			log.Debug("searching file",
				zap.String("p", path),
				zap.String("q", q),
				zap.Bool("R", isRegex),
			)
			sp := search.FileParam(path, q, isRegex)
			sr_ch, err_ch = search.SearchLine(sp)
		}

		count := 0
		for {
			// fmt.Println("RESULT TIME!", len(sr_ch), len(err_ch))
			select {
			case sr := <-sr_ch:
				// if !ok {
				// 	fmt.Print(sr)
				// 	return nil
				// }
				fmt.Printf("FOUND %s:%d\n", sr.Name, sr.Position)
			case err := <-err_ch:
				count++
				if err != nil {
					fmt.Println(err)
				}
				if count == threads+4 {
					fmt.Println("DONZO")
					return nil
				}
			}
		}

		// srs, err := search.SearchLineMany(
		// 	util.SliceMap(args[:len(args)-1], func(fp string) *search.SearchParam {
		// 		return &search.SearchParam{
		// 			IsRegex:  isRegex,
		// 			Query:    args[0],
		// 			FilePath: fp,
		// 		}
		// 	}), threads)
		// if err != nil {
		// 	return err
		// }
		// return err
	},
}

func init() {
	searchCmd.Flags().StringP("path", "p",
		".", "path to file or directory")
	searchCmd.Flags().BoolP("regex", "R",
		false, "regex search")
	searchCmd.Flags().BoolP("recursive", "r",
		false, "recursive search") // filepath must be a directory
	searchCmd.Flags().IntP("threads", "t",
		runtime.NumCPU()/2, "max of threads to use")
	rootCmd.AddCommand(searchCmd)
}
