package retry

import (
	"errors"
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

type counter struct {
	expectedRetryCount int
	actualRetroCount   int
}

func newCounter(expectedRetryCount int) *counter {
	return &counter{
		expectedRetryCount: expectedRetryCount - 1,
		actualRetroCount:   0,
	}

}

var (
	notDoneErr = errors.New("NOT DONE!")
)

func TestRetryMultipleCases(t *testing.T) {

	var testCases = []struct {
		expectedRetryCount int
		isFuncResErr       bool
		doneChanErr        error
	}{
		{1, false, nil},
		{2, false, nil},
		{3, false, nil},
		{4, true, notDoneErr},
		{5, true, notDoneErr},
		{6, true, notDoneErr},
	}

	for idx := range testCases {
		testCase := testCases[idx]
		t.Run(fmt.Sprint(idx), func(t *testing.T) {
			_ = RetryDo(t, testCase.expectedRetryCount, testCase.isFuncResErr)
		})
	}
}

func RetryDo(t *testing.T, expectedRetryCount int, isFuncResErr bool) (doneErr error) {

	c := newCounter(expectedRetryCount)

	var fn RetryableFunc
	{
		fn = func() error {
			t.Logf("expected retry count: [%v], actual retry count: [%v]",
				c.expectedRetryCount, c.actualRetroCount)

			if c.expectedRetryCount == c.actualRetroCount {
				return nil
			}
			c.actualRetroCount++
			return errors.New("not yet")
		}
	}

	doneChan := make(chan struct{})
	{
		err := Do(fn, doneChan)
		assert.Equal(t, err != nil, isFuncResErr)
	}

	// check done chan
	{
		tick := time.Tick(time.Millisecond * 1)
		// check doneChan
		for {
			select {
			case <-doneChan:
				return nil
			case <-tick:
				return notDoneErr
			default:
			}
		}
	}
}
