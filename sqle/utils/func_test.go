package utils

import (
	"context"
	"errors"
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func elapsedFunc(s int, e error) error {
	time.Sleep(time.Second * time.Duration(s))
	return e
}

var errTestMsg = errors.New("test err message")

func cancelFn(cancel context.CancelFunc, timeout int) {
	time.Sleep(time.Second * time.Duration(timeout))
	cancel()
}

func getTimeoutCtx(timeout int) (context.Context, context.CancelFunc) {
	return context.WithTimeout(context.Background(), time.Second*time.Duration(timeout))
}

func TestAsyncCallTimeout(t *testing.T) {
	cases := []struct {
		timeout          int
		elapsedFunc      int
		elapsedFuncError error
		cancelSleep      int

		expectedErr error
	}{
		{5, 1, nil, 10, nil},
		{1, 5, nil, 10, context.DeadlineExceeded},
		{10, 5, nil, 1, context.Canceled},
		{10, 1, errTestMsg, 5, errTestMsg},
		{10, 5, errTestMsg, 1, context.Canceled},
		{1, 5, errTestMsg, 10, context.DeadlineExceeded},
	}

	for i := range cases {
		c := cases[i]
		t.Run(fmt.Sprintf("case %v", i), func(t *testing.T) {
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
