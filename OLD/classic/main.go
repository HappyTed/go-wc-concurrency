package main

import (
	"bufio"
	"errors"
	"fmt"
	"os"
	"strings"
)

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

func (f *File) read() error {

	file, err := os.Open(f.path)
	if err != nil {
		return errors.New("Ошибка открытия файла")
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)

	for scanner.Scan() {
		f.lines++
		f.bytes += len(scanner.Bytes())

		line := scanner.Text()
		f.words += len(strings.Split(line, " "))

		if err = scanner.Err(); err != nil {
			return fmt.Errorf("Ошибка чтения файла: %w", err)
		}
	}

	return nil
}

type WordCount struct {
	sb    strings.Builder
	files []*File
}

func (wc *WordCount) calculate() string {

	for _, f := range wc.files {
		f.read()
		s := fmt.Sprintf(" %d %d %d %s\n", f.lines, f.words, f.bytes, f.name)
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

	wc := WordCount{sb: strings.Builder{}, files: files}

	res := wc.calculate()

	fmt.Println(res)
}
