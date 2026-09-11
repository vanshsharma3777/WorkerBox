package api

import (
	"encoding/json"
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
	fmt.Println("worker registered successfully")
	return nil
}

func TestWorker(w http.ResponseWriter, r *http.Request) {

	fmt.Println("Test worker route hit")
	q := queue.NewQueue()

	worker := worker.NewWorker(q)

	worker.Register("test_2", TestHandler1)

	q.Enqueue(models.Job{
		JobType: "test_2",
	})
	fmt.Println("Enqueue job done")

	worker.Start()
	fmt.Println("Starting worker")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"msg": "Worker started successfully",
	})
}
