package util

import (
	"bytes"
	"compress/gzip"
	"fmt"
	"io/ioutil"
)

func GzipString(a string) ([]byte, error) {
	var b bytes.Buffer
	{
		gz, err := gzip.NewWriterLevel(&b, gzip.BestCompression)
		if nil != err {
			return nil, err
		}
		if _, err := gz.Write([]byte(a)); nil != err {
			return nil, err
		}
		if err := gz.Flush(); nil != err {
			return nil, err
		}
		if err := gz.Close(); nil != err {
			return nil, err
		}
	}
	fmt.Printf("%v:%v\n", len(a), len(b.Bytes()))
	return b.Bytes(), nil
}

func MustGzipString(a string) []byte {
	bs, err := GzipString(a)
	if nil != err {
		panic(err.Error())
	}
	return bs
}

func UnGzipString(input []byte) (string, error) {
	reader, err := gzip.NewReader(bytes.NewReader(input))
	if nil != err {
		return "", err
	}
	bs, err := ioutil.ReadAll(reader)
	if nil != err {
		return "", err
	}
	return string(bs), nil
}

func MustUnGzipString(input []byte) string {
	s, err := UnGzipString(input)
	if nil != err {
		panic(err.Error())
	}
	return s
}
