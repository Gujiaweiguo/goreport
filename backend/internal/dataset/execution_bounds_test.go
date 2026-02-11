package dataset

import (
	"context"
	"testing"
	"time"
)

func TestWithDatasetQueryTimeout(t *testing.T) {
	ctx, cancel := withDatasetQueryTimeout(context.Background())
	defer cancel()

	deadline, ok := ctx.Deadline()
	if !ok {
		t.Fatalf("expected deadline to be set")
	}

	remaining := time.Until(deadline)
	if remaining <= 0 || remaining > datasetQueryTimeout {
		t.Fatalf("expected remaining timeout within (0, %s], got %s", datasetQueryTimeout, remaining)
	}
}

func TestWithDatasetPreviewTimeout(t *testing.T) {
	ctx, cancel := withDatasetPreviewTimeout(context.Background())
	defer cancel()

	deadline, ok := ctx.Deadline()
	if !ok {
		t.Fatalf("expected deadline to be set")
	}

	remaining := time.Until(deadline)
	if remaining <= 0 || remaining > datasetPreviewTimeout {
		t.Fatalf("expected remaining timeout within (0, %s], got %s", datasetPreviewTimeout, remaining)
	}
}
