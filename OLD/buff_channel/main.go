package main

import (
	"context"
	"errors"
	"fmt"
	"os"
	"strings"
	"sync"
	"time"
)

const NUM_WORKERS = 5
const KILL_TIME_SECONDS = 4

type IFile interface {
	read() (string, error)
}

type File struct {
	path  string
	name  string
	lines int
	words int
	bytes int
}

func NewFile(path string) *File {
	f := &File{path: path}

	j := strings.Split(f.path, "/")
	f.name = j[len(j)-1]

	return f
}

func (f *File) read() (string, error) {
	data, err := os.ReadFile(f.path)
	if err != nil {
		return "", err
	}
	content := string(data)
	f.bytes = len(content)
	f.words = len(strings.Fields(content))
	f.lines = len(strings.Split(content, "\n"))

	return fmt.Sprintf(" %d %d %d %s",
		f.lines, f.words, f.bytes, f.name), nil
}

type WordCount struct {
	ctx   context.Context
	sb    strings.Builder
	files []*File
}

var (
	tasks   = make(chan File, NUM_WORKERS)
	results = make(chan string, NUM_WORKERS)
	wg      sync.WaitGroup
)

func worker(wg *sync.WaitGroup) {
	defer wg.Done()

	for file := range tasks { // Читаем пока канал не закрыт
		r, err := file.read()
		if err != nil {
			results <- fmt.Sprintf("Error reading %s: %v", file.name, err)
		} else {
			results <- r
		}
	}
	return

	select {
	case file := <-tasks:
		r, _ := file.read()
		results <- r
	case <-time.After(KILL_TIME_SECONDS * time.Second):
		return
	}
}

func (wc *WordCount) create(numWorkers int) error {

	for i := 0; i < numWorkers; i++ {
		wg.Add(1)
		go worker(&wg)
	}

	go func() {
		for _, task := range wc.files {
			tasks <- *task
		}
		close(tasks)
	}()

	return nil
}

func (wc *WordCount) complete() string {
	go func() {
		wg.Wait()
		close(results)
	}()

	for r := range results {
		s := fmt.Sprintf(" %s\n", r)
		wc.sb.WriteString(s)
	}

	return wc.sb.String()
}

func run() ([]*File, error) {
	a := os.Args

	if len(a) < 2 {
		return nil, errors.New("")
	}

	var files []*File

	for _, p := range a[1:] {
		f := NewFile(p)

		files = append(files, f)
	}

	return files, nil
}

func main() {

	files, err := run()
	if err != nil {
		fmt.Print("ERROR: Use args: file1, file2, ..., fileN", err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	wc := WordCount{ctx: ctx, sb: strings.Builder{}, files: files}

	wc.create(NUM_WORKERS)
	result := wc.complete()

	fmt.Print(result)
}
