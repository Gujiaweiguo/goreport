package dataset

import (
	"context"
	"time"
)

const (
	datasetQueryTimeout   = 15 * time.Second
	datasetPreviewTimeout = 8 * time.Second
)

func withDatasetQueryTimeout(ctx context.Context) (context.Context, context.CancelFunc) {
	return context.WithTimeout(ctx, datasetQueryTimeout)
}

func withDatasetPreviewTimeout(ctx context.Context) (context.Context, context.CancelFunc) {
	return context.WithTimeout(ctx, datasetPreviewTimeout)
}
