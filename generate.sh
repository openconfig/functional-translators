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

OC_VERSION="v5.3.0"

git clone https://github.com/openconfig/public.git --branch $OC_VERSION
git clone https://github.com/YangModels/yang.git

# go/keep-sorted start
bash arista/aristainterface/generate.sh
bash ciscoxr/ciscoxr8000icresource/generate.sh
bash ciscoxr/ciscoxrarp/generate.sh
bash ciscoxr/ciscoxrcarrier/generate.sh
bash ciscoxr/ciscoxrfabric/generate.sh
bash ciscoxr/ciscoxrfpd/generate.sh
bash ciscoxr/ciscoxrfragment/generate.sh
bash ciscoxr/ciscoxrlagmac/generate.sh
bash ciscoxr/ciscoxrlaser/generate.sh
bash ciscoxr/ciscoxrmount/generate.sh
bash ciscoxr/ciscoxrpower/generate.sh
bash ciscoxr/ciscoxrqos/generate.sh
bash ciscoxr/ciscoxrsubcounters/generate.sh
# go/keep-sorted end

rm -rf public
rm -rf yang
