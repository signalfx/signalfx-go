# Writers

See the [godoc](https://godoc.org/github.com/signalfx/golib/writer) for
information on what they are and how to use them.

## Compiling

The buffer and writer are generated from a common template in the [`template`](./template)
package.  This module has a `//go:generate` comment that will be recognized when you
run `go generate ./...` on the repo.  For the generate script to work, you must
install the code generation tool with `go get github.com/mauricelam/genny`.

Also, `span_[writer|buffer]_test.go` is automatically generated from
`datapoint_[writer|buffer]_test.go` by the same `go generate` command, so be
sure to make changes in the datapoint test module.

