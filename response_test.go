package response_test

import (
	"net/http"
	"testing"

	"github.com/go-carrot/response"
	"github.com/stretchr/testify/assert"
)

type DummyResult struct {
	Value1 string
	Value2 string
}

type StackWriter struct {
	HeaderInt int
	Stack     []string
}

func (sw *StackWriter) Write(p []byte) (n int, err error) {
	sw.Stack = append(sw.Stack, string(p))
	return 0, nil
}

func (sw *StackWriter) WriteHeader(header int) {
	sw.HeaderInt = header
}

func (sw *StackWriter) Header() http.Header {
	return make(http.Header, 0)
}

func (sw *StackWriter) Peek() string {
	if len(sw.Stack) > 0 {
		return sw.Stack[len(sw.Stack)-1]
	}
	return ""
}

func TestResponseNotSet(t *testing.T) {
	sw := new(StackWriter)
	resp := response.New(sw)

	resp.Output()
	assert.Equal(t, "{\"meta\":{\"success\":false,\"status_code\":500,\"status_text\":\"Internal Server Error\",\"error_details\":null},\"content\":null}", sw.Peek())
}

func TestResponseSingleDetail(t *testing.T) {
	sw := new(StackWriter)
	resp := response.New(sw)

	resp.SetErrorDetails("Missing Auth")
	resp.Output()
	assert.Equal(t, sw.HeaderInt, 500)
	assert.Equal(t, "{\"meta\":{\"success\":false,\"status_code\":500,\"status_text\":\"Internal Server Error\",\"error_details\":\"Missing Auth\"},\"content\":null}", sw.Peek())
}

func TestSuccessfulResult(t *testing.T) {
	sw := new(StackWriter)
	resp := response.New(sw)

	resp.SetResult(http.StatusOK,
		&DummyResult{
			Value1: "Hello World",
			Value2: "Wow",
		},
	)
	resp.Output()
	assert.Equal(t, sw.HeaderInt, 200)
	assert.Equal(t, "{\"meta\":{\"success\":true,\"status_code\":200,\"status_text\":\"OK\",\"error_details\":null},\"content\":{\"Value1\":\"Hello World\",\"Value2\":\"Wow\"}}", sw.Peek())
}

func TestJsonRenderFailure(t *testing.T) {
	defer func() {
		recover()
	}()
	sw := new(StackWriter)
	resp := response.New(sw)

	resp.SetResult(http.StatusOK, func() {})
	resp.Output()
	t.Error("JsonRenderer should fail with content that can not be serialized to JSON")
}
