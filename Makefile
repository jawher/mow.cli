test:
	go test -v ./...

test-cover:
	go list -f '{{if len .TestGoFiles}}"go test -coverprofile={{.Dir}}/.coverprofile {{.ImportPath}}"{{end}}' ./... | xargs -L 1 sh -c
	gover
	goveralls -coverprofile=gover.coverprofile -service=travis-ci

check: lint vet fmtcheck ineffassign readmecheck

lint:
	golint -set_exit_status ./...

vet:
	go vet

doc:
	autoreadme -f

fmtcheck:
	@ export output="$$(gofmt -s -d .)"; \
		[ -n "$${output}" ] && echo "$${output}" && export status=1; \
		exit $${status:-0}

ineffassign:
	ineffassign .

readmecheck:
	sed '$ d' README.md > README.original.md
	autoreadme -f
	sed '$ d' README.md > README.generated.md
	diff README.generated.md README.original.md

setup:
	go get -u -m github.com/gordonklaus/ineffassign
	go get -u -m github.com/golang/lint/golint
	go get -u -m golang.org/x/tools/cmd/cover
	go get -u -m github.com/mattn/goveralls
	go get -u -m github.com/modocache/gover
	go get -u -m github.com/divan/autoreadme
	go get -t -u ./...

.PHONY: test check lint vet fmtcheck ineffassign readmecheck
