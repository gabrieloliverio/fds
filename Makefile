test:
	go test -cover ./...

coverage:
	go test ./... -coverprofile=coverage.out \
		&& go tool cover -html=coverage.out

clean:
	rm -f fds coverage.out
