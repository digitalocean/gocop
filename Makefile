.phony: unit
unit:
	GOFLAGS=-mod=vendor go test -v -cover github.com/digitalocean/gocop/gocop

.phony: component
component:
	GOFLAGS=-mod=vendor go test github.com/digitalocean/gocop

.phony: gen-samples
gen-samples:
	GOFLAGS=-mod=vendor go test -count=1 github.com/digitalocean/gocop/sample/... 2>&1 | tee gocop/testdata/run0.txt
	GOFLAGS=-mod=vendor go test -count=1 github.com/digitalocean/gocop/sample/... 2>&1 | tee gocop/testdata/run1.txt
	GOFLAGS=-mod=vendor go test -count=1 github.com/digitalocean/gocop/sample/... 2>&1 | tee gocop/testdata/run2.txt
	GOFLAGS=-mod=vendor go test -count=1 github.com/digitalocean/gocop/sample/... 2>&1 | tee gocop/testdata/run3.txt
