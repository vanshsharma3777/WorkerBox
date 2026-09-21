package api

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/vanshsharma3777/WorkerBox/helper"
	"github.com/vanshsharma3777/WorkerBox/internal/queue"
	"github.com/vanshsharma3777/WorkerBox/internal/worker"
	"github.com/vanshsharma3777/WorkerBox/models"
)

func Test(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Server running")

	w.Header().Set("Content-Type", "application/json")

	err := json.NewEncoder(w).Encode(map[string]string{
		"msg": "Server running.",
	})
	if err != nil {
		json.NewEncoder(w).Encode("Internal Server Error in /test")
	}
}

func TestWorker(w http.ResponseWriter, r *http.Request) {

	q := queue.NewQueue()
	dlq := queue.NewDLQ()
	concurrency := 3
	worker := worker.NewWorker(q, dlq, concurrency)

	worker.Register("test", helper.TestHandler)

	for i := 0; i < 10; i++ {
		job := models.Job{
			ID:          uuid.New().String(),
			JobType:     "test",
			Attempts:    0,
			MaxAttempts: 3,
		}

		q.Enqueue(job)
	}

	worker.Start()

	time.Sleep(1 * time.Second)

	fmt.Println(">>> SHUTDOWN REQUESTED")

	worker.Shutdown()

	fmt.Println(">>> SHUTDOWN COMPLETE")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"msg": "Worker started successfully",
	})
}
