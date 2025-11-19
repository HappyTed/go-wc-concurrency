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

type Job struct {
	r io.Reader
}

func NewJob(r io.Reader) *Job {
	return &Job{r: r}
}

// Scan io.Reader line by line and calculate: lines, words and bytes
func (j *Job) Calculate() (*entity.OutputData, error) {
	res := &entity.OutputData{}

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
