# Response

[![Build Status](https://travis-ci.org/dcstack/response.svg?branch=master)](https://travis-ci.org/dcstack/response) [![codecov](https://codecov.io/gh/dcstack/response/branch/master/graph/badge.svg)](https://codecov.io/gh/dcstack/response)

Response is a library that flexibly generates output that follows [carrot/restful-api-spec](https://github.com/carrot/restful-api-spec).

## Output Interface

Following [carrot/restful-api-spec](https://github.com/carrot/restful-api-spec), the output conforms to this structure (doesn't have to be JSON):

```json
{
    "success": true,
    "status_code": 200,
    "status_text": "OK",
    "error_details": [
        {
            "code": 1,
            "text": "Invalid email"
        }
    ],
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
        resp.AddErrorDetail(errors.MissingSlugParameter)
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

## Custom Errors

In our output interface, we have a `error_details` array.  This is the place where we're to put additional details about a specific error that aren't fully explained with HTTP status codes.  This array is intended to be both human and machine readable, so each error detail has both a code and text.

To be able to call `AddErrorDetail(...)` on a response, you must first register error codes to this library.  You can accomplish this by calling `SetErrorMap`, which takes a `map[int]string`.

```go
const (
    ErrorMissingAuth      = 1
    ErrorMissingParameter = 2
)
myErrors := map[int]string{
    ErrorMissingAuth:      "Missing Auth",
    ErrorMissingParameter: "Missing Parameter",
}
response.SetErrorMap(myErrors)
```

Now that we've registered all of our errors, we can simply call `AddErrorDetail` on our responses:

```
resp.AddErrorDetail(ErrorMissingAuth)
```

## License

[MIT](LICENSE.md)
