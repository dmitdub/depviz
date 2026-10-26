#!/usr/bin/env sh
set -eu

BIN="./bin/depviz"
echo "Сборка $BIN"
mkdir -p bin
go build -o "$BIN" ./src/

CFG="${1:-testdata/config.test.toml}"
if [ ! -f "$CFG" ]; then
    echo "Конфигурация не найдена: $CFG" >&2
    exit 1
fi

echo "Запуск с конфигом $CFG"
exec "$BIN" "$CFG"