#!/usr/bin/env bash
set -euo pipefail

ROOT="$(dirname "${BASH_SOURCE[0]}")/.."
go build -o "$ROOT/gotify-listen" "$ROOT"
