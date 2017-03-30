<a href="https://engineering.carrot.is/"><p align="center"><img src="https://cloud.githubusercontent.com/assets/2105067/24525319/d3d26516-1567-11e7-9506-7611b3287d53.png" alt="Go Carrot" width="350px" align="center;" /></p></a>
# Response

[![Build Status](https://travis-ci.org/go-carrot/response.svg?branch=master)](https://travis-ci.org/go-carrot/response) [![codecov](https://codecov.io/gh/go-carrot/response/branch/master/graph/badge.svg)](https://codecov.io/gh/go-carrot/response)

Response is a library that flexibly generates output that follows [carrot/restful-api-spec](https://github.com/carrot/restful-api-spec).

## Output Interface

Following [carrot/restful-api-spec](https://github.com/carrot/restful-api-spec), the output conforms to this structure (doesn't have to be JSON):

```javascript
{
  "meta": {
    "success": true,
    "status_code": 200,
    "status_text": "OK",
    "error_details": "Invalid email"
  },
  "content": {
    // ...
  }
}
```

## Sample Usage

Before jumping in and explaining all of the components, here's a quick look at the usage:

```go
func (c *MyController) SomeFunction() error {
    resp := response.New()
    defer resp.Output()

    // Some API logic...

    if err != nil {
        resp.setResult(http.StatusBadRequest, nil)
        resp.SetErrorDetails("Missing slug parameter")
        return err
    } else {
        resp.SetResult(http.StatusOK, myContent)
        return nil
    }
}
```

## Custom Renderer

A Renderer is effectively the piece that defines how we generate the output.  Response is not built against any specific framework/router, so you're actually going to have to build out a custom renderer to fit your setup.

Before we build one out, here is the interface that defines a Renderer:

```go
type Renderer interface {
    Render(*Response) string
}
```

The `Render` function is passed the `Response` struct that it is bound to, so you can generate the output yourself. This setup, while a little more heavy, does not force this library into a specific format of output.  You can generate JSON / XML / YAML, or whatever your use case calls for!

If we wanted to build out a Renderer that simply generated JSON and printed out the result it would look something like this:

```go
type PrettyJsonRenderer int

func (r *PrettyJsonRenderer) Render(resp *response.Response) string {
    b, err := json.MarshalIndent(resp, "", "    ")
    if err != nil {
        panic("Unable to json.Marshal our Response")
    }
    fmt.Println(b)
    return string(b)
}
```

You will have to set the Renderer before calling `Output()` on the response.  It's easiest to just set this immediately:

```go
resp := response.New().SetRenderer(new(PrettyJsonRenderer))
```

## License

[MIT](LICENSE.md)
