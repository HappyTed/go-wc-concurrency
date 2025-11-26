package pkg

import (
	"errors"
	"go-wc-concurrency/internal/entity"
	"go-wc-concurrency/internal/logic"
	"sync"
)

type WorkerFunc func(*sync.WaitGroup, chan logic.IJob, chan *entity.OutputData)

type WoorkerPool struct {
	wg         *sync.WaitGroup
	numWorkers uint8
	jobsCh     chan logic.IJob
	outputCh   chan *entity.OutputData
	workerFunc WorkerFunc
}

type Option func(*WoorkerPool) error

func WithWaitGroup(wg *sync.WaitGroup) Option {
	return func(wp *WoorkerPool) error {
		if wg == nil {
			return errors.New("WaitGroup can't be nil pointer")
		}
		wp.wg = wg
		return nil
	}
}

func WithWorkersCount(c uint8) Option {
	return func(wp *WoorkerPool) error {
		wp.numWorkers = c
		return nil
	}
}

func WithJobsChannel(j chan logic.IJob) Option {
	return func(wp *WoorkerPool) error {
		if j == nil {
			return errors.New("Jobs channel can't be nil chan")
		}
		wp.jobsCh = j
		return nil
	}
}

func WithOutputChannel(op chan *entity.OutputData) Option {
	return func(wp *WoorkerPool) error {
		if op == nil {
			return errors.New("Result's channel can't be nil chan")
		}
		wp.outputCh = op
		return nil
	}
}

func WithWorkerFunc(f WorkerFunc) Option {
	return func(wp *WoorkerPool) error {
		if f == nil {
			return errors.New("Worker Func can't be nil pointer")
		}
		wp.workerFunc = f
		return nil
	}
}

func MakePool(opts ...Option) (*WoorkerPool, error) {
	wp := &WoorkerPool{wg: &sync.WaitGroup{}}
	for _, opt := range opts {
		if err := opt(wp); err != nil {
			return nil, err
		}
	}
	return wp, nil
}

func (wp *WoorkerPool) CreateWorkers() error {

	for i := 0; i < int(wp.numWorkers); i++ {
		wp.wg.Add(1)
		go wp.workerFunc(wp.wg, wp.jobsCh, wp.outputCh)
	}

	return nil
}

func (wp *WoorkerPool) Complete() error {
	go func() {
		wp.wg.Wait()

		close(wp.outputCh)
	}()

	return nil
}
