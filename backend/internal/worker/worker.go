package worker

import (
	"fmt"
	"time"

	"github.com/vanshsharma3777/WorkerBox/helper"
	"github.com/vanshsharma3777/WorkerBox/internal/queue"
	"github.com/vanshsharma3777/WorkerBox/models"
)

type Worker struct {
	handlers map[string]func(models.Job) error
	queue    *queue.Queue
	dlq      *queue.DLQ
}

func NewWorker(q *queue.Queue, dlq *queue.DLQ) *Worker {
	return &Worker{
		handlers: make(map[string]func(models.Job) error),
		queue:    q,
		dlq:      dlq,
	}
}

func (w *Worker) Register(jobType string, handler func(models.Job) error) {
	w.handlers[jobType] = handler
	fmt.Println("Handler registered for job type:", jobType)
}

func (w *Worker) Start() {
	for {
		job := w.queue.Dequeue()

		fmt.Println(job.Attempts+1, " Attempt...")

		handler, ok := w.handlers[job.JobType]

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

				w.dlq.Add(job)

				continue
			}
			delay := helper.CalculateBackoff(job.Attempts)
			fmt.Println("Retrying after:", delay)

			time.Sleep(delay)

			fmt.Println()

			w.queue.Enqueue(job)
			continue
		}

		fmt.Println("Job completed:", job.JobType)

	}
}
