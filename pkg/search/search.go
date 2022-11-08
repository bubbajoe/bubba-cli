package search

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"regexp"
	"strings"

	"github.com/bubbajoe/bubba-cli/pkg/util"
)

func Reverse(input string) (result string) {
	for _, c := range input {
		result = string(c) + result
	}
	return result
}

type SearchResult struct {
	Name     string
	Position int64
}

type SearchParam struct {
	IsRegex  bool
	Match    string
	FilePath string
	// ReduceFunc func(string) bool

}

func SearchLine(sp *SearchParam) (<-chan *SearchResult, chan error) {
	// open file
	errc := make(chan error, 1)
	sr := make(chan *SearchResult)
	var rdr io.Reader = nil
	var name string
	if sp.FilePath == "-" {
		rdr = os.Stdin
		name = ":stdin"
	} else {
		file, err := os.Open(sp.FilePath)
		if err != nil {
			if err == os.ErrNotExist {
				fmt.Println("asdsdsa")
			}
			// s, _ := os.Getwd()
			// fmt.Println(s)
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
		if sp.IsRegex {
			re, err := regexp.Compile(sp.Match)
			if err != nil {
				errc <- err
				return sr, errc
			}
			if re.MatchString(scanner.Text()) {
				sr <- &SearchResult{name, line}
			}
		} else {
			if strings.Contains(scanner.Text(), sp.Match) {
				sr <- &SearchResult{name, line}
			}
		}

		line++
	}
	close(sr)

	return sr, errc
}

func SearchLineMany(srs []*SearchParam, workers int) ([]*SearchResult, error) {
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
	src, errc := SearchLineManyChan(spc, workers)
	r := util.ChanToSlice(src)
	return r, <-errc
}

func SearchLineManyChan(spc <-chan *SearchParam, workers int) (<-chan *SearchResult, <-chan error) {
	src := make(chan *SearchResult)
	errc := make(chan error, 1)
	// create n workers to consme the data
	if workers < 1 {
		workers = 1
	}

	for i := 0; i < workers; i++ {
		go searchLineWorker(i, spc, src, errc)
	}

	return src, errc
}

func searchLineWorker(workerId int, spc <-chan *SearchParam, src chan<- *SearchResult, errc chan<- error) {
	jobIndex := 0
	for sp := range spc {
		fmt.Printf("Worker (%d) / Job ID (%d-%d) - Searching for '%s' in '%s'\n", workerId,
			workerId, jobIndex, sp.Match, sp.FilePath)
		sc, ec := SearchLine(sp)
		for {
			select {
			case result, ok := <-sc:
				if result != nil {
					fmt.Printf("Job ID (%d-%d) %s:%d", workerId,
						jobIndex, result.Name, result.Position)
					src <- result
				} else {
					fmt.Printf("Job ID (%d-%d) - No result", workerId, jobIndex)
				}
				if !ok {
					close(src)
					errc <- nil
					return
				}
			case err := <-ec:
				if err != nil {
					// print the error
					fmt.Printf("Job ID (%d-%d) ERROR: %s", workerId, jobIndex, err)
				}
				close(src)
				errc <- err
				return
			}
		}
	}
}
