package workerpool

import (
	"fmt"
	"go-wc-concurrency/internal/entity"
	"sync"
)

// type WoorkerPool struct {
// 	woorkerNum int
// 	worker     func(wg *sync.WaitGroup)
// 	jobsCh     <-chan any
// 	resultsCh  chan<- any
// 	errorsCh   chan error
// }

// type Worker func(wg *sync.WaitGroup)

// func WithWorker(wp *WoorkerPool)

// func New(workers int, jobsCh <-chan any, resultsCh chan<- any, errorsCh chan error) *WoorkerPool {
// 	wp := &WoorkerPool{
// 		woorkerNum: workers,
// 		jobsCh:     jobsCh,
// 		resultsCh:  resultsCh,
// 		errorsCh:   errorsCh,
// 	}
// 	return wp
// }

func worker(
	wg *sync.WaitGroup,
	jobsCh chan *entity.Job,
	outputCh chan *entity.OutputData,
) {
	defer wg.Done()

	for job := range jobsCh { // Читаем пока канал не закрыт
		// TODO smth
		fmt.Println(job)
		outputCh <- &entity.OutputData{}
	}
}

func CreateWorkers(
	wg *sync.WaitGroup,
	numWorkers int,
	jobsCh chan *entity.Job,
	outputCh chan *entity.OutputData,
) error {

	for i := 0; i < numWorkers; i++ {
		wg.Add(1)
		go worker(wg, jobsCh, outputCh)
	}

	return nil
}

func CompleteWorkers(
	wg *sync.WaitGroup,
	jobsCh chan *entity.Job,
	outputCh chan *entity.OutputData,
) string {
	go func() {
		wg.Wait()
		close(outputCh)
	}()

	for i := range outputCh {
		fmt.Println(i)
		// TODO: processing
	}

	return ""
}
