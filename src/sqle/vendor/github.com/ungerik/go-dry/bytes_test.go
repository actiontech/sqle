package dry

import (
	"bytes"
	"io"
	"reflect"
	"strings"
	"testing"
)

type MyString struct {
	str string
}

func (t MyString) String() string {
	return t.str
}

type MyError struct {
	str string
}

func (t MyError) Error() string {
	return t.str
}

func Test_BytesReader(t *testing.T) {
	expected := []byte("hello")
	testBytesReaderFn := func(input interface{}) {
		result := make([]byte, 5)
		returnedIoReader := BytesReader(input)
		n, _ := returnedIoReader.Read(result)
		if n != 5 {
			t.FailNow()
		}
		for i, _ := range result {
			if result[i] != expected[i] {
				t.FailNow()
			}
		}
		n, err := returnedIoReader.Read(result)
		if n != 0 || err != io.EOF {
			t.FailNow()
		}
	}

	testBytesReaderFn(strings.NewReader("hello"))

	bytesInput := []byte("hello")
	testBytesReaderFn(bytesInput)

	testBytesReaderFn("hello")

	myStr := MyString{"hello"}
	testBytesReaderFn(myStr)

	myErr := MyError{"hello"}
	testBytesReaderFn(myErr)
}

func assertStringsEqual(t *testing.T, str1, str2 string) {
	if str1 != str2 {
		t.FailNow()
	}
}

func Test_BytesMD5(t *testing.T) {
	assertStringsEqual(
		t, "5d41402abc4b2a76b9719d911017c592",
		BytesMD5("hello"))
}

func Test_BytesEncodeBase64(t *testing.T) {
	assertStringsEqual(
		t, "aGVsbG8=",
		BytesEncodeBase64("hello"))
}

func Test_BytesDecodeBase64(t *testing.T) {
	assertStringsEqual(
		t, "hello",
		BytesDecodeBase64("aGVsbG8="))
}

func Test_BytesEncodeHex(t *testing.T) {
	assertStringsEqual(
		t, "68656c6c6f",
		BytesEncodeHex("hello"))
}

func Test_BytesDecodeHex(t *testing.T) {
	assertStringsEqual(
		t, "hello",
		BytesDecodeHex("68656C6C6F"))
}

func testCompressDecompress(t *testing.T,
	compressFunc func([]byte) []byte,
	decompressFunc func([]byte) []byte) {
	testFn := func(testData []byte) {
		compressedData := compressFunc(testData)
		uncompressedData := decompressFunc(compressedData)
		if !bytes.Equal(testData, uncompressedData) {
			t.FailNow()
		}
	}

	go testFn([]byte("hello123"))
	go testFn([]byte("gopher456"))
	go testFn([]byte("dry789"))
}

func Test_BytesDeflateInflate(t *testing.T) {
	testCompressDecompress(t, BytesDeflate, BytesInflate)
}

func Test_BytesGzipUnGzip(t *testing.T) {
	testCompressDecompress(t, BytesGzip, BytesUnGzip)
}

func bytesHeadTailTestHelper(
	t *testing.T,
	testMethod func([]byte, int) ([]string, []byte),
	lines []byte, n int,
	expected_lines []string, expected_rest []byte) {
	result_lines, result_rest := testMethod(lines, n)
	if !reflect.DeepEqual(result_lines, expected_lines) {
		t.FailNow()
	}
	if !bytes.Equal(result_rest, expected_rest) {
		t.FailNow()
	}
}

func Test_BytesHead(t *testing.T) {
	bytesHeadTailTestHelper(
		t, BytesHead,
		[]byte("line1\nline2\r\nline3\nline4\nline5"), 3,
		[]string{"line1", "line2", "line3"}, []byte("line4\nline5"))
	bytesHeadTailTestHelper(
		t, BytesHead,
		[]byte("line1\nline2\r\nline3\nline4\nline5"), 6,
		[]string{"line1", "line2", "line3", "line4", "line5"}, []byte(""))
}

func Test_BytesTail(t *testing.T) {
	bytesHeadTailTestHelper(
		t, BytesTail,
		[]byte("line1\nline2\nline3\nline4\r\nline5"), 2,
		[]string{"line5", "line4"}, []byte("line1\nline2\nline3"))
	bytesHeadTailTestHelper(
		t, BytesTail,
		[]byte("line1\nline2\r\nline3\nline4\nline5"), 6,
		[]string{"line5", "line4", "line3", "line2", "line1"}, []byte(""))
}

func Test_BytesMap(t *testing.T) {
	upper := func(b byte) byte {
		return b - ('a' - 'A')
	}
	result := BytesMap(upper, []byte("hello"))
	correct := []byte("HELLO")
	if len(result) != len(correct) {
		t.Fail()
	}
	for i, _ := range result {
		if result[i] != correct[i] {
			t.Fail()
		}
	}
}

func Test_BytesFilter(t *testing.T) {
	azFunc := func(b byte) bool {
		return b >= 'A' && b <= 'Z'
	}
	result := BytesFilter(azFunc, []byte{1, 2, 3, 'A', 'f', 'R', 123})
	correct := []byte{'A', 'R'}
	if len(result) != len(correct) {
		t.Fail()
	}
	for i, _ := range result {
		if result[i] != correct[i] {
			t.Fail()
		}
	}
}
