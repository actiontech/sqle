//go:build enterprise
// +build enterprise

package slowquery

import (
	"io"

	"github.com/nxadm/tail"
	"github.com/percona/pmm-agent/agents/mysql/slowlog/parser"
)

// ContinuousFileReader implements
// Reader(github.com/percona/pmm-agent/agents/mysql/parser) interface. It
// reads lines from the single file. When EOF is reached, it will wait for
// more data to become available(like tail -f).
type ContinuousFileReader struct {
	fileName string
	l        Logger

	tf *tail.Tail
}

func NewContinuousFileReader(fineName string, l Logger) (*ContinuousFileReader, error) {
	tf, err := tail.TailFile(fineName, tail.Config{Follow: true})
	if err != nil {
		return nil, err
	}

	return &ContinuousFileReader{
		fileName: fineName,
		tf:       tf,
		l:        l,
	}, nil
}

func (r *ContinuousFileReader) NextLine() (string, error) {
	line, ok := <-r.tf.Lines
	if !ok {
		err := r.tf.Err()
		if err == nil {
			err = io.EOF
		}
		return "", err
	}
	return line.Text, nil
}

func (r *ContinuousFileReader) Close() error {
	return r.tf.Stop()
}

func (r *ContinuousFileReader) Metrics() *parser.ReaderMetrics { return nil }
