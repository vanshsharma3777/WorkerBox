package worker

import (
	"fmt"

	"github.com/vanshsharma3777/WorkerBox/internal/queue"
	"github.com/vanshsharma3777/WorkerBox/models"
)

type Worker struct {
	handlers map[string]func(models.Job) error
	queue    *queue.Queue
}

func NewWorker(q *queue.Queue) *Worker {
	return &Worker{
		handlers: make(map[string]func(models.Job) error),
		queue:    q,
	}
}

func (w *Worker) Register(jobType string, handler func(models.Job) error) {
	w.handlers[jobType] = handler
	fmt.Println("Handler registered for job type:", jobType)
}

func (w *Worker) Start() {
	for {
		job := w.queue.Dequeue()

		handler, ok := w.handlers[job.JobType]

		if !ok {
			fmt.Println("No handler registered for job type:", job.JobType)
			continue
		}

		err := handler(job)

		if err != nil {
			fmt.Println("Job failed:", err)
			continue
		}

		fmt.Println("Job completed:", job.JobType)

	}
}
