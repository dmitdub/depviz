#!/usr/bin/env sh
set -eu

BIN="./bin/depviz"
SRC="./src/"

if [ ! -x "$BIN" ]; then
    echo "==> Сборка $BIN"
    mkdir -p bin
    go build -o "$BIN" "$SRC"
fi

CFG="${1:-testdata/config.example.toml}"
if [ ! -f "$CFG" ]; then
    echo "Конфигурация не найдена: $CFG" >&2
    exit 1
fi

echo "==> Запуск с конфигом $CFG"
exec "$BIN" "$CFG"