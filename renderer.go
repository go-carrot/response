package response

import (
	"encoding/json"
)

type Renderer interface {
	Render(*Response) string
}

type JsonRenderer int

func (r *JsonRenderer) Render(resp *Response) string {
	b, err := json.Marshal(resp)
	if err != nil {
		panic("Unable to json.Marshal our Response")
	}
	return string(b)
}
