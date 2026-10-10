#!/usr/bin/env bash
# Run govulncheck on edgelet module code paths and allow documented exceptions
# listed in SECURITY.md (sync ALLOWED_VULNS with that file).
set -euo pipefail

ROOT="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
cd "$ROOT"

# Documented upstream exceptions — see SECURITY.md § Known vulnerability exceptions.
# Space-separated GO-* IDs; empty when no exceptions are active.
ALLOWED_VULNS=""

out="$(mktemp)"
trap 'rm -f "$out"' EXIT

set +e
govulncheck -format=text ./pkg/... >"$out" 2>&1
status=$?
set -e

cat "$out"

if [[ $status -eq 0 ]]; then
	echo "govulncheck: no vulnerabilities affecting edgelet call paths"
	exit 0
fi

if [[ $status -ne 3 ]]; then
	echo "govulncheck: unexpected exit status $status" >&2
	exit "$status"
fi

found="$(grep -oE 'GO-[0-9]{4}-[0-9]+' "$out" | sort -u || true)"
if [[ -z "$found" ]]; then
	echo "govulncheck: failed but no GO-* IDs parsed; see output above" >&2
	exit 3
fi

unexpected=""
while IFS= read -r id; do
	[[ -z "$id" ]] && continue
	allowed=false
	for a in $ALLOWED_VULNS; do
		if [[ "$id" == "$a" ]]; then
			allowed=true
			break
		fi
	done
	if [[ "$allowed" == false ]]; then
		unexpected="${unexpected} ${id}"
	fi
done <<<"$found"

if [[ -n "${unexpected// /}" ]]; then
	echo "govulncheck: unexpected vulnerabilities (not in SECURITY.md exceptions):${unexpected}" >&2
	exit 3
fi

echo "govulncheck: only documented exceptions remain; see SECURITY.md"
exit 0
