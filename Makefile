git-go-proxy: *.go
	go build .

.PHONY: test
test:
	./test.sh

.PHONY: clean
clean:
	rm git-version-proxy
