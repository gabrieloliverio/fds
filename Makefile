packages := ./input ./match ./replace

fds: $(packages) main.go
	go build .

test:
	go test -v $(packages)

clean:
	rm -f fds
