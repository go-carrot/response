package response

import (
	"net/http"
)

type Meta struct {
	Success      bool   `json:"success"`
	StatusCode   int    `json:"status_code"`
	StatusText   string `json:"status_text"`
	ErrorDetails string `json:"error_details"`
}

type Response struct {
	Renderer Renderer    `json:"-"`
	Meta     Meta        `json:"meta"`
	Content  interface{} `json:"content"`
}

// New instantiates a new Response struct
func New() *Response {
	r := new(Response)
	r.Renderer = new(JsonRenderer)
	r.Meta = Meta{}
	r.Meta.Success = false
	r.Meta.StatusCode = http.StatusInternalServerError
	r.Meta.StatusText = http.StatusText(http.StatusInternalServerError)
	return r
}

// SetRenderer sets the renderer for this response
func (r *Response) SetRenderer(renderer Renderer) *Response {
	r.Renderer = renderer
	return r
}

// AddErrorDetail appends an error to the response via an Error Code.
func (r *Response) SetErrorDetails(errorDetails string) *Response {
	r.Meta.ErrorDetails = errorDetails
	return r
}

// SetResult sets the result status code and content.
func (r *Response) SetResult(httpStatusCode int, content interface{}) *Response {
	r.Meta.StatusCode = httpStatusCode
	r.Meta.StatusText = http.StatusText(httpStatusCode)
	r.Content = content
	return r
}

// SetResultWithStatusText sets the result status code, status text, and content.
func (r *Response) SetResultWithStatusText(httpStatusCode int, statusText string, content interface{}) *Response {
	r.Meta.StatusCode = httpStatusCode
	r.Meta.StatusText = statusText
	r.Content = content
	return r
}

// Output sets the appropriate status and passes it to the Renderer to render
func (r *Response) Output() string {
	if r.Meta.StatusCode >= 200 && r.Meta.StatusCode < 300 {
		r.Meta.Success = true
	}
	return r.Renderer.Render(r)
}
