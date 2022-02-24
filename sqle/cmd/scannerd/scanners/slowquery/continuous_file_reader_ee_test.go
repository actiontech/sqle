//go:build enterprise
// +build enterprise

package slowquery

import (
	"io/ioutil"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestContinuousFileReader(t *testing.T) {
	f, err := ioutil.TempFile("", "slowquery-reader")
	assert.NoError(t, err)

	_, err = f.WriteString("line-0\n")
	assert.NoError(t, err)

	r, err := NewContinuousFileReader(f.Name(), &testLogger{})
	assert.NoError(t, err)

	lineCh := make(chan string, 10)
	go func() {
		for {
			l, err := r.NextLine()
			if err != nil {
				close(lineCh)
				return
			}
			lineCh <- l
		}
	}()

	// test read from start of file(it is different from partial line read)
	assert.Equal(t, "line-0", <-lineCh)

	f.WriteString("line-1\nline-2")
	assert.Equal(t, "line-1", <-lineCh)
	assert.Equal(t, "line-2", <-lineCh)

	f.WriteString("line-3\n")

	err = r.Close()
	assert.NoError(t, err)

	line, ok := <-lineCh
	assert.False(t, ok)
	assert.Empty(t, line)
}
