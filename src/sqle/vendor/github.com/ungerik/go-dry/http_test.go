package dry

import (
	"compress/gzip"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestHTTPCompressHandlerFunc(t *testing.T) {
	for i := 0; i < 100; i++ {
		handlerFunc := HTTPCompressHandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			fmt.Fprintf(w, "hello world!")
		})

		request, err := http.NewRequest("GET", "/foobar", nil)
		if err != nil {
			t.Fatalf("http.NewRequest failed: %v", err)
		}
		request.Header.Set("Accept-Encoding", "gzip, deflate")
		responseWriter := httptest.NewRecorder()

		handlerFunc.ServeHTTP(responseWriter, request)

		t.Logf("responseWriter.Body = %v", responseWriter.Body.Bytes())

		reader, err := gzip.NewReader(responseWriter.Body)
		if err != nil {
			t.Fatalf("gzip.NewReader failed: %v", err)
		}

		readData, err := ioutil.ReadAll(reader)
		if err != nil {
			t.Fatalf("Reading from body failed: %v", err)
		}

		if string(readData) != "hello world!" {
			t.Fatalf("Body content: expected \"hello world!\", got %s instead.", string(readData))
		}
	}
}

func TestHTTPCompressHandler(t *testing.T) {
	for i := 0; i < 100; i++ {
		handler := &HTTPCompressHandler{&helloWorldHandler{}}

		request, err := http.NewRequest("GET", "/foobar", nil)
		if err != nil {
			t.Fatalf("http.NewRequest failed: %v", err)
		}
		request.Header.Set("Accept-Encoding", "gzip")
		responseWriter := httptest.NewRecorder()

		handler.ServeHTTP(responseWriter, request)

		t.Logf("responseWriter.Body = %v", responseWriter.Body.Bytes())

		reader, err := gzip.NewReader(responseWriter.Body)
		if err != nil {
			t.Fatalf("gzip.NewReader failed: %v", err)
		}

		readData, err := ioutil.ReadAll(reader)
		if err != nil {
			t.Fatalf("Reading from body failed: %v", err)
		}

		if string(readData) != "hallo welt." {
			t.Fatalf("Body content: expected \"hallo welt.\", got %s instead.", string(readData))
		}
	}
}

type helloWorldHandler struct{}

func (h *helloWorldHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "hallo welt.")
}
