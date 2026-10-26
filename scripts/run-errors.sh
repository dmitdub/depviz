#!/usr/bin/env sh
set -u

BIN="./bin/depviz"
echo "Сборка $BIN"
mkdir -p bin
go build -o "$BIN" ./src/

TMP="$(mktemp -d)"
trap 'rm -rf "$TMP"' EXIT

check_fail() {
    name="$1"; shift
    echo "$name"
    if "$BIN" "$@" >/dev/null 2>"$TMP/err"; then
        echo "FAIL: ожидался код 1"; return 1
    fi
    echo "OK (код 1): $(cat "$TMP/err")"
}

echo "1. Несуществующий файл"
check_fail "missing file" "$TMP/nope.toml"

cat > "$TMP/empty_pkg.toml" <<EOF
repository_url = "x"
output_image = "o"
max_depth = 1
EOF

cat > "$TMP/empty_repo.toml" <<EOF
package_name = "p"
output_image = "o"
max_depth = 1
EOF

cat > "$TMP/bad_depth.toml" <<EOF
package_name = "p"
repository_url = "x"
output_image = "o"
max_depth = 0
EOF

cat > "$TMP/empty_image.toml" <<EOF
package_name = "p"
repository_url = "x"
max_depth = 1
EOF

cat > "$TMP/test_no_path.toml" <<EOF
package_name = "p"
output_image = "o"
max_depth = 1
test_repository = true
EOF

echo "2. package_name пуст"
check_fail "empty package_name" "$TMP/empty_pkg.toml"
echo "3. repository_url пуст"
check_fail "empty repository_url" "$TMP/empty_repo.toml"
echo "4. max_depth = 0"
check_fail "bad max_depth" "$TMP/bad_depth.toml"
echo "5. output_image пуст"
check_fail "empty output_image" "$TMP/empty_image.toml"
echo "6. test_repository без test_file_path"
check_fail "test without path" "$TMP/test_no_path.toml"

echo
echo "Все проверки ошибок пройдены."