
.PHONY: build
build: bin
	go build -o bin/stack .

bin:
	mkdir -p bin

.PHONY: clean
clean:
	rm -rf bin
