test:
	go test ./...

check: readmecheck
	bin/golangci-lint run

doc:
	autoreadme -f

readmecheck:
	head -n -1 README.md > README.original.md
	autoreadme -f
	head -n -1 README.md > README.generated.md
	diff README.generated.md README.original.md

lint.setup:
	curl -sfL https://install.goreleaser.com/github.com/golangci/golangci-lint.sh | sh -s v1.28.0


readmecheck.setup:
	go get github.com/jawher/autoreadme

.PHONY: test check doc readmecheck lint.setup readmecheck.setup
