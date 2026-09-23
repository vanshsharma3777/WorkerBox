package worker

import (
	"fmt"
	"sync"
	"time"

	"github.com/vanshsharma3777/WorkerBox/helper"
	"github.com/vanshsharma3777/WorkerBox/internal/queue"
	"github.com/vanshsharma3777/WorkerBox/internal/repository"
	"github.com/vanshsharma3777/WorkerBox/models"
)

type Worker struct {
	Handlers    map[string]func(models.Job) error
	Queue       *queue.Queue
	Dlq         *queue.DLQ
	Concurrency int
	Repo        *repository.JobRepo

	wg sync.WaitGroup

	shutdowmOnce sync.Once
}

func NewWorker(q *queue.Queue, dlq *queue.DLQ, concurrency int, repo *repository.JobRepo) *Worker {
	if concurrency <= 0 {
		concurrency = 1
	}
	w := &Worker{
		Handlers:    make(map[string]func(models.Job) error),
		Queue:       q,
		Dlq:         dlq,
		Concurrency: concurrency,
		Repo:        repo,
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

		error := w.Repo.UpdateJobStatus(
			job.ID,
			models.StatusProcessing,
		)
		handler, ok := w.Handlers[job.JobType]

		if error != nil {
			fmt.Println("Failed to update job status:", error)
			continue
		}
		if !ok {
			fmt.Println("No handler registered for job type:", job.JobType)
			w.Repo.UpdateJobStatus(
				job.ID,
				models.StatusFailed,
			)
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
		} else {
			w.Repo.UpdateJobStatus(
				job.ID,
				models.StatusCompleted,
			)
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

func (w *Worker) RecoverJobs() error {
	jobs, err := w.Repo.GetRecoverableJobs()

	if err != nil {
		return err
	}

	for _, job := range jobs {
		if err := w.Queue.Enqueue(job); err != nil {
			return err
		}
	}

	return nil
}
