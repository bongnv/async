// Package async provides APIs to handle asynchronous tasks without using channels.
package async

import (
	"context"
)

// Future provides a mechanism to access the future result of asynchronous works.
type Future[T any] interface {
	// Get waits for the async work to be done and returns the result.
	// If the provided context is cancelled or its deadline passes,
	// the function will return the context error.
	Get(ctx context.Context) (T, error)

	// Done returns a channel that's closed when the work is done.
	// fut.Done is provided for use in select statements:
	//
	//	fut := async.Go(ctx, someWorkFn)
	//
	//	select {
	//	case <-fut.Done():
	//		return fut.Get(ctx)
	//	case <-ctx.Done():
	//		// manual hanndling when context is done
	//		return nil, ctx.Err()
	//	}
	Done() <-chan struct{}
}

// Go runs fn in a different goroutine and returns an instance of Future.
// That instance of Future can be used to access result of the asynchronous function.
//
// Example:
//
//	fut := Go(ctx, func(ctx context.Context) (MyStruct, error) {
//		// Doing some stuff
//		return MyStruct{}, nil
//	}}
//
//	// Do other stuff meanwhile then wait for the response
//	// Call Get to access the response and the error. Get is a blocking call and will wait for the async func to finish
//	resp, err := fut.Get(ctx)
//
// Check Future APIs for more detail.
func Go[T any](ctx context.Context, fn func(ctx context.Context) (T, error)) Future[T] {
	fut := &futureImpl[T]{
		doneCh: make(chan struct{}),
	}

	go func() {
		val, err := fn(ctx)
		fut.value = val
		fut.err = err
		close(fut.doneCh)
	}()

	return fut
}

// futureImpl is the an implementation of Feature.
type futureImpl[T any] struct {
	doneCh chan struct{}
	value  T
	err    error
}

func (f *futureImpl[T]) Done() <-chan struct{} {
	return f.doneCh
}

func (f *futureImpl[T]) Get(ctx context.Context) (resp T, err error) {
	select {
	case <-f.doneCh:
		resp = f.value
		err = f.err
		return
	case <-ctx.Done():
		err = ctx.Err()
		return
	}
}
