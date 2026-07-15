package queue

import (
	"context"
)

type Publisher interface {
	PublishImageJob(ctx context.Context, job ImageJob) error
}
