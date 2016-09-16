package response_test

import (
	"encoding/json"
	"fmt"
	"github.com/dcstack/response"
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

func (suite *ResponseTestSuite) SetupTest() {
	myErrors := map[int]string{
		ErrorMissingAuth:      "Missing Auth",
		ErrorMissingParameter: "Missing Parameter",
	}
	response.SetErrorMap(myErrors)
}

func (suite *ResponseTestSuite) TestResponseNotSet() {
	resp := response.New()
	result := resp.Output()
	assert.Equal(suite.T(), "{\"success\":false,\"status_code\":500,\"status_text\":\"Internal Server Error\",\"error_details\":null,\"content\":null}", result)
}

func (suite *ResponseTestSuite) TestResponseSingleDetail() {
	resp := response.New()
	resp.AddErrorDetail(1)
	result := resp.Output()
	assert.Equal(suite.T(), "{\"success\":false,\"status_code\":500,\"status_text\":\"Internal Server Error\",\"error_details\":[{\"code\":1,\"text\":\"Missing Auth\"}],\"content\":null}", result)
}

func (suite *ResponseTestSuite) TestResponseMultipleDetails() {
	resp := response.New()
	resp.AddErrorDetail(ErrorMissingAuth, ErrorMissingParameter)
	result := resp.Output()
	assert.Equal(suite.T(), "{\"success\":false,\"status_code\":500,\"status_text\":\"Internal Server Error\",\"error_details\":[{\"code\":1,\"text\":\"Missing Auth\"},{\"code\":2,\"text\":\"Missing Parameter\"}],\"content\":null}", result)
}

func (suite *ResponseTestSuite) TestInvalidErrorDetailCode() {
	defer func() {
		recover()
	}()
	resp := response.New()
	resp.AddErrorDetail(3)
	suite.T().Error("AddErrorDetail should fail with an invalid error code")
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
	assert.Equal(suite.T(), "{\"success\":true,\"status_code\":200,\"status_text\":\"OK\",\"error_details\":null,\"content\":{\"Value1\":\"Hello World\",\"Value2\":\"Wow\"}}", result)
}

func (suite *ResponseTestSuite) TestCustomRenderer() {
	resp := response.New().SetRenderer(new(PrettyJsonRenderer))
	result := resp.Output()
	fmt.Print(result)
	expectedResult := `{
    "success": false,
    "status_code": 500,
    "status_text": "Internal Server Error",
    "error_details": null,
    "content": null
}`
	assert.Equal(suite.T(), expectedResult, result)
}

func (suite *ResponseTestSuite) TestJsonRenderFailure() {
	defer func() {
		recover()
	}()
	resp := response.New()
	resp.SetResult(http.StatusOK, make(map[int]int))
	resp.Output()
	suite.T().Error("JsonRenderer should fail with content that can not be serialized to JSON")
}

// In order for 'go test' to run this suite, we need to create
// a normal test function and pass our suite to suite.Run
func TestResponseTestSuite(t *testing.T) {
	suite.Run(t, new(ResponseTestSuite))
}
