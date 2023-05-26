package utils

import "context"

type asyncFunc func() error

func AsyncCallTimeout(ctx context.Context, fn asyncFunc) error {
	if fn == nil {
		return nil
	}

	errChan := make(chan error)
	go func() {
		errChan <- fn()
	}()

	select {
	case <-ctx.Done():
		return ctx.Err()
	case err := <-errChan:
		return err
	}
}
