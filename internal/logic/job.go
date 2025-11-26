package logic

import (
	"bufio"
	"go-wc-concurrency/internal/entity"
	"io"
	"strings"
)

type IJob interface {
	Calculate() (*entity.OutputData, error)
}

type CloseReader func() error

type Job struct {
	name  string
	r     io.Reader
	close CloseReader
}

func NewJob(name string, r io.Reader, close CloseReader) *Job {
	return &Job{name: name, r: r, close: close}
}

// Scan io.Reader line by line and calculate: lines, words and bytes
func (j *Job) Calculate() (*entity.OutputData, error) {
	defer j.close() // если это файл, закрываем его

	res := &entity.OutputData{Name: j.name}

	scanner := bufio.NewScanner(j.r)
	scanner.Split(bufio.ScanLines)

	for scanner.Scan() {
		if scanner.Err() != nil {
			return nil, scanner.Err()
		}
		res.Lines++
		res.Bytes += uint(len(scanner.Bytes()))
		res.Words += uint(len(strings.Split(scanner.Text(), " ")))
	}

	return res, nil
}
