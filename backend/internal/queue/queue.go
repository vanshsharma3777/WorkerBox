package queue

import (
	"sync"

	"github.com/vanshsharma3777/WorkerBox/models"
)

type Queue struct {
	mu   sync.Mutex
	cond *sync.Cond
	jobs []models.Job
}

func NewQueue() *Queue {
	q := &Queue{
		jobs: make([]models.Job, 0),
	}

	q.cond = sync.NewCond(&q.mu)

	return q
}

func (q *Queue) Enqueue(job models.Job) {
	q.mu.Lock()
	defer q.mu.Unlock()

	q.jobs = append(q.jobs, job)

	q.cond.Signal()
}

func (q *Queue) Dequeue() models.Job {
	q.mu.Lock()
	defer q.mu.Unlock()

	for len(q.jobs) == 0 {
		q.cond.Wait()
	}

	job := q.jobs[0]

	q.jobs = q.jobs[1:]

	return job
}
