package response

import (
	"net/http"
)

type Response struct {
	Renderer     Renderer      `json:"-"`
	Success      bool          `json:"success"`
	StatusCode   int           `json:"status_code"`
	StatusText   string        `json:"status_text"`
	ErrorDetails []ErrorDetail `json:"error_details"`
	Content      interface{}   `json:"content"`
}

// New instantiates a new Response struct
func New() *Response {
	r := new(Response)
	r.Renderer = new(JsonRenderer)
	r.Success = false
	r.StatusCode = http.StatusInternalServerError
	return r
}

// SetRenderer sets the renderer for this response
func (r *Response) SetRenderer(renderer Renderer) *Response {
	r.Renderer = renderer
	return r
}

// AddErrorDetail appends an error to the response via an Error Code.
func (r *Response) AddErrorDetail(codes ...int) *Response {
	for _, code := range codes {
		r.ErrorDetails = append(r.ErrorDetails, ErrorDetail{code, ErrorDetailText(code)})
	}
	return r
}

// SetResult sets the result status code and content.
func (r *Response) SetResult(httpStatusCode int, content interface{}) *Response {
	r.StatusCode = httpStatusCode
	r.Content = content
	return r
}

// Output sets the appropriate status and passes it to the Renderer to render
func (r *Response) Output() string {
	if r.StatusCode >= 200 && r.StatusCode < 300 {
		r.Success = true
	}
	r.StatusText = http.StatusText(r.StatusCode)
	return r.Renderer.Render(r)
}
