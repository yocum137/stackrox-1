#!/usr/bin/env bats

load "../helpers.bash"

out_dir=""
original_flavor=""
default_repository="docker.io/stackrox"
custom_repository="example.io/rhacs"
custom_repository_2="my.repo.com"

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

  roxctl
}

setup() {
  out_dir="$(mktemp -d -u)"
  set_image_flavor "development_build"
}

teardown() {
  rm -rf "$out_dir"
  set_image_flavor $original_flavor
}

@test "[development_build] roxctl sensor generate default values" {
  generate_bundle k8s
  assert_components_registry "$out_dir" "$default_repository" "sensor" "collector"
}

@test "[development_build] roxctl sensor generate custom main repository" {
  generate_bundle k8s "--main-image-repository=$custom_repository"
  assert_components_registry "$out_dir" "$custom_repository" "sensor" "collector"
}

@test "[development_build] roxctl sensor generate custom main repository and different collector repository" {
  generate_bundle k8s "--main-image-repository=$custom_repository" "--collector-image-repository=$custom_repository_2"
  assert_components_registry "$out_dir" "$custom_repository" "sensor"
  assert_components_registry "$out_dir" "$custom_repository_2" "collector"
}

@test "[development_build] roxctl sensor generate custom collector repository" {
  generate_bundle k8s "--collector-image-repository=$custom_repository"
  assert_components_registry "$out_dir" "$default_repository" "sensor"
  assert_components_registry "$out_dir" "$custom_repository" "collector"
}
