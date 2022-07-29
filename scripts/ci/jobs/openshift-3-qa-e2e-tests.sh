#!/usr/bin/env bash

ROOT="$(cd "$(dirname "${BASH_SOURCE[0]}")"/../../.. && pwd)"
source "$ROOT/scripts/ci/lib.sh"

set -euo pipefail

run() {
    info "openshift 3 cluster??"

    oc get nodes -o wide || true
}

run "$*"
