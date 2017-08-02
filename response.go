package response

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

type Meta struct {
	Success      bool    `json:"success"`
	StatusCode   int     `json:"status_code"`
	StatusText   string  `json:"status_text"`
	ErrorDetails *string `json:"error_details"`
}

type Response struct {
	Writer  io.Writer   `json:"-"`
	Meta    Meta        `json:"meta"`
	Content interface{} `json:"content"`
}

// New instantiates a new Response struct
func New(writer io.Writer) *Response {
	r := new(Response)
	r.Writer = writer
	r.Meta = Meta{}
	r.Meta.Success = false
	r.Meta.StatusCode = http.StatusInternalServerError
	r.Meta.StatusText = http.StatusText(http.StatusInternalServerError)
	return r
}

// SetErrorDetails appends an error to the response via an Error Code.
func (r *Response) SetErrorDetails(errorDetails string) *Response {
	r.Meta.ErrorDetails = &errorDetails
	return r
}

// SetResult sets the result status code and content.
func (r *Response) SetResult(httpStatusCode int, content interface{}) *Response {
	r.Meta.StatusCode = httpStatusCode
	r.Meta.StatusText = http.StatusText(httpStatusCode)
	r.Content = content
	return r
}

// Output sets the appropriate status and writes the response JSON to the writer
func (r *Response) Output() {
	// Determine success
	if r.Meta.StatusCode >= 200 && r.Meta.StatusCode < 300 {
		r.Meta.Success = true
	}

	// Write header, if this is a ResponseWriter
	switch v := r.Writer.(type) {
	case http.ResponseWriter:
		v.WriteHeader(r.Meta.StatusCode)
	}

	// Write Body
	b, err := json.Marshal(r)
	if err != nil {
		panic("Unable to json.Marshal our Response")
	}
	fmt.Fprint(r.Writer, string(b))
}
