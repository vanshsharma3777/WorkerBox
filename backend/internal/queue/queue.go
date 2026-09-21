package queue

import (
	"errors"
	"fmt"
	"sync"

	"github.com/vanshsharma3777/WorkerBox/models"
)

type Queue struct {
	mu       sync.Mutex
	cond     *sync.Cond
	jobs     []models.Job
	shutdown bool
}

func NewQueue() *Queue {
	q := &Queue{
		jobs: make([]models.Job, 0),
	}

	q.cond = sync.NewCond(&q.mu)

	return q
}

func (q *Queue) Enqueue(job models.Job) error {
	q.mu.Lock()
	defer q.mu.Unlock()

	if q.shutdown {
		return errors.New("queue is shutting down")
	}

	q.jobs = append(q.jobs, job)

	q.cond.Signal()

	return nil
}

func (q *Queue) Dequeue() (models.Job, bool) {
	q.mu.Lock()
	defer q.mu.Unlock()

	for len(q.jobs) == 0 && !q.shutdown {
		q.cond.Wait()
	}

	if q.shutdown && len(q.jobs) == 0 {
		return models.Job{}, false
	}

	job := q.jobs[0]

	q.jobs = q.jobs[1:]

	return job, true
}

func (q *Queue) Shutdown() {
	q.mu.Lock()
	defer q.mu.Unlock()

	q.shutdown = true
	q.cond.Broadcast()
	fmt.Println("Broadcast shutdown msg to all workers")
}
