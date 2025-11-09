#!/usr/bin/env bash
set -euo pipefail
DIR="$(cd "$(dirname "$0")" && pwd)"

# Generate the client sources from backend OpenAPI
pushd "$DIR/../.." >/dev/null
if ! make -C backend client-gen; then
  echo "Warning: client generation failed (likely offline). Preserving existing client sources." >&2
fi
popd >/dev/null

# Compile the client to dist using the web TypeScript compiler
WEB_DIR="$(cd "$DIR/.." && pwd)"
"$WEB_DIR/node_modules/.bin/tsc" -p "$WEB_DIR/api-client/tsconfig.json"

echo "Client generated and built under web/api-client/dist"
