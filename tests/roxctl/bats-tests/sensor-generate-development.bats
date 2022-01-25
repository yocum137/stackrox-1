#!/usr/bin/env bats

load "helpers.bash"

out_dir=""
original_flavor=""

setup_file() {
  echo "Testing roxctl version: '$(roxctl-development version)'" >&3
  command -v yq || skip "Tests in this file require yq"
  [[ -n "$API_ENDPOINT" ]] || skip "API_ENDPOINT environment variable required"
  [[ -n "$ROX_PASSWORD" ]] || skip "ROX_PASSWORD environment variable required"
  original_flavor=$(kubectl -n stackrox exec -it deployment/central -- env | grep -i ROX_IMAGE_FLAVOR | sed 's/ROX_IMAGE_FLAVOR=//')
}

set_image_flavor() {
  flavor=$1; shift
  kubectl -n stackrox set env deployment/central ROX_IMAGE_FLAVOR="$flavor"
}

roxctl_cmd() {
  roxctl --insecure-skip-tls-verify -e "$API_ENDPOINT" -p "$ROX_PASSWORD" "$@"
}

setup() {
  out_dir="$(mktemp -d -u)"
}

teardown() {
  rm -rf "$out_dir"
  set_image_flavor $original_flavor
}

no_override_test() {
  local main_default="$1"; shift
  local collector_default="$2"; shift
  generate_bundle k8s
  assert_image_starts_with "sensor" "$main_default"
  assert_image_starts_with "collector" "$collector_default"
}

test_main_repository_override() {
  generate_bundle k8s '--main-image-repository=example.io/rhacs/main'
  assert_image_starts_with "sensor" "example.io/rhacs"
  assert_image_starts_with "collector" "example.io/rhacs"
}

@test "[development_build] roxctl sensor generate: derive collector image using --main-image-repository override" {
  set_image_flavor "development_build"
  test_main_repository_override
}

@test "[stackrox.io] roxctl sensor generate: derive collector image using --main-image-repository override" {
  set_image_flavor "stackrox.io"
  test_main_repository_override
}

@test "[rhacs] roxctl sensor generate: derive collector image using --main-image-repository override" {
  set_image_flavor "rhacs"
  test_main_repository_override
}

@test "[development_build]"

#@test "[development_build] roxctl-development sensor generate k8s generates bundle for docker.io" {
#  set_image_flavor "development_build"
#  generate_bundle k8s
#  assert_image_starts_with "sensor" "docker.io/stackrox/main"
#  assert_image_starts_with "collector" "docker.io/stackrox/collector"
#}
#
#@test "[development_build] roxctl-development sensor generate k8s generates bundle with main override" {
#  set_image_flavor "development_build"
#  generate_bundle k8s '--main-image-repository=example.io/rhacs/main'
#  assert_image_starts_with "sensor" "example.io/rhacs/main"
#  assert_image_starts_with "collector" "example.io/rhacs/collector"
#}
#
#@test "[development_build] roxctl-development sensor generate k8s generates bundle with main override" {
#  set_image_flavor "development_build"
#  generate_bundle k8s '--main-image-repository=example.io/rhacs/main' '--collector-image-repository=collector.example.io/rhacs/collector'
#  assert_image_starts_with "sensor" "example.io/rhacs/main"
#  assert_image_starts_with "collector" "collector.example.io/rhacs/collector"
#}
