package dry

import (
	"bytes"
	"encoding/json"
	"encoding/xml"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
)

type wrappedResponseWriter struct {
	io.Writer
	http.ResponseWriter
}

func (wrapped wrappedResponseWriter) Write(data []byte) (int, error) {
	return wrapped.Writer.Write(data)
}

// HTTPCompressHandlerFunc wraps a http.HandlerFunc so that the response gets
// gzip or deflate compressed if the Accept-Encoding header of the request allows it.
func HTTPCompressHandlerFunc(handlerFunc http.HandlerFunc) http.HandlerFunc {
	return func(response http.ResponseWriter, request *http.Request) {
		NewHTTPCompressHandlerFromFunc(handlerFunc).ServeHTTP(response, request)
	}
}

// HTTPCompressHandler wraps a http.Handler so that the response gets
// gzip or deflate compressed if the Accept-Encoding header of the request allows it.
type HTTPCompressHandler struct {
	http.Handler
}

func NewHTTPCompressHandler(handler http.Handler) *HTTPCompressHandler {
	return &HTTPCompressHandler{handler}
}

func NewHTTPCompressHandlerFromFunc(handler http.HandlerFunc) *HTTPCompressHandler {
	return &HTTPCompressHandler{handler}
}

func (h *HTTPCompressHandler) ServeHTTP(response http.ResponseWriter, request *http.Request) {
	accept := request.Header.Get("Accept-Encoding")
	if strings.Contains(accept, "gzip") {
		response.Header().Set("Content-Encoding", "gzip")
		writer := Gzip.GetWriter(response)
		defer Gzip.ReturnWriter(writer)
		response = wrappedResponseWriter{Writer: writer, ResponseWriter: response}
	} else if strings.Contains(accept, "deflate") {
		response.Header().Set("Content-Encoding", "deflate")
		writer := Deflate.GetWriter(response)
		defer Deflate.ReturnWriter(writer)
		response = wrappedResponseWriter{Writer: writer, ResponseWriter: response}
	}
	h.Handler.ServeHTTP(response, request)
}

// HTTPPostJSON marshalles data as JSON
// and sends it as HTTP POST request to url.
// If the response status code is not 200 OK,
// then the status is returned as an error.
func HTTPPostJSON(url string, data interface{}) error {
	b, err := json.MarshalIndent(data, "", "  ")
	if err != nil {
		return err
	}
	response, err := http.Post(url, "application/json", bytes.NewBuffer(b))
	if err == nil && (response.StatusCode < 200 || response.StatusCode > 299) {
		err = errors.New(response.Status)
	}
	return err
}

// HTTPPostXML marshalles data as XML
// and sends it as HTTP POST request to url.
// If the response status code is not 200 OK,
// then the status is returned as an error.
func HTTPPostXML(url string, data interface{}) error {
	b, err := xml.MarshalIndent(data, "", "  ")
	if err != nil {
		return err
	}
	response, err := http.Post(url, "application/xml", bytes.NewBuffer(b))
	if err == nil && (response.StatusCode < 200 || response.StatusCode > 299) {
		err = errors.New(response.Status)
	}
	return err
}

// HTTPDelete performs a HTTP DELETE request
func HTTPDelete(url string) (statusCode int, statusText string, err error) {
	request, err := http.NewRequest("DELETE", url, nil)
	if err != nil {
		return 0, "", err
	}
	response, err := http.DefaultClient.Do(request)
	if err != nil {
		return 0, "", err
	}
	return response.StatusCode, response.Status, nil
}

// HTTPPostForm performs a HTTP POST request with data as application/x-www-form-urlencoded
func HTTPPostForm(url string, data url.Values) (statusCode int, statusText string, err error) {
	request, err := http.NewRequest("POST", url, strings.NewReader(data.Encode()))
	if err != nil {
		return 0, "", err
	}
	request.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	response, err := http.DefaultClient.Do(request)
	if err != nil {
		return 0, "", err
	}
	return response.StatusCode, response.Status, nil
}

// HTTPPutForm performs a HTTP PUT request with data as application/x-www-form-urlencoded
func HTTPPutForm(url string, data url.Values) (statusCode int, statusText string, err error) {
	request, err := http.NewRequest("PUT", url, strings.NewReader(data.Encode()))
	if err != nil {
		return 0, "", err
	}
	request.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	response, err := http.DefaultClient.Do(request)
	if err != nil {
		return 0, "", err
	}
	return response.StatusCode, response.Status, nil
}

