BIN := bin/depviz
SRC := ./src/

.PHONY: build test run run-apk run-graph run-reverse run-render demo vet clean

build:
	go build -o $(BIN) $(SRC)

test:
	go test ./tests/... -v

run: run-graph

run-apk: build
	./$(BIN) testdata/config.example.toml

run-graph: build
	./$(BIN) testdata/config.test.toml

run-reverse: build
	./$(BIN) testdata/config.reverse.toml

run-render: build
	./$(BIN) testdata/config.curl.toml

demo: build
	bash scripts/run-all.sh

vet:
	go vet ./src/... ./tests/...

clean:
	rm -rf bin