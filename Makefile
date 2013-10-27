.PHONY: all test fmt benchmark git-add-hook clean

all: test fmt git-version-proxy

git-version-proxy: *.go
	go build .

test:
	go test ./...

fmt:
	go fmt ./...

benchmark:
	go test ./... -bench=".*"

git-pre-commit-hook:
	curl -s 'http://tip.golang.org/misc/git/pre-commit?m=text' > .git/hooks/pre-commit
	chmod +x .git/hooks/pre-commit

clean:
	go clean .