// HTTPRespondMarshalJSON marshals response as JSON to responseWriter, sets Content-Type to application/json
// and compresses the response if Content-Encoding from the request allows it.
func HTTPRespondMarshalJSON(response interface{}, responseWriter http.ResponseWriter, request *http.Request) (err error) {
	NewHTTPCompressHandlerFromFunc(func(responseWriter http.ResponseWriter, request *http.Request) {
		var data []byte
		if data, err = json.Marshal(response); err == nil {
			responseWriter.Header().Set("Content-Type", "application/json")
			_, err = responseWriter.Write(data)
		}
	}).ServeHTTP(responseWriter, request)
	return err
}

// HTTPRespondMarshalIndentJSON marshals response as JSON to responseWriter, sets Content-Type to application/json
// and compresses the response if Content-Encoding from the request allows it.
// The JSON will be marshalled indented according to json.MarshalIndent
func HTTPRespondMarshalIndentJSON(response interface{}, prefix, indent string, responseWriter http.ResponseWriter, request *http.Request) (err error) {
	NewHTTPCompressHandlerFromFunc(func(responseWriter http.ResponseWriter, request *http.Request) {
		var data []byte
		if data, err = json.MarshalIndent(response, prefix, indent); err == nil {
			responseWriter.Header().Set("Content-Type", "application/json")
			_, err = responseWriter.Write(data)
		}
	}).ServeHTTP(responseWriter, request)
	return err
}

// HTTPRespondMarshalXML marshals response as XML to responseWriter, sets Content-Type to application/xml
// and compresses the response if Content-Encoding from the request allows it.
// If rootElement is not empty, then an additional root element with this name will be wrapped around the content.
func HTTPRespondMarshalXML(response interface{}, rootElement string, responseWriter http.ResponseWriter, request *http.Request) (err error) {
	NewHTTPCompressHandlerFromFunc(func(responseWriter http.ResponseWriter, request *http.Request) {
		var data []byte
		if data, err = xml.Marshal(response); err == nil {
			responseWriter.Header().Set("Content-Type", "application/xml")
			if rootElement == "" {
				_, err = fmt.Fprintf(responseWriter, "%s%s", xml.Header, data)
			} else {
				_, err = fmt.Fprintf(responseWriter, "%s<%s>%s</%s>", xml.Header, rootElement, data, rootElement)
			}
		}
	}).ServeHTTP(responseWriter, request)
	return err
}

// HTTPRespondMarshalIndentXML marshals response as XML to responseWriter, sets Content-Type to application/xml
// and compresses the response if Content-Encoding from the request allows it.
// The XML will be marshalled indented according to xml.MarshalIndent.
// If rootElement is not empty, then an additional root element with this name will be wrapped around the content.
func HTTPRespondMarshalIndentXML(response interface{}, rootElement string, prefix, indent string, responseWriter http.ResponseWriter, request *http.Request) (err error) {
	NewHTTPCompressHandlerFromFunc(func(responseWriter http.ResponseWriter, request *http.Request) {
		var data []byte
		contentPrefix := prefix
		if rootElement != "" {
			contentPrefix += indent
		}
		if data, err = xml.MarshalIndent(response, contentPrefix, indent); err == nil {
			responseWriter.Header().Set("Content-Type", "application/xml")
			if rootElement == "" {
				_, err = fmt.Fprintf(responseWriter, "%s%s\n%s", prefix, xml.Header, data)
			} else {
				_, err = fmt.Fprintf(responseWriter, "%s%s%s<%s>\n%s\n%s</%s>", prefix, xml.Header, prefix, rootElement, data, prefix, rootElement)
			}
		}
	}).ServeHTTP(responseWriter, request)
	return err
}

// HTTPRespondText sets Content-Type to text/plain
// and compresses the response if Content-Encoding from the request allows it.
func HTTPRespondText(response string, responseWriter http.ResponseWriter, request *http.Request) (err error) {
	NewHTTPCompressHandlerFromFunc(func(responseWriter http.ResponseWriter, request *http.Request) {
		responseWriter.Header().Set("Content-Type", "text/plain")
		_, err = responseWriter.Write([]byte(response))
	}).ServeHTTP(responseWriter, request)
	return err
}

// HTTPUnmarshalRequestBodyJSON reads a http.Request body and unmarshals it as JSON to result.
func HTTPUnmarshalRequestBodyJSON(request *http.Request, result interface{}) error {
	defer request.Body.Close()
	body, err := ioutil.ReadAll(request.Body)
	if err != nil {
		return err
	}
	return json.Unmarshal(body, result)
}
