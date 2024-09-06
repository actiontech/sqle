//go:build enterprise
// +build enterprise

package scanners

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
	fileName       string
	l              Logger
	skipBlankLines bool

	tf *tail.Tail
}

func NewContinuousFileReader(fineName string, l Logger, skipBlankLines bool) (*ContinuousFileReader, error) {
	tf, err := tail.TailFile(fineName, tail.Config{Follow: true, ReOpen: true, CompleteLines: true})
	if err != nil {
		return nil, err
	}

	return &ContinuousFileReader{
		fileName:       fineName,
		tf:             tf,
		l:              l,
		skipBlankLines: skipBlankLines,
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
	// Empty lines cause SlowLogParser.Run() to panic
	if len(line.Text) == 0 && r.skipBlankLines {
		return r.NextLine()
	}
	return line.Text, nil
}

func (r *ContinuousFileReader) Close() error {
	return r.tf.Stop()
}

func (r *ContinuousFileReader) Metrics() *parser.ReaderMetrics { return nil }
