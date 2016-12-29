// HTTP Client which handles generic error routing and marshaling
package httpclient

import (
	"crypto/tls"
	"errors"
	"io/ioutil"
	"net/http"
	"strings"
	"sync"
	"time"
	"github.com/ContainX/go-utils/encoding"
	"github.com/ContainX/go-utils/logger"
)

var log = logger.GetLogger("httpclient")

type Response struct {
	// Status is the underlying HTTP Status code
	Status  int
	// Content is RAW content/body as a string
	Content string
	// Elapsed time
	Elapsed time.Duration
	// Error is the error captured or nil
	Error   error
}

type Request struct {
	// Http Method type
	method Method
	// Complete URL including params
	url string
	// Post data
	data string
	// Expected data type
	result interface{}
	// encoding type (optional : default JSON)
	encodingType encoding.EncoderType
}

type HttpClientConfig struct {
	sync.RWMutex
	// Http Basic Auth Username
	HttpUser string
	// Http Basic Auth Password
	HttpPass string
	// Request timeout
	RequestTimeout int
	// TLS Insecure Skip Verify
	TLSInsecureSkipVerify bool
}

type httpClient struct {
	config HttpClientConfig
	http   *http.Client
}

type HttpClient interface {
	// Get a resource from the specified url and unmarshal
	// the response into the result if it is not nil
	Get(url string, result interface{}) *Response

	// Put a data resource to the specified url and unmarshal
	// the response into the result if it is not nil
	Put(url string, data interface{}, result interface{}) *Response

	// Delete a resource from the specified url.  If data is not nil
	// then a body is submitted in the request.  If a result is
	// not nil then the response body will be unmarshalled into the
	// result if it is not nil
	Delete(url string, data interface{}, result interface{}) *Response

	// Post the data against the specified url and unmarshal the
	// response into the result if it is not nil
	Post(url string, data interface{}, result interface{})
}

var (
	// invalid or error response
	ErrorInvalidResponse = errors.New("Invalid response from Remote")
	// some resource does not exists
	ErrorNotFound = errors.New("The resource does not exist")
	// Generic Error Message
	ErrorMessage = errors.New("Unknown error message was captured")
	// Not Authorized - 403
	ErrorNotAuthorized = errors.New("Not Authorized to perform this action - Status: 403")
	// Not Authenticated 401
	ErrorNotAuthenticated = errors.New("Not Authenticated to perform this action - Status: 401")

	// singleton client used for static function based calls
	sclient = DefaultHttpClient()
)

// NewDefaultConfig creates a HttpClientConfig wth default options
func NewDefaultConfig() *HttpClientConfig {
	return &HttpClientConfig{RequestTimeout: 30, TLSInsecureSkipVerify: false}
}

// DefaultHttpClient provides a basic default http client
func DefaultHttpClient() *httpClient {
	return NewHttpClient(*NewDefaultConfig())
}

func NewHttpClient(config HttpClientConfig) *httpClient {
	hc := &httpClient{
		config: config,
		http: &http.Client{
			Timeout: (time.Duration(config.RequestTimeout) * time.Second),
		},
	}
	if config.TLSInsecureSkipVerify {
		tr := &http.Transport{
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
		}
		hc.http.Transport = tr
	}
	return hc
}

func NewResponse(status int, elapsed time.Duration, content string, err error) *Response {
	return &Response{Status: status, Elapsed: elapsed, Content: content, Error: err}
}

// Get is a non-instance based call which uses the default configuration
// for custom overrides and control you should use the HttpClient.
//
// Usage: Get a resource from the specified url and unmarshal
// the response into the result if it is not nil
func Get(url string, result interface{}) *Response {
	return sclient.Get(url, result)
}

// Put is a non-instance based call which uses the default configuration.
// For custom overrides and control you should use the HttpClient.
//
// Usage: Put a data resource to the specified url and unmarshal
// the response into the result if it is not nil
func Put(url string, data interface{}, result interface{}) *Response {
	return sclient.Put(url, data, result)
}

// Delete is a non-instance based call which uses the default configuration.
// For custom overrides and control you should use the HttpClient.
//
// Usage: Delete a resource from the specified url.  If data is not nil
// then a body is submitted in the request.  If a result is
// not nil then the response body will be unmarshalled into the
// result if it is not nil
func Delete(url string, data interface{}, result interface{}) *Response {
	return sclient.Delete(url, data, result)
}

// Post is a non-instance based call which uses the default configuration.
// For custom overrides and control you should use the HttpClient.
//
// Usage: Post the data against the specified url and unmarshal the
// response into the result if it is not nil
func Post(url string, data interface{}, result interface{}) *Response {
	return sclient.Post(url, data, result)
}

func (h *httpClient) Get(url string, result interface{}) *Response {
	return h.invoke(&Request{method: GET, url: url, result: result})
}

func (h *httpClient) Put(url string, data interface{}, result interface{}) *Response {
	return h.httpCall(PUT, url, data, result)
}

func (h *httpClient) Delete(url string, data interface{}, result interface{}) *Response {
	return h.httpCall(DELETE, url, data, result)
}

func (h *httpClient) Post(url string, data interface{}, result interface{}) *Response {
	return h.httpCall(POST, url, data, result)
}

func (h *httpClient) httpCall(method Method, url string, data interface{}, result interface{}) *Response {
	var body string
	if data != nil {
		body = h.convertBody(data)
	}
	return h.invoke(&Request{method: method, url: url, data: body, result: result})
}

func (h *httpClient) invoke(r *Request) *Response {

	log.Debugf("%s - %s, Body:\n%s", r.method.String(), r.url, r.data)

	request, err := http.NewRequest(r.method.String(), r.url, strings.NewReader(r.data))

	if err != nil {
		return &Response{Error: err}
	}

	addHeaders(request)
	addAuthentication(h.config, request)

	req_start := time.Now()
	response, err := h.http.Do(request)
	req_elapsed := time.Now().Sub(req_start)

	if err != nil {
		return NewResponse(0, req_elapsed, "", err)
	}

	status := response.StatusCode
	var content string
	if response.ContentLength != 0 {
		defer response.Body.Close()
		rc, err := ioutil.ReadAll(response.Body)
		if err != nil {
			return NewResponse(status, req_elapsed, "", err)
		}
		content = string(rc)
	}

	log.Debugf("Status: %v, RAW: %s", status, content)

	if status >= 200 && status < 300 {
		if r.result != nil {
			h.convert(r, content)
		}
		return NewResponse(status, req_elapsed, content, nil)
	}

	switch status {
	case 500:
		return NewResponse(status, req_elapsed, content, ErrorInvalidResponse)
	case 404:
		return NewResponse(status, req_elapsed, content, ErrorNotFound)
	case 403:
		return NewResponse(status, req_elapsed, content, ErrorNotAuthorized)
	case 401:
		return NewResponse(status, req_elapsed, content, ErrorNotAuthenticated)
	}

	return NewResponse(status, req_elapsed, content, ErrorMessage)
}

func (h *httpClient) convertBody(data interface{}) string {
	if data == nil {
		return ""
	}
	encoder, _ := encoding.NewEncoder(encoding.JSON)
	body, _ := encoder.Marshal(data)
	return body
}

func (h *httpClient) convert(r *Request, content string) error {
	um, _ := encoding.NewEncoder(encoding.JSON)
	if r.encodingType != 0 {
		um, _ = encoding.NewEncoder(r.encodingType)
	}
	um.UnMarshalStr(content, r.result)
	return nil

}

func addHeaders(req *http.Request) {
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Accept", "application/json")
}

func addAuthentication(c HttpClientConfig, req *http.Request) {
	if c.HttpUser != "" {
		req.SetBasicAuth(c.HttpUser, c.HttpPass)
	}
}
