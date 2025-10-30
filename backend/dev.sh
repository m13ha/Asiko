#!/bin/bash
# Deprecated: use `make dev` instead.
echo "[DEPRECATED] backend/dev.sh is deprecated. Use: make dev" >&2
make -C "$(dirname "$0")" dev
