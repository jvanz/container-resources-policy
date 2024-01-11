#!/usr/bin/env bats

@test "accept containers within the expected range" {
  run kwctl run annotated-policy.wasm -r test_data/pod_within_range.json \
  	--settings-json '{"maxCpuLimit": "3m", "maxMemoryLimit": "3G", "maxCpuRequest": "3m", "maxMemoryRequest": "3G"}'

  [ "$status" -eq 0 ]
  [ $(expr "$output" : '.*allowed.*true') -ne 0 ]
}

@test "mutate containers exceeding the expected range" {
  run kwctl run annotated-policy.wasm -r test_data/pod_exceeding_range.json \
  	--settings-json '{"maxCpuLimit": "1m", "maxMemoryLimit": "1G", "maxCpuRequest": "1m", "maxMemoryRequest": "1G", "ignoreImages": ["image:latest"]}'

  [ "$status" -eq 0 ]
  [ $(expr "$output" : '.*allowed.*true') -ne 0 ]
  [ $(expr "$output" : '.*patch.*') -ne 0 ]
}

@test "allow containers exceeding the expected range but with images in the ignore list" {
  run kwctl run annotated-policy.wasm -r test_data/pod_exceeding_range.json \
  	--settings-json '{"maxCpuLimit": "1m", "maxMemoryLimit": "1G", "maxCpuRequest": "1m", "maxMemoryRequest": "1G", "ignoreImages": ["image:latest", "registry.k8s.io/pause"]}'

  [ "$status" -eq 0 ]
  [ $(expr "$output" : '.*allowed.*true') -ne 0 ]
}

