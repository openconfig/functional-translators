#!/bin/bash
# Copyright 2025 Google LLC
#
# Licensed under the Apache License, Version 2.0 (the "License");
# you may not use this file except in compliance with the License.
# You may obtain a copy of the License at
#
#     https://www.apache.org/licenses/LICENSE-2.0
#
# Unless required by applicable law or agreed to in writing, software
# distributed under the License is distributed on an "AS IS" BASIS,
# WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
# See the License for the specific language governing permissions and
# limitations under the License.

OC_OUT_PATH=ciscoxr/ciscoxrfpd/yang/openconfig

OC_YANG_FILES=(
  public/release/models/platform/openconfig-platform.yang
  public/release/models/platform/openconfig-platform-types.yang
  public/release/models/platform/openconfig-platform-port.yang
  public/release/models/platform/openconfig-platform-common.yang
  public/release/models/platform/openconfig-platform-integrated-circuit.yang
  public/release/models/platform/openconfig-platform-transceiver.yang
  public/release/models/interfaces/openconfig-interfaces.yang
  public/release/models/types/openconfig-types.yang
  public/release/models/optical-transport/openconfig-transport-types.yang
  public/release/models/openconfig-extensions.yang
  public/release/models/types/openconfig-yang-types.yang
  public/release/models/system/openconfig-alarm-types.yang
  public/release/models/platform/openconfig-platform-pipeline-counters.yang
)

mkdir -p "$OC_OUT_PATH"

go run github.com/openconfig/ygot/generator \
  -annotations \
  -compress_paths=false \
  -exclude_modules=ietf-interfaces \
  -output_dir="${OC_OUT_PATH}" \
  -package_name=openconfig \
  -generate_fakeroot \
  -fakeroot_name=device \
  -typedef_enum_with_defmod \
  -enum_suffix_for_simple_union_enums \
  -generate_simple_unions \
  -generate_append \
  -generate_getters \
  -generate_rename \
  -generate_delete \
  -ignore_circdeps \
  -ignore_unsupported=true \
  -structs_split_files_count=5 \
  -path=public/release/models,public/third_party \
  "${OC_YANG_FILES[@]}"

gofmt -w "${OC_OUT_PATH}"/*.go
goimports -w "${OC_OUT_PATH}"/*.go