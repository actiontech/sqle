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

func (self *StringBuilder) Write(strings ...string) *StringBuilder {
	for _, str := range strings {
		self.buffer.WriteString(str)
	}
	return self
}

func (self *StringBuilder) Printf(format string, args ...interface{}) *StringBuilder {
	fmt.Fprintf(&self.buffer, format, args...)
	return self
}

func (self *StringBuilder) Byte(value byte) *StringBuilder {
	self.buffer.WriteByte(value)
	return self
}

func (self *StringBuilder) WriteBytes(bytes []byte) *StringBuilder {
	self.buffer.Write(bytes)
	return self
}

func (self *StringBuilder) Int(value int) *StringBuilder {
	self.buffer.WriteString(strconv.Itoa(value))
	return self
}

func (self *StringBuilder) Uint(value uint) *StringBuilder {
	self.buffer.WriteString(strconv.FormatUint(uint64(value), 10))
	return self
}

func (self *StringBuilder) Float(value float64) *StringBuilder {
	self.buffer.WriteString(strconv.FormatFloat(value, 'f', -1, 64))
	return self
}

func (self *StringBuilder) Bool(value bool) *StringBuilder {
	self.buffer.WriteString(strconv.FormatBool(value))
	return self
}

func (self *StringBuilder) WriteTo(writer io.Writer) (n int64, err error) {
	return self.buffer.WriteTo(writer)
}

func (self *StringBuilder) Bytes() []byte {
	return self.buffer.Bytes()
}

func (self *StringBuilder) String() string {
	return self.buffer.String()
}
