#!/usr/bin/env bash
set -euo pipefail

if [[ $# -ne 1 ]]; then
  echo "usage: release.sh <tag>  (e.g. v1.0.0)"
  echo "latest tag: $(git -C "$(dirname "${BASH_SOURCE[0]}")" describe --tags --abbrev=0 2>/dev/null || echo 'none')"
  exit 1
fi

TAG="$1"
ROOT="$(dirname "${BASH_SOURCE[0]}")/.."

echo "Building..."
bash "$(dirname "${BASH_SOURCE[0]}")/build.sh"

echo "Creating release $TAG..."
gh release create "$TAG" \
  --repo "$(gh repo view --json nameWithOwner -q .nameWithOwner)" \
  --title "$TAG" \
  --generate-notes \
  "$ROOT/gotify-listen"

echo "Done: $TAG"
