# Go client library for SignalFx

Note: This library is experimental. Do not rely on it yet.

This is a programmatic interface in Go for SignalFx's metadata and ingest APIs.

# TODO

* Settle on request/response body deserialization (we're thinking of using OpenTracing generated models to be in sync with docs!)
* Include APIs for metric reading (signalflow, etc!)

# Example

```
import "github.com/signalfx/signalfx-go"

// The client can be customized by backing options onto the end. Check the
// for more info!

// Instantiate your own client if you want to customize its options
// or test with a RoundTripper
httpClient := &http.Client{…}
client := new Client("your-token-here", HTTPClient(httpClient))
```
