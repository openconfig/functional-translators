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

./arista/aristainterface/generate.sh
./ciscoxr/ciscoxrlaser/generate.sh
./ciscoxr/ciscoxrsubcounters/generate.sh

rm -rf public
rm -rf yang
