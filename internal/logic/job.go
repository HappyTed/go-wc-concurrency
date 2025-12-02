package logic

import (
	"bufio"
	"errors"
	"go-wc-concurrency/internal/config"
	"go-wc-concurrency/internal/entity"
	"io"
)

type IJob interface {
	Calculate() (*entity.OutputData, error)
}

type CloseReader func() error

type Job struct {
	br        *bufio.Reader
	options   config.OptionFlag
	closeFunc CloseReader
	result    *entity.OutputData
}

type Option func(*Job) error

func WithBytesCount(c uint64) Option {
	return func(j *Job) error {
		if j.result == nil {
			return errors.New("entity.OutputData is nil")
		}
		j.result.Bytes = c
		return nil
	}
}

func WithFlags(o config.OptionFlag) Option {
	return func(j *Job) error {
		j.options = o
		return nil
	}
}

func WithCloseFunc(f CloseReader) Option {
	return func(j *Job) error {
		if f == nil {
			return errors.New("close reader func is nil")
		}
		j.closeFunc = f
		return nil
	}
}

func WithName(name string) Option {
	return func(j *Job) error {
		if j.result == nil {
			return errors.New("entity.OutputData is nil")
		}
		j.result.Name = name
		return nil
	}
}

func NewJob(r io.Reader, options ...Option) (*Job, error) {

	j := &Job{
		br:        bufio.NewReader(r),
		closeFunc: func() error { return nil },
		options:   config.BYTES | config.LINES | config.WORDS,
		result:    &entity.OutputData{},
	}

	for _, opt := range options {
		if err := opt(j); err != nil {
			return nil, err
		}
	}

	return j, nil
}

func (j *Job) Calculate() (*entity.OutputData, error) {
	defer j.closeFunc() // если это файл, закрываем его

	// если нужны только байты, их могли передать вначале (Если это файл, то количество байт заранее известно)
	if j.options&config.BYTES != 0 &&
		(j.options&config.LINES == 0 && j.options&config.WORDS == 0) &&
		j.result.Bytes != 0 {
		return j.result, nil
	}

	var (
		lines, bytes, words uint64
		isWord              bool
		buffer              = make([]byte, 1024)
	)

	for {
		n, err := j.br.Read(buffer)

		if err == io.EOF {
			break
		} else if err != nil {
			return nil, err
		}

		if n == 0 {
			break
		}

		bytes += uint64(n)

		chunk := buffer[0:n]
		for _, b := range chunk {
			switch b {
			case ' ', '\t', '\n', '\r', '\f', '\v':
				if isWord {
					words++
					isWord = false
				}
				if b == '\n' {
					lines++
				}
			default:
				isWord = true
			}
		}
	}

	if isWord {
		words++
	}

	j.result.Bytes = bytes
	j.result.Lines = lines
	j.result.Words = words

	return j.result, nil
}
