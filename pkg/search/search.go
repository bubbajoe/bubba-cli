package search

import (
	"bufio"
	"fmt"
	"os"
	"regexp"
	"strings"
)

func Reverse(input string) (result string) {
	for _, c := range input {
		result = string(c) + result
	}
	return result
}

type SearchResult struct {
	File     *os.File
	Position int64
}

type SearchParam struct {
	IsRegex  bool
	Match    string
	Filename string
	// ReduceFunc func(string) bool

}

func SearchLine(sp *SearchParam) (<-chan *SearchResult, chan error) {
	// open file
	errc := make(chan error, 1)
	sr := make(chan *SearchResult)
	file, err := os.Open(sp.Filename)
	if err != nil {
		if err == os.ErrNotExist {
			fmt.Println("asdsdsa")
		}
		// s, _ := os.Getwd()
		// fmt.Println(s)
		errc <- err
		return sr, errc
	}

	scanner := bufio.NewScanner(file)
	var line int64 = 1

	for scanner.Scan() {
		if sp.IsRegex {
			re, err := regexp.Compile(sp.Match)
			if err != nil {
				errc <- err
				return sr, errc
			}
			if re.MatchString(scanner.Text()) {
				sr <- &SearchResult{file, line}
			}
		} else {
			if strings.Contains(scanner.Text(), sp.Match) {
				sr <- &SearchResult{
					File:     file,
					Position: line,
				}
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
	r := ChanToSlice(src)
	return r, <-errc
}

func ChanToSlice[T any](chv <-chan T) []T {
	slv := make([]T, 0)
	for {
		v, ok := <-chv
		if !ok {
			return slv
		}
		slv = append(slv, v)
	}
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
			workerId, jobIndex, sp.Match, sp.Filename)
		sc, ec := SearchLine(sp)
		for {
			select {
			case result, ok := <-sc:
				if result != nil {
					fmt.Printf("Job ID (%d-%d) %s:%d", workerId,
						jobIndex, result.File.Name(), result.Position)
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
