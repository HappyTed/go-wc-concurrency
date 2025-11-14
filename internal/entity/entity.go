package entity

import (
	"errors"
	"fmt"
	"os"
)

type OutputData struct {
	lines int
	words int
	bytes int
}

type JobType uint8

const (
	BYTES = iota
	FILE
)

var JobTypeEnum = map[JobType]JobType{
	BYTES: BYTES,
	FILE:  FILE,
}

type Job struct {
	InputData []byte
	Type      JobType
}

func NewJob(inputData []byte) *Job {
	j := &Job{InputData: inputData}

	_, error := os.Stat(string(inputData))
	if !errors.Is(error, os.ErrNotExist) { // это файл
		j.Type = JobTypeEnum[FILE]
	} else {
		fmt.Println("File does'not exist") // TODO: нужна более логичная проверка на то что это файл path или например вывод `cat file.txt | ./my_prog`
		j.Type = JobTypeEnum[BYTES]
	}

	return j
}
