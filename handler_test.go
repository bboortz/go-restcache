package restcache

import (
	"fmt"
	"strings"
	"testing"
	//	"reflect"
	"encoding/json"
	"io"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	//	"github.com/julienschmidt/httprouter"
	//	"github.com/davecgh/go-spew/spew"
	"github.com/stretchr/testify/assert"
)

type TestItemCreate struct {
	Name string `json:"Name"`
}

func genericHandlerApiTest(t *testing.T, method string, url string, expectedStatusCode int) []byte {
	return genericHandlerApiTestWithRequestBody(t, method, url, expectedStatusCode, nil)
}

func genericHandlerApiTestWithRequestBody(t *testing.T, method string, url string, expectedStatusCode int, requestBody io.Reader) []byte {
	assert := assert.New(t)
	router := NewRouter()

	req, err := http.NewRequest(method, url, requestBody)
	recorder := httptest.NewRecorder()
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

func TestHandlerIndexRead(t *testing.T) {
	assert := assert.New(t)
	body := genericHandlerApiTest(t, "GET", "/", 200)

	bodyResponse := Api{}
	if err := json.Unmarshal(body, &bodyResponse); err != nil {
		fmt.Println("ERROR: ", err)
	}
	assert.NotNil(bodyResponse)
	assert.NotEmpty(bodyResponse.ApiName)
	assert.NotEmpty(bodyResponse.ApiVersion)
}

func TestHandlerAliveRead(t *testing.T) {
	assert := assert.New(t)
	body := genericHandlerApiTest(t, "GET", "/alive", 200)

	bodyResponse := Alive{}
	if err := json.Unmarshal(body, &bodyResponse); err != nil {
		fmt.Println("ERROR: ", err)
	}
	assert.NotNil(bodyResponse)
	assert.True(bodyResponse.Alive)
}

func TestHandlerMethodNotAllowed(t *testing.T) {
	assert := assert.New(t)
	requestStruct := TestItemCreate{Name: "go-test"}
	requestJson, _ := json.Marshal(requestStruct)
	requestBody := string(requestJson)
	body := genericHandlerApiTestWithRequestBody(t, "POST", "/alive", 405, strings.NewReader(requestBody))

	bodyResponse := Alive{}
	if err := json.Unmarshal(body, &bodyResponse); err != nil {
		fmt.Println("ERROR: ", err)
	}
	assert.NotNil(bodyResponse)
	assert.Empty(bodyResponse.Alive)
}

func TestHandlerCreateNotFound(t *testing.T) {
	assert := assert.New(t)
	requestStruct := TestItemCreate{Name: "go-test"}
	requestJson, _ := json.Marshal(requestStruct)
	requestBody := string(requestJson)
	body := genericHandlerApiTestWithRequestBody(t, "POST", "/notfound", 404, strings.NewReader(requestBody))

	bodyResponse := Api{}
	if err := json.Unmarshal(body, &bodyResponse); err != nil {
		fmt.Println("ERROR: ", err)
	}
	assert.NotNil(bodyResponse)
	assert.Empty(bodyResponse.ApiName)
	assert.Empty(bodyResponse.ApiVersion)
}
