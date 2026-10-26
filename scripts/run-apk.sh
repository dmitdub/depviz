#!/usr/bin/env sh
set -eu

BIN="./bin/depviz"
echo "Сборка $BIN"
mkdir -p bin
go build -o "$BIN" ./src/

echo "Запуск с конфигом testdata/config.example.toml"
exec "$BIN" testdata/config.example.toml