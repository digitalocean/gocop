//go:build tools

// To install the following tools at the version used by this repo run:
// $ go generate -tags tools tools/tools.go

package tools

//go:generate go install gotest.tools/gotestsum@v1.8.2
