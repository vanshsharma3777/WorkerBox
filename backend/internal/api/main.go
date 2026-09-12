package api

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

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

func TestHandler1(job models.Job) error {
	fmt.Println("test handler called")

	return errors.New("deliberate test failure")
}

func TestWorker(w http.ResponseWriter, r *http.Request) {

	q := queue.NewQueue()
	dlq := queue.NewDLQ()
	worker := worker.NewWorker(q, dlq)

	worker.Register("test_1", TestHandler1)

	q.Enqueue(models.Job{
		JobType:     "test_1",
		MaxAttempts: 3,
	})
	fmt.Println("Enqueue job done")

	worker.Start()
	fmt.Println("Starting worker")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"msg": "Worker started successfully",
	})
}
