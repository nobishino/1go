CFLAGS=-std=c11 -g -static

1go: main.go
	go build

test: 1go
	./test.sh

clean:
	rm -f 1go *.o *~ tmp*

.PHONY: test clean