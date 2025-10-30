#!/bin/bash
# Deprecated: use `make -C backend client-gen` instead.
echo "[DEPRECATED] backend/generate-client.sh is deprecated. Use: make -C backend client-gen" >&2
make -C "$(dirname "$0")" client-gen
