GOFLAGS := -mod=vendor
export GOFLAGS

build:
	mkdir -p .build
	rm -f .build/gocop-linux-amd64
	go build -o .build/gocop-linux-amd64

.phony: unit
unit:
	go test -v -cover github.com/digitalocean/gocop/gocop

.phony: component
component:
	go test github.com/digitalocean/gocop

.phony: gen-samples
gen-samples:
	go test -count=1 github.com/digitalocean/gocop/sample/... 2>&1 | tee gocop/testdata/run0.txt
	go test -count=1 github.com/digitalocean/gocop/sample/... 2>&1 | tee gocop/testdata/run1.txt
	go test -count=1 github.com/digitalocean/gocop/sample/... 2>&1 | tee gocop/testdata/run2.txt
	go test -count=1 github.com/digitalocean/gocop/sample/... 2>&1 | tee gocop/testdata/run3.txt
