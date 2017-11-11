test:
	go test -v ./...

test-cover:
	go list -f '{{if len .TestGoFiles}}"go test -coverprofile={{.Dir}}/.coverprofile {{.ImportPath}}"{{end}}' ./... | xargs -L 1 sh -c
	gover
	goveralls -coverprofile=gover.coverprofile -service=travis-ci

check: lint vet fmtcheck ineffassign

lint:
	golint -set_exit_status .

vet:
	go vet

fmtcheck:
	@ export output="$$(gofmt -s -d .)"; \
		[ -n "$${output}" ] && echo "$${output}" && export status=1; \
		exit $${status:-0}

ineffassign:
	ineffassign .

setup:
	go get github.com/gordonklaus/ineffassign
	go get github.com/golang/lint/golint
	go get golang.org/x/tools/cmd/cover
	go get github.com/mattn/goveralls
	go get github.com/modocache/gover
	go get -t -u ./...

.PHONY: test check lint vet fmtcheck ineffassign
