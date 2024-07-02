package api

import (
	"context"

	"github.com/xhermitx/gitpulse-tracker/profiler-service/models"
)

type Cache interface {
	Set(ctx context.Context, profile models.TopCandidates) error
	Get(ctx context.Context, jobID uint) ([]models.TopCandidates, error)
}

type Queue interface {
	Publish(data any, queueName string) error
}
