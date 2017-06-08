package response_test

import (
	"encoding/json"
	"github.com/go-carrot/response"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"net/http"
	"testing"
)

const (
	ErrorMissingAuth      = 1
	ErrorMissingParameter = 2
)

type PrettyJsonRenderer int

func (r *PrettyJsonRenderer) Render(resp *response.Response) string {
	b, err := json.MarshalIndent(resp, "", "    ")
	if err != nil {
		panic("Unable to json.Marshal our Response")
	}
	return string(b)
}

type ResponseTestSuite struct {
	suite.Suite
}

func (suite *ResponseTestSuite) TestResponseNotSet() {
	resp := response.New()
	result := resp.Output()
	assert.Equal(suite.T(), "{\"meta\":{\"success\":false,\"status_code\":500,\"status_text\":\"Internal Server Error\",\"error_details\":\"\"},\"content\":null}", result)
}

func (suite *ResponseTestSuite) TestResponseSingleDetail() {
	resp := response.New()
	resp.SetErrorDetails("Missing Auth")
	result := resp.Output()
	assert.Equal(suite.T(), "{\"meta\":{\"success\":false,\"status_code\":500,\"status_text\":\"Internal Server Error\",\"error_details\":\"Missing Auth\"},\"content\":null}", result)
}

type DummyResult struct {
	Value1 string
	Value2 string
}

func (suite *ResponseTestSuite) TestSuccessfulResult() {
	resp := response.New()
	resp.SetResult(http.StatusOK,
		&DummyResult{
			Value1: "Hello World",
			Value2: "Wow",
		},
	)
	result := resp.Output()
	assert.Equal(suite.T(), "{\"meta\":{\"success\":true,\"status_code\":200,\"status_text\":\"OK\",\"error_details\":\"\"},\"content\":{\"Value1\":\"Hello World\",\"Value2\":\"Wow\"}}", result)
}

func (suite *ResponseTestSuite) TestSuccessfulResultWithStatusText() {
	resp := response.New()
	resp.SetResultWithStatusText(http.StatusOK, "Status OK",
		&DummyResult{
			Value1: "Hello World",
			Value2: "Wow",
		},
	)
	result := resp.Output()
	assert.Equal(suite.T(), "{\"meta\":{\"success\":true,\"status_code\":200,\"status_text\":\"Status OK\",\"error_details\":\"\"},\"content\":{\"Value1\":\"Hello World\",\"Value2\":\"Wow\"}}", result)
}

func (suite *ResponseTestSuite) TestCustomRenderer() {
	resp := response.New().SetRenderer(new(PrettyJsonRenderer))
	result := resp.Output()
	expectedResult := `{
    "meta": {
        "success": false,
        "status_code": 500,
        "status_text": "Internal Server Error",
        "error_details": ""
    },
    "content": null
}`
	assert.Equal(suite.T(), expectedResult, result)
}

func (suite *ResponseTestSuite) TestJsonRenderFailure() {
	defer func() {
		recover()
	}()
	resp := response.New()
	resp.SetResult(http.StatusOK, func() {})
	resp.Output()
	suite.T().Error("JsonRenderer should fail with content that can not be serialized to JSON")
}

// In order for 'go test' to run this suite, we need to create
// a normal test function and pass our suite to suite.Run
func TestResponseTestSuite(t *testing.T) {
	suite.Run(t, new(ResponseTestSuite))
}
