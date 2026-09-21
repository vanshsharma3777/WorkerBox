package worker

import (
	"fmt"
	"sync"
	"time"

	"github.com/vanshsharma3777/WorkerBox/helper"
	"github.com/vanshsharma3777/WorkerBox/internal/queue"
	"github.com/vanshsharma3777/WorkerBox/models"
)

type Worker struct {
	Handlers    map[string]func(models.Job) error
	Queue       *queue.Queue
	Dlq         *queue.DLQ
	Concurrency int

	wg sync.WaitGroup

	shutdowmOnce sync.Once
}

func NewWorker(q *queue.Queue, dlq *queue.DLQ, concurrency int) *Worker {
	if concurrency <= 0 {
		concurrency = 1
	}
	w := &Worker{
		Handlers:    make(map[string]func(models.Job) error),
		Queue:       q,
		Dlq:         dlq,
		Concurrency: concurrency,
	}
	return w
}

func (w *Worker) Register(jobType string, handler func(models.Job) error) {
	w.Handlers[jobType] = handler
}

func (w *Worker) proccessJobs(jobNumber int) {
	for {
		job, ok := w.Queue.Dequeue()
		if !ok {
			fmt.Println("Shutdown the processess, all jobs exit")
			return
		}
		handler, ok := w.Handlers[job.JobType]

		if !ok {
			fmt.Println("No handler registered for job type:", job.JobType)
			continue
		}

		err := handler(job)

		if err != nil {
			fmt.Println(job.Attempts+1, "Attempt failed....Error :", err)
			job.LastError = err.Error()

			job.Attempts++

			if job.Attempts >= job.MaxAttempts {
				fmt.Println("3 Attempts reached... Failed to process the Job")

				w.Dlq.Add(job)

				continue
			}
			delay := helper.CalculateBackoff(job.Attempts)
			fmt.Println("Retrying after:", delay)

			time.Sleep(delay)

			w.Queue.Enqueue(job)
			continue
		}
	}
}

func (w *Worker) Start() {
	for i := 0; i < w.Concurrency; i++ {
		w.wg.Add(1)

		go func(workerNo int) {
			defer w.wg.Done()
			w.proccessJobs(workerNo)

		}(i + 1)
	}
}

func (w *Worker) Shutdown() {
	w.shutdowmOnce.Do(func() {
		w.Queue.Shutdown()
	})

	w.wg.Wait()

}
