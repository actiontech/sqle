package utils

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func elapsedFunc(s int, e error) error {
	time.Sleep(time.Millisecond * time.Duration(s))
	return e
}

var errTestMsg = errors.New("test err message")

func cancelFn(cancel context.CancelFunc, timeout int) {
	time.Sleep(time.Millisecond * time.Duration(timeout))
	cancel()
}

func getTimeoutCtx(timeout int) (context.Context, context.CancelFunc) {
	return context.WithTimeout(context.Background(), time.Millisecond*time.Duration(timeout))
}

func TestAsyncCallTimeout(t *testing.T) {
	cases := []struct {
		timeout          int
		elapsedFunc      int
		elapsedFuncError error
		cancelSleep      int

		expectedErr error
	}{
		{3, 2, nil, 4, nil},
		{2, 3, nil, 4, context.DeadlineExceeded},
		{4, 3, nil, 2, context.Canceled},
		{4, 3, errTestMsg, 4, errTestMsg},
		{4, 3, errTestMsg, 2, context.Canceled},
		{2, 3, errTestMsg, 4, context.DeadlineExceeded},
	}

	for i := range cases {
		c := cases[i]
		t.Run("", func(t *testing.T) {
			ctx, cancel := getTimeoutCtx(c.timeout)
			go cancelFn(cancel, c.cancelSleep)
			err := AsyncCallTimeout(ctx, func() error {
				return elapsedFunc(c.elapsedFunc, c.elapsedFuncError)
			})
			assert.Equal(t, c.expectedErr, err)
			assert.True(t, errors.Is(c.expectedErr, err))
		})
	}
}
