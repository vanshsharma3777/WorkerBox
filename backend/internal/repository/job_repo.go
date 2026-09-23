package repository

import (
	"fmt"
	"log/slog"
	"time"

	"github.com/vanshsharma3777/WorkerBox/models"
	"gorm.io/gorm"
)

type JobRepo struct {
	DB *gorm.DB
}

func NewJobRepository(db *gorm.DB) *JobRepo {
	slog.Info("Job repository initialized")
	fmt.Println()
	fmt.Println()

	return &JobRepo{
		DB: db,
	}
}

func (r *JobRepo) CreateJob(job *models.Job) error {
	slog.Info(
		"creating job",
		"job_id", job.ID,
		"job_type", job.JobType,
		"status", job.Status,
	)

	err := r.DB.Create(job).Error

	if err != nil {
		slog.Error(
			"failed to create job",
			"job_id", job.ID,
			"error", err,
		)
		return err
	}

	slog.Info(
		"job created successfully",
		"job_id", job.ID,
	)
	fmt.Println()
	fmt.Println()

	return nil
}

func (r *JobRepo) UpdateJobStatus(id string, status models.JobStatus) error {
	slog.Info(
		"updating job status",
		"job_id", id,
		"status", status,
	)

	err := r.DB.
		Model(&models.Job{}).
		Where("id = ?", id).
		Updates(map[string]interface{}{
			"status":     status,
			"updated_at": time.Now(),
		}).Error

	if err != nil {
		slog.Error(
			"failed to update job status",
			"job_id", id,
			"status", status,
			"error", err,
		)
		return err
	}

	slog.Info(
		"job status updated",
		"job_id", id,
		"status", status,
	)

	fmt.Println()
	fmt.Println()

	return nil
}

func (r *JobRepo) GetRecoverableJobs() ([]models.Job, error) {
	slog.Info("recovering jobs from database")

	var jobs []models.Job

	err := r.DB.Where(
		"status IN ?",
		[]models.JobStatus{
			models.StatusProcessing,
			models.StatusQueued,
		},
	).Find(&jobs).Error

	if err != nil {
		slog.Error(
			"failed to recover jobs",
			"error", err,
		)
		return nil, err
	}

	slog.Info(
		"jobs recovered successfully",
		"count", len(jobs),
	)
	fmt.Println()
	fmt.Println()

	return jobs, nil
}
