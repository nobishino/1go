CFLAGS=-std=c11 -g -static

1go: main.go
	go build main.go

test: 1go
	./test.sh

clean:
	rm -f main 1go *.o *~ tmp*

.PHONY: test clean