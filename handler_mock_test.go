package restcache

import (
	"bytes"
	"fmt"
	"strings"
	"testing"
	//	"reflect"
	"io"
	"io/ioutil"
	"net/http"
	//	"net/http/httptest"
	"encoding/json"
	//	"github.com/julienschmidt/httprouter"
	//	"github.com/davecgh/go-spew/spew"
	"github.com/stretchr/testify/assert"
)

type mockResponseWriter struct {
	Code int
	//	Body	*bytes.Buffer
	Body io.Reader
	//	Body []byte
}

func (m *mockResponseWriter) Header() (h http.Header) {
	return http.Header{}
}

func (m *mockResponseWriter) Write(p []byte) (n int, err error) {
	m.Body = bytes.NewReader(p)
	return len(p), nil
}

func (m *mockResponseWriter) WriteString(s string) (n int, err error) {
	return len(s), nil
}

func (m *mockResponseWriter) WriteHeader(s int) { m.Code = s }

type handlerStruct struct {
	handeled *bool
}

func (h handlerStruct) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	*h.handeled = true
}

/*
func MockHandler(t *testing.T) {
	router := httprouter.New()

	routed := false
	router.Handle("GET", "/user/:name", func(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
		routed = true
		want := httprouter.Params{httprouter.Param{"name", "gopher"}}
		if !reflect.DeepEqual(ps, want) {
			t.Fatalf("wrong wildcard values: want %v, got %v", want, ps)
		}
	})

	w := new(mockResponseWriter)

	req, _ := http.NewRequest("GET", "/user/gopher", nil)
	router.ServeHTTP(w, req)

	if !routed {
		t.Fatal("routing failed")
	}
}
*/
func genericHandlerApiMock(t *testing.T, method string, url string, expectedStatusCode int) []byte {
	return genericHandlerApiMockWithRequestBody(t, method, url, expectedStatusCode, nil)
}

func genericHandlerApiMockWithRequestBody(t *testing.T, method string, url string, expectedStatusCode int, requestBody io.Reader) []byte {
	assert := assert.New(t)
	router := NewRouter()

	req, err := http.NewRequest(method, url, requestBody)
	recorder := new(mockResponseWriter)
	router.ServeHTTP(recorder, req)
	assert.Nil(err)
	assert.NotNil(req)
	assert.NotNil(recorder)
	assert.Equal(expectedStatusCode, recorder.Code)

	body, err := ioutil.ReadAll(io.LimitReader(recorder.Body, 1048576))
	if err != nil {
		panic(err)
	}
	assert.NotNil(body)

	return body
}

func TestHandlerIndexReadMock(t *testing.T) {
	assert := assert.New(t)
	body := genericHandlerApiMock(t, "GET", "/", 200)

	bodyResponse := Api{}
	if err := json.Unmarshal(body, &bodyResponse); err != nil {
		fmt.Println("ERROR: ", err)
	}
	assert.NotNil(bodyResponse)
	assert.NotEmpty(bodyResponse.ApiName)
	assert.NotEmpty(bodyResponse.ApiVersion)
}

func TestHandlerAliveReadMock(t *testing.T) {
	assert := assert.New(t)
	body := genericHandlerApiMock(t, "GET", "/alive", 200)

	bodyResponse := Alive{}
	if err := json.Unmarshal(body, &bodyResponse); err != nil {
		fmt.Println("ERROR: ", err)
	}
	assert.NotNil(bodyResponse)
	assert.True(bodyResponse.Alive)
}

func TestHandlerMethodNotAllowedMock(t *testing.T) {
	assert := assert.New(t)
	requestStruct := ServiceCreate{Name: "go-test"}
	requestJson, _ := json.Marshal(requestStruct)
	requestBody := string(requestJson)
	body := genericHandlerApiMockWithRequestBody(t, "POST", "/alive", 405, strings.NewReader(requestBody))

	bodyResponse := Service{}
	if err := json.Unmarshal(body, &bodyResponse); err != nil {
		fmt.Println("ERROR: ", err)
	}
	assert.NotNil(bodyResponse)
	assert.Empty(bodyResponse.Id)
	assert.Empty(bodyResponse.Name)
	assert.False(bodyResponse.Completed)
	assert.Empty(bodyResponse.Due)
}

func TestHandlerNotFoundMock(t *testing.T) {
	assert := assert.New(t)
	requestStruct := ServiceCreate{Name: "go-test"}
	requestJson, _ := json.Marshal(requestStruct)
	requestBody := string(requestJson)
	body := genericHandlerApiMockWithRequestBody(t, "POST", "/notfound", 404, strings.NewReader(requestBody))

	bodyResponse := Service{}
	if err := json.Unmarshal(body, &bodyResponse); err != nil {
		fmt.Println("ERROR: ", err)
	}
	assert.NotNil(bodyResponse)
	assert.Empty(bodyResponse.Id)
	assert.Empty(bodyResponse.Name)
	assert.False(bodyResponse.Completed)
	assert.Empty(bodyResponse.Due)
}
