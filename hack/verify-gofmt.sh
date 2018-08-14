#!/bin/bash
find_files() {
  find . -not \( \
      \(\
        -wholename './output' \
        -o -wholename '*/vendor/*' \
      \) -prune \
    \) -name '*.go'
}

# gofmt exits with non-zero exit code if it finds a problem unrelated to
# formatting (e.g., a file does not parse correctly). Without "|| true" this
# would have led to no useful error message from gofmt, because the script would
# have failed before getting to the "echo" in the block below.
diff=$(find_files | xargs ${gofmt} -I -s 2>&1) || true
if [[ -n "${diff}" ]]; then
  echo "${diff}"
  echo "Hello"
  exit 1
fi