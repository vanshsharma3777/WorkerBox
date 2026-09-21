package helper

import (
	"fmt"
	"sync/atomic"
	"time"

	"github.com/vanshsharma3777/WorkerBox/models"
)

var running int32
var maxRunning int32

func TestHandler(job models.Job) error {
	current := atomic.AddInt32(&running, 1)

	fmt.Printf(
		"Job %v started | attempt: %d | running: %d\n",
		job.ID,
		job.Attempts+1,
		current,
	)

	// Simulate actual work
	time.Sleep(2 * time.Second)

	atomic.AddInt32(&running, -1)

	fmt.Printf(
		"Job %s finished\n",
		job.ID,
	)

	return nil
}
