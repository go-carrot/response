<a href="https://engineering.carrot.is/"><p align="center"><img src="https://cloud.githubusercontent.com/assets/2105067/24525319/d3d26516-1567-11e7-9506-7611b3287d53.png" alt="Go Carrot" width="350px" align="center;" /></p></a>
# Response

[![Build Status](https://travis-ci.org/go-carrot/response.svg?branch=master)](https://travis-ci.org/go-carrot/response) [![codecov](https://codecov.io/gh/go-carrot/response/branch/master/graph/badge.svg)](https://codecov.io/gh/go-carrot/response)

Response is a library that generates output that follows [carrot/restful-api-spec](https://github.com/carrot/restful-api-spec).

## Output Interface

The output conforms to this structure:

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

## Usage

Dig into [response.go](/response.go) for specifics, but this code sample covers all usage:

```go
func (c *MyController) Index(w http.ResponseWriter, r *http.Request) {
    resp := response.New(w)
    defer resp.Output()

    models, err := getModels()
    if err != nil {
        resp.SetErrorDetails(err.Error())
        resp.setResult(http.StatusInternalServerError, nil)
        return
    }

    resp.SetResult(http.StatusOK, models)
}
```

## License

[MIT](LICENSE.md)
