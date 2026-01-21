package main

import (
	"fmt"
	"go-wc-concurrency/internal/config"
	"go-wc-concurrency/internal/entity"
	"go-wc-concurrency/internal/logic"
	"go-wc-concurrency/pkg"
	"log"
	"os"
	"strings"
	"sync"
)

func printing(flags config.OptionFlag, out []*entity.OutputData) {
	var header strings.Builder
	var sep strings.Builder

	if flags&config.LINES != 0 {
		header.WriteString(fmt.Sprintf("%-6s", "lines"))
		sep.WriteString(strings.Repeat("-", len("lines")+1))
	}
	if flags&config.WORDS != 0 {
		header.WriteString(fmt.Sprintf("%-6s", "words"))
		sep.WriteString(strings.Repeat("-", len("words")+1))
	}
	if flags&config.BYTES != 0 {
		header.WriteString(fmt.Sprintf("%-6s", "bytes"))
		sep.WriteString(strings.Repeat("-", len("bytes")+1))
	}

	fmt.Println(header.String())
	fmt.Println(sep.String())

	var sb strings.Builder
	var fullResult = entity.OutputData{}
	for _, r := range out {

		fullResult.Lines += r.Lines
		fullResult.Words += r.Words
		fullResult.Bytes += r.Bytes

		if flags&config.LINES != 0 {
			sb.WriteString(fmt.Sprintf("%5d", r.Lines))
		}
		if flags&config.WORDS != 0 {
			sb.WriteString(fmt.Sprintf("%6d", r.Words))
		}
		if flags&config.BYTES != 0 {
			sb.WriteString(fmt.Sprintf("%6d", r.Bytes))
		}
		sb.WriteString(fmt.Sprintf(" %s", r.Name))

		// REV: принты юзать такое, только логер
		fmt.Println(sb.String())
		sb.Reset()
	}

	if len(out) > 1 {
		if flags&config.LINES != 0 {
			sb.WriteString(fmt.Sprintf("%5d", fullResult.Lines))
		}
		if flags&config.WORDS != 0 {
			sb.WriteString(fmt.Sprintf("%6d", fullResult.Words))
		}
		if flags&config.BYTES != 0 {
			sb.WriteString(fmt.Sprintf("%6d", fullResult.Bytes))
		}
		fmt.Println()
		sb.WriteString(fmt.Sprintf(" %s", "total"))
		fmt.Println(sb.String())
	}
}

// REV: если уж начал правильно по папкам раскидывать, нужно вынести запуск в cmd/app-name/main.go
// там только вызов инициализирующей функции.
// А саму инициализирующую функцию в internal/application/app-name(или просто app, если сервис один)/app.go.
// Тут инициализруешь конфиги, репозитории и все что вообще нужно
func main() {

	cfg := config.ReadConfig()

	var (
		wg       sync.WaitGroup
		results  []*entity.OutputData
		jobsCh   = make(chan logic.IJob, cfg.NumWorkers)
		outputCh = make(chan *entity.OutputData, cfg.NumWorkers)
	)

	// Создаём функцию-обработчик задач
	workerFunc := func(wg *sync.WaitGroup, jobsCh chan logic.IJob, outCh chan *entity.OutputData) {
		// REV: wg лучше юзать не в самой workFunc, а перед тем где она вызывается,
		// проще будет отслеживать закрыл или нет.
		// А еще есть errgroup, ее тоже можно как авто закрывающуюся wg юзать.
		// В го 1.25 появился метод wg.Go(), который сам будет закрывать wg
		defer wg.Done()

		for job := range jobsCh {
			result, err := job.Calculate()
			if err != nil {
				fmt.Println(result.Name, "counting error:", err)
				continue
			}
			outCh <- result
		}
	}

	// Создаём горутины на чтение результатов выполнения задач
	wg.Add(1)
	go func() {
		for out := range outputCh {
			results = append(results, out)
		}
		wg.Done()
	}()

	// REV: With функции обычно используются, когда делаешь либу, которая будет использоваться
	// как пакет, а так было бы достаточно и просто структуру сделать
	//
	// Заводим пул воркеров для обработки задач
	wp, err := pkg.MakePool(
		pkg.WithWorkersCount(cfg.NumWorkers),
		pkg.WithJobsChannel(jobsCh),
		pkg.WithOutputChannel(outputCh),
		pkg.WithWorkerFunc(workerFunc),
	)
	if err != nil {
		log.Fatal(err)
	}

	err = wp.CreateWorkers()
	if err != nil {
		// REV: log устаревший, норм пацаны юзают slog, там есть куча классных реализаций,
		// я для дева такую использую "github.com/golang-cz/devslog"
		log.Fatal(err)
	}
	// REV: этот ебейший if else лучше разбить на функции и вызывать их, легче будет читаться и тестить потом
	//
	// Отправляем задачи на обработку
	if len(cfg.Files) > 0 { // если это файлы, заводи jobs на их чтение
		wg.Add(1)

		// REV: а зачем эта горутина? ты же все равно ждешь ее завершения сразу после
		// было бы ок, если бы читал файлы в нескольких горутинах, а так все равно по очереди идешь
		go func() {
			defer close(jobsCh)
			defer wg.Done()

			for _, path := range cfg.Files {
				file, err := os.Open(path)
				if err != nil {
					// REV: логгер
					fmt.Println("error: file: No such file or directory`:", path)
					continue
				}

				info, err := file.Stat()
				if err != nil {
					fmt.Println("error: failed get file info:", err)
					continue
				}

				// REV: тоже самое с with, но в целом ок
				job, err := logic.NewJob(
					file,
					logic.WithName(file.Name()),
					logic.WithFlags(cfg.Options),
					logic.WithCloseFunc(file.Close),
					logic.WithBytesCount(uint64(info.Size())),
				)
				if err != nil {
					fmt.Println("failed read file:", err)
					continue
				}

				jobsCh <- job
			}
		}()
	} else { // иначе пытаемся считать из stdin (или передача через linux pipe)
		wg.Add(1)
		go func() {
			defer wg.Done()
			defer close(jobsCh)

			job, err := logic.NewJob(
				os.Stdin,
				logic.WithFlags(cfg.Options),
			)
			if err != nil {
				fmt.Println("failed read file:", err)
				return
			}

			jobsCh <- job
		}()
	}

	// ждём, когда все задачи будут выполнены
	wp.Complete()

	// ждём завершения обработки результатов
	wg.Wait()

	printing(cfg.Options, results)
}
