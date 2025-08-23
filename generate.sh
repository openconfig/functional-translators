#!/bin/bash

OC_VERSION="v5.3.0"

git clone https://github.com/openconfig/public.git --branch $OC_VERSION
git clone https://github.com/YangModels/yang.git

./arista/aristainterface/generate.sh
./ciscoxr/ciscoxrlaser/generate.sh

rm -rf public
rm -rf yang
