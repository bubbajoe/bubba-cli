package search

import (
	"bufio"
	"context"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/bubbajoe/bubba-cli/pkg/util"
)

type SearchResult struct {
	Name     string
	Position int64
	Text     string
}

type SearchParam struct {
	IsRegex  bool
	Query    string
	FilePath string
	// ReduceFunc func(string) bool

}

func SearchLine(sp *SearchParam) (<-chan *SearchResult, <-chan error) {
	// open file
	errc := make(chan error)
	defer close(errc)
	sr := make(chan *SearchResult)
	defer close(sr)
	var rdr io.Reader = nil
	var name string

	if sp.FilePath == "-" {
		rdr = os.Stdin
		name = "{stdin}"
	} else {
		file, err := os.Open(sp.FilePath)
		if err != nil {
			if err == os.ErrNotExist {
				fmt.Printf("file does not exist: %s\n", sp.FilePath)
			}
			// s, _ := os.Getwd()
			fmt.Println(err)
			errc <- err
			return sr, errc
		}
		name = file.Name()
		defer file.Close()
		rdr = file
	}

	scanner := bufio.NewScanner(rdr)
	var line int64 = 1

	for scanner.Scan() {
		text := string(scanner.Bytes())
		if sp.IsRegex {
			re, err := regexp.Compile(sp.Query)
			if err != nil {
				errc <- err
				return sr, errc
			}
			if re.MatchString(text) {
				sr <- &SearchResult{name, line, text}
			}
		} else {
			// fmt.Println("Searching : ", sp.Query, string(scanner.Text()))
			if strings.Contains(text, sp.Query) {
				// fmt.Println("Found : ", sp.Query, string(scanner.Text()))
				sr <- &SearchResult{name, line, text}
			}
		}

		line++
	}

	return sr, errc
}

func SearchLineMany(ctx context.Context, srs []*SearchParam, workers int) ([]*SearchResult, error) {
	if workers < 1 {
		workers = 1
	}

	spc := make(chan *SearchParam)
	// in case its a buffered channel,
	// we will use a goroutine to fill it
	go func() {
		for _, sr := range srs {
			spc <- sr
		}
		close(spc)
	}()
	src, errc := SearchLineManyChan(ctx, spc, workers)
	r := util.ChanToSlice(src)
	return r, <-errc
}

func SearchLineManyChan(ctx context.Context, spc <-chan *SearchParam, workers int) (<-chan *SearchResult, <-chan error) {
	src := make(chan *SearchResult)
	errc := make(chan error, 1)
	// create n workers to consme the data
	if workers < 1 {
		workers = 1
	}

	for i := 0; i < workers; i++ {
		go searchLineWorker(i, ctx, spc, src, errc)
	}

	return src, errc
}

func searchLineWorker(
	workerId int,
	ctx context.Context,
	spc <-chan *SearchParam,
	src chan<- *SearchResult,
	errc chan<- error,
) {
	jobIndex := 0
	defer fmt.Printf("Worker (%d) / Job (%d/%d) - CLOSED\n",
		workerId, workerId, jobIndex)
	for sp := range spc {
		fmt.Printf("Worker (%d) / Job (%d/%d) - Searching for '%s' in '%s'\n",
			workerId, workerId, jobIndex, sp.Query, sp.FilePath)
		sc, ec := SearchLine(sp)
		go func() {
			for {
				select {
				case <-ctx.Done():
					return
				case result, ok := <-sc:
					if !ok {
						// fmt.Printf("Job ID (%d-%d) - DONE\n", workerId, jobIndex)
						errc <- nil
						return
					}
					if result != nil {
						// fmt.Printf("Job ID (%d-%d) %s:%d", workerId,
						// 	jobIndex, result.Name, result.Position)
						src <- result
					} else {
						// fmt.Printf("Job ID (%d-%d) - No result\n", workerId, jobIndex)
					}
				case err := <-ec:
					if err != nil {
						errc <- err
						// fmt.Printf("Job ID (%d-%d) ERROR: %s\n", workerId, jobIndex, err)
					}
					return
				}
			}
		}()
	}
}

func RecursiveParams(dir, query string, isRegex bool) chan *SearchParam {
	c := make(chan *SearchParam)
	go func() {
		filepath.WalkDir(dir, func(path string, de os.DirEntry, err error) error {
			if err != nil {
				return err
			}
			if de.IsDir() {
				return nil
			}
			c <- &SearchParam{isRegex, query, filepath.Join(path, de.Name())}
			return nil
		})
		close(c)
	}()
	return c
}

func DirectoryParams(dir, query string, isRegex bool) chan *SearchParam {
	c := make(chan *SearchParam)
	go func() {
		defer close(c)
		files, err := os.ReadDir(dir)
		if err != nil {
			return
		}
		for _, file := range files {
			if file.IsDir() {
				continue
			}
			fmt.Println("Adding file to search: ", file.Name())
			c <- &SearchParam{isRegex, query, filepath.Join(dir, file.Name())}
		}
	}()
	return c
}

func FileParam(file, query string, isRegex bool) *SearchParam {
	return &SearchParam{isRegex, query, file}
}

func IsDir(path string) (bool, error) {
	f, err := os.Open(path)
	if err != nil {
		return false, err
	}
	defer f.Close()
	fi, err := f.Stat()
	if err != nil {
		return false, err
	}
	return fi.IsDir(), nil
}
