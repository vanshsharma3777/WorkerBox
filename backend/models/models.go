package models

import "time"

type JobStatus string

const (
	StatusQueued     JobStatus = "queued"
	StatusProcessing JobStatus = "processing"
	StatusCompleted  JobStatus = "completed"
	StatusFailed     JobStatus = "failed"
	StatusDead       JobStatus = "dead"
)

type Job struct {
	ID string `gorm:"type:uuid;default:gen_random_uuid();primaryKey"`

	JobType string `gorm:"not null"`

	Payload string `gorm:"type:jsonb;not null"`

	Status      JobStatus `gorm:"not null;default:'queued'"`
	Attempts    int       `gorm:"not null;default:0"`
	MaxAttempts int       `gorm:"not null;default:3"`

	LastError string `gorm:"type:text"`

	CreatedAt time.Time `gorm:"not null;autoCreateTime"`
	UpdatedAt time.Time `gorm:"not null;autoUpdateTime"`

	AvailableAt time.Time `gorm:"not null;index"`
}
