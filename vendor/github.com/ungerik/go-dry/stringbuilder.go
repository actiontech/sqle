package dry

import (
	"bytes"
	"fmt"
	"io"
	"strconv"
)

type StringBuilder struct {
	buffer bytes.Buffer
}

func (s *StringBuilder) Write(strings ...string) *StringBuilder {
	for _, str := range strings {
		s.buffer.WriteString(str)
	}
	return s
}

func (s *StringBuilder) Printf(format string, args ...interface{}) *StringBuilder {
	fmt.Fprintf(&s.buffer, format, args...)
	return s
}

func (s *StringBuilder) Byte(value byte) *StringBuilder {
	s.buffer.WriteByte(value)
	return s
}

func (s *StringBuilder) WriteBytes(bytes []byte) *StringBuilder {
	s.buffer.Write(bytes)
	return s
}

func (s *StringBuilder) Int(value int) *StringBuilder {
	s.buffer.WriteString(strconv.Itoa(value))
	return s
}

func (s *StringBuilder) Uint(value uint) *StringBuilder {
	s.buffer.WriteString(strconv.FormatUint(uint64(value), 10))
	return s
}

func (s *StringBuilder) Float(value float64) *StringBuilder {
	s.buffer.WriteString(strconv.FormatFloat(value, 'f', -1, 64))
	return s
}

func (s *StringBuilder) Bool(value bool) *StringBuilder {
	s.buffer.WriteString(strconv.FormatBool(value))
	return s
}

func (s *StringBuilder) WriteTo(writer io.Writer) (n int64, err error) {
	return s.buffer.WriteTo(writer)
}

func (s *StringBuilder) Bytes() []byte {
	return s.buffer.Bytes()
}

func (s *StringBuilder) String() string {
	return s.buffer.String()
}
