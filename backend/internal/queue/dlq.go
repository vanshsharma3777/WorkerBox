package queue

import (
	"fmt"
	"sync"

	"github.com/vanshsharma3777/WorkerBox/models"
)

type DLQ struct {
	jobs []models.Job
	mu   sync.Mutex
}

func NewDLQ() *DLQ {
	return &DLQ{
		jobs: make([]models.Job, 0),
	}
}

func (d *DLQ) Add(job models.Job) {
	d.mu.Lock()
	defer d.mu.Unlock()

	d.jobs = append(d.jobs, job)
	fmt.Println("Job added in DLQ JOBTYPE", job.JobType)
	fmt.Println("Last Error: ", job.LastError)
}
