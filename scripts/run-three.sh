#!/usr/bin/env sh
set -eu

for cfg in testdata/config.curl.toml testdata/config.vim.toml testdata/config.openssl.toml; do
    echo "=== $cfg ==="
    bash scripts/run-render.sh "$cfg"
    echo
done