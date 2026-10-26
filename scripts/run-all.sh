#!/usr/bin/env sh
set -eu

cd "$(dirname "$0")/.."

echo "run.sh (успешный запуск)"
bash scripts/run.sh testdata/config.example.toml

echo
echo "run-errors.sh (обработка ошибок)"
bash scripts/run-errors.sh