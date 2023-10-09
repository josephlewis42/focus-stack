
.PHONY: build
build: bin
	go build -o bin/focus-stack .

bin:
	mkdir -p bin

.PHONY: clean
clean:
	rm -rf bin
