package main

import (
	"go-wc-concurrency/internal/entity"
	"go-wc-concurrency/pkg/workerpool"
	"runtime"
	"sync"
)

func main() {
	workersNum := runtime.NumCPU()

	jobs := make([]*entity.Job, 1)
	jobs = append(jobs, entity.NewJob([]byte{1, 2, 3}))

	var wg *sync.WaitGroup
	var (
		jobsCh   = make(chan *entity.Job, workersNum)
		outputCh = make(chan *entity.OutputData, workersNum)
	)

	workerpool.CreateWorkers(wg, workersNum, jobsCh, outputCh)

	go func() {
		for _, job := range jobs {
			jobsCh <- job
		}
		close(jobsCh)
	}()

	workerpool.CompleteWorkers(wg, jobsCh, outputCh)
}
