package main

import (
	"flag"
	"fmt"
	"go-wc-concurrency/internal/config"
	"go-wc-concurrency/internal/entity"
	"go-wc-concurrency/internal/logic"
	"go-wc-concurrency/pkg"
	"log"
	"os"
	"sync"
)

func Printing(out []*entity.OutputData) {
	// var sb strings.Builder
	for _, r := range out {
		fmt.Println(*r)
	}
}

func main() {

	var wg sync.WaitGroup
	var jobsCh = make(chan logic.IJob)
	var outputCh = make(chan *entity.OutputData)
	var closeFilesCh = make(chan func() error)

	cfg := config.ReadConfig()

	if len(cfg.FilesPath) > 0 {
		wg.Add(1)
		go func() {
			defer close(jobsCh)
			defer close(closeFilesCh)
			defer wg.Done()

			for _, path := range flag.Args() {
				file, err := os.Open(path)
				if err != nil {
					log.Fatal("Failed to open file:", path)
				}
				closeFilesCh <- file.Close

				job := logic.NewJob(file)
				jobsCh <- job
			}
		}()
	} else { // пытаемся считать из stdin (передача через pipe)
		wg.Add(1)
		go func() {
			job := logic.NewJob(os.Stdin)
			jobsCh <- job

			close(jobsCh)
			close(closeFilesCh)

			wg.Done()
		}()
	}

	workerFunc := func(wg *sync.WaitGroup, jobs chan logic.IJob, out chan *entity.OutputData) {
		defer wg.Done()

		for job := range jobs {
			result, _ := job.Calculate()
			out <- result
		}

		// close(out)
	}

	wp, err := pkg.MakePool(
		pkg.WithWorkersCount(cfg.NumWorkers),
		pkg.WithJobsChannel(jobsCh),
		pkg.WithOutputChannel(outputCh),
		pkg.WithWorkerFunc(workerFunc),
	)
	if err != nil {
		log.Fatal(err)
	}

	var outs []*entity.OutputData
	wg.Add(1)
	go func() {
		for out := range outputCh {
			outs = append(outs, out)
		}
		wg.Done()
	}()

	err = wp.CreateWorkers()
	if err != nil {
		log.Fatal(err)
	}

	err = wp.Complete()
	if err != nil {
		log.Fatal(err)
	}

	wg.Add(1)
	go func() {
		for close := range closeFilesCh {
			close()
		}
		wg.Done()
	}()

	wg.Wait()

	Printing(outs)
}
