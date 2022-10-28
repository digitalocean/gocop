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

.phony: tools
tools:
	go generate --tags tools tools/tools.go

.phony: gen-samples
gen-samples:
	gotestsum --format=standard-verbose --hide-summary=all --jsonfile=gocop/testdata/v2/run0.json -- -cover -count=1 -tags="sample" github.com/digitalocean/gocop/sample/... > gocop/testdata/v2/run0.txt 2>&1 || true
	gotestsum --format=standard-verbose --hide-summary=all --jsonfile=gocop/testdata/v2/run1.json -- -cover -count=1 -tags="sample" github.com/digitalocean/gocop/sample/... > gocop/testdata/v2/run1.txt 2>&1 || true
	gotestsum --format=standard-verbose --hide-summary=all --jsonfile=gocop/testdata/v2/run2.json -- -cover -count=1 -tags="sample" github.com/digitalocean/gocop/sample/... > gocop/testdata/v2/run2.txt 2>&1 || true
	gotestsum --format=standard-verbose --hide-summary=all --jsonfile=gocop/testdata/v2/run3.json -- -cover -count=1 -tags="sample" github.com/digitalocean/gocop/sample/... > gocop/testdata/v2/run3.txt 2>&1 || true

.phony: gen-samples-legacy
gen-samples-legacy:
	go test -count=1 -tags="sample" github.com/digitalocean/gocop/sample/... 2>&1 | tee gocop/testdata/run0.txt
	go test -count=1 -tags="sample" github.com/digitalocean/gocop/sample/... 2>&1 | tee gocop/testdata/run1.txt
	go test -count=1 -tags="sample" github.com/digitalocean/gocop/sample/... 2>&1 | tee gocop/testdata/run2.txt
	go test -count=1 -tags="sample" github.com/digitalocean/gocop/sample/... 2>&1 | tee gocop/testdata/run3.txt
