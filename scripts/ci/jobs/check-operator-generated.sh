#!/bin/env bash

ROOT="$(cd "$(dirname "${BASH_SOURCE[0]}")"/../../.. && pwd)"
source "$ROOT/scripts/ci/lib.sh"

set -euo pipefail

go mod tidy

# shellcheck disable=SC2016
echo 'Check operator files are up to date (If this fails, run `make -C operator manifests generate bundle` and commit the result.)'
function check-operator-generated-files-up-to-date() {
    make -C operator/ generate
    make -C operator/ manifests
    echo 'Checking for diffs after making generate and manifests...'
    git diff --exit-code HEAD
    make -C operator/ bundle
    echo 'Checking for diffs after making bundle...'
    echo 'If this fails, check if the invocation of the normalize-metadata.py script in operator/Makefile'
    echo 'needs to change due to formatting changes in the generated files.'
    git diff --exit-code HEAD
}
check-operator-generated-files-up-to-date
