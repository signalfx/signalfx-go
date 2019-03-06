# Go client library for SignalFx

Note: This library is experimental. Do not rely on it yet.

This is a programmatic interface in Go for SignalFx's metadata and ingest APIs.

# Example

```
import "github.com/signalfx/signalfx-go"

// The client can be customized by backing options onto the end. Check the
// for more info!

// Instantiate your own client if you want to customize it's options
// or test with a RoundTripper
httpClient := &http.Client{…}
client := new Client("your-token-here", HTTPClient(httpClient))
```
