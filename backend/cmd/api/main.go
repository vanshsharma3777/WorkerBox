package main

import (
	"fmt"

	"net/http"
	"os"
	"time"

	"github.com/joho/godotenv"
	"github.com/vanshsharma3777/WorkerBox/internal/api"
	"github.com/vanshsharma3777/WorkerBox/internal/db"
	"github.com/vanshsharma3777/WorkerBox/internal/queue"
	"github.com/vanshsharma3777/WorkerBox/internal/repository"
	"github.com/vanshsharma3777/WorkerBox/internal/worker"
	"github.com/vanshsharma3777/WorkerBox/models"
)

var repo *repository.JobRepo
var appWorker *worker.Worker

func main() {
	fmt.Println("Starting WorkerBox...")

	if err := godotenv.Load(); err != nil {
		fmt.Println("Error loading .env")
	}

	databaseURL := os.Getenv("DATABASE_URL")
	if databaseURL == "" {
		fmt.Println("DATABASE_URL is not set")
	}

	if err := db.Init(databaseURL); err != nil {
		fmt.Println(err)
	}

	repo = repository.NewJobRepository(db.DB)

	q := queue.NewQueue()
	dlq := queue.NewDLQ()
	workersCount := 2

	appWorker = worker.NewWorker(q, dlq, workersCount, repo)

	fmt.Println("Number of active workers = ", workersCount)

	appWorker.Register("email", func(job models.Job) error {
		fmt.Println("Executing:", job.ID)

		time.Sleep(2 * time.Second)

		return nil
	})

	jobs, err := repo.GetRecoverableJobs()
	if err != nil {
		fmt.Println("failed to recover jobs:", err)
	}

	for _, job := range jobs {
		if err := q.Enqueue(job); err != nil {
			fmt.Println("failed to enqueue recovered job:", err)
		}
	}
	fmt.Println(" old job enqueued")
	fmt.Println("Recovered jobs:", len(jobs))
	appWorker.Start()
	apiServer := api.NewAPI(repo, appWorker)

	mux := http.NewServeMux()

	mux.HandleFunc("GET /test-worker", apiServer.TestWorker)
	http.ListenAndServe(":8080", mux)

}
