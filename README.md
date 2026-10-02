# req

Authenticated HTTP requests as a chain. `Do` reads and closes the body for you, and returns a `*RequestError` tagged with the stage it failed at.

```
go get github.com/digiconvent/req
```

## basic usage

```go
body, resp, err := req.Get(url).Do()
if err != nil {
    // body and resp are nil - something failed before or while getting a response
}
```

## posting a body

```go
body, resp, err := req.Post(url, strings.NewReader(payload)).Do()
```

## adding headers

`.Header` chains, so add as many as you need, on either `Get` or `Post`:

```go
body, resp, err := req.Get(url).
    Header("Accept", "application/vnd.github.v3+json").
    Header("Authorization", "Bearer "+token).
    Do()
```

## decoding JSON directly

`Into` is an alternative to `Do` - it sends the request and unmarshals the body into `v`, returning just the error:

```go
var repo struct {
    Name string `json:"name"`
}
err := req.Get(url).Into(&repo)
```

## telling errors apart

Every error `Do` returns is a `*req.RequestError`, tagged with the `Stage` it happened at (`StageBuildRequest`, `StageSendRequest`, `StageReadBody`, `StageStatus`):

```go
body, resp, err := req.Get(url).Do()

var reqErr *req.RequestError
if errors.As(err, &reqErr) {
    switch reqErr.Stage {
    case req.StageStatus:
        // got a response, non-2xx - body and resp are still populated
    default:
        // never reached the server
    }
}
```

A non-2xx response wraps a `*req.StatusError` with the code and body, so you can get at it directly without checking the stage:

```go
var statusErr *req.StatusError
if errors.As(err, &statusErr) {
    statusErr.StatusCode
    statusErr.Body
}
```

## swapping the client

`req.Client` is a plain `*http.Client` var - point it at a test server, a client with a timeout, whatever:

```go
req.Client = &http.Client{Timeout: 5 * time.Second}
```
