BIN := bin/depviz
SRC := ./src/

.PHONY: build test run demo vet clean

build:
	go build -o $(BIN) $(SRC)

test:
	go test ./tests/... -v

run: build
	./$(BIN) testdata/config.example.toml

demo: build
	bash scripts/run-all.sh

vet:
	go vet ./src/... ./tests/...

clean:
	rm -rf bin