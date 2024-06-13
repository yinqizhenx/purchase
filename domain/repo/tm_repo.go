package repo

import (
	"context"

	"purchase/domain/entity/tm"
)

type TMRepo interface {
	GetTasksStats(context.Context, string, []string) ([]*tm.TMState, error)
}
