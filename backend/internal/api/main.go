package api

import (
	"encoding/json"
	"net/http"

	"github.com/google/uuid"

	"github.com/vanshsharma3777/WorkerBox/internal/repository"
	"github.com/vanshsharma3777/WorkerBox/internal/worker"
	"github.com/vanshsharma3777/WorkerBox/models"
)

type API struct {
	Repo      *repository.JobRepo
	AppWorker *worker.Worker
}

func NewAPI(
	repo *repository.JobRepo,
	appWorker *worker.Worker,
) *API {
	return &API{
		Repo:      repo,
		AppWorker: appWorker,
	}
}

func (a *API) TestWorker(w http.ResponseWriter, r *http.Request) {
	concurrency := 3
	for i := 0; i < concurrency; i++ {

		job := models.Job{
			ID:          uuid.New().String(),
			JobType:     "email",
			Attempts:    0,
			MaxAttempts: 3,
			Status:      "queued",
			Payload:     "{}",
		}
		if err := a.Repo.CreateJob(&job); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		if err := a.AppWorker.Queue.Enqueue(job); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}

	json.NewEncoder(w).Encode(map[string]interface{}{
		"msg": "Jobs added successfully",
	})
}
