package async_test

import (
	"context"
	"testing"
	"time"

	"github.com/bongnv/async"
)

func TestAsync_Get(t *testing.T) {
	t.Run("should return an error when context is cancelled", func(t *testing.T) {
		ctx, cancel := context.WithCancel(context.Background())
		cancel()

		fut := async.Go(ctx, func(ctx context.Context) (int, error) {
			return 1, nil
		})

		_, err := fut.Get(ctx)

		if err != context.Canceled {
			t.Fatalf("Expected %v, but got %v", ctx.Err(), err)
		}
	})

	t.Run("should return a response when there is no error", func(t *testing.T) {
		fut := async.Go(context.Background(), func(ctx context.Context) (int, error) {
			return 1, nil
		})

		resp, err := fut.Get(context.Background())
		if err != nil {
			t.Fatalf("Expected no error, but got %v", err)
		}

		if resp != 1 {
			t.Fatalf("Expected a response, but got %v", resp)
		}
	})

	t.Run("should return an error when context's deadline passes", func(t *testing.T) {
		testEndCh := make(chan struct{})
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Millisecond)
		defer cancel()

		fut := async.Go(context.Background(), func(ctx context.Context) (int, error) {
			<-testEndCh
			return 1, nil
		})

		_, err := fut.Get(ctx)
		if err != context.DeadlineExceeded {
			t.Fatalf("Expected %v, but got %v", ctx.Err(), err)
		}
	})
}

func TestAsync_Done(t *testing.T) {
	t.Run("should return a response when there is no error", func(t *testing.T) {
		fut := async.Go(context.Background(), func(ctx context.Context) (int, error) {
			return 1, nil
		})

		select {
		case <-fut.Done():
			resp, err := fut.Get(context.Background())
			if err != nil {
				t.Fatalf("Expected no error, but got %v", err)
			}

			if resp != 1 {
				t.Fatalf("Expected a response, but got %v", resp)
			}
		case <-time.After(100 * time.Millisecond):
			t.Fatal("test timed out")
		}
	})
}
