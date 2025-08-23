#!/bin/bash

    # cmd = "$(locations //third_party/golang/ygot/generator:generator) " +
    #       # We explicitly do not compress paths so we can get/set State paths.
    #       "-compress_paths=false -exclude_modules=ietf-interfaces " +
    #       "-output_dir=$(RULEDIR) -package_name=interfaces " +
    #       "-generate_fakeroot -fakeroot_name=device " +
    #       "-typedef_enum_with_defmod " +
    #       "-enum_suffix_for_simple_union_enums " +
    #       "-generate_simple_unions " +
    #       "-ignore_circdeps " +
    #       "-ignore_unsupported=true " +
    #       "-structs_split_files_count=5 " +
    #       "-path=third_party/openconfig/public/release/models,third_party/openconfig/public/third_party " +
    #       "-ygot_path=google3/third_party/golang/ygot/ygot/ygot " +
    #       "-ytypes_path=google3/third_party/golang/ygot/ytypes/ytypes " +
    #       "-goyang_path=google3/third_party/golang/goyang/pkg/yang/yang " +
    #       "third_party/openconfig/public/release/models/interfaces/openconfig-interfaces.yang " +
    #       "third_party/openconfig/public/release/models/interfaces/openconfig-if-ethernet.yang " +
    #       "third_party/openconfig/public/release/models/lacp/openconfig-lacp.yang " +
    #       "third_party/openconfig/functional_translators/arista/interfaces/yang/openconfig/restriction.yang" +
    #       " && $(locations //third_party/go:gofmt) -w $(OUTS)" +
    #       " && //third_party/golang/go_tools/cmd/goimports:goimports -w $(OUTS)",

YANG_FILES=(
  public/release/models/interfaces/openconfig-interfaces.yang
  public/release/models/interfaces/openconfig-if-ethernet.yang
  public/release/models/lacp/openconfig-lacp.yang
  arista/aristainterface/yang/openconfig/restriction.yang
)

go run github.com/openconfig/ygot/generator \
  -compress_paths=false \
  -exclude_modules=ietf-interfaces \
  -output_dir=arista/aristainterface/yang/openconfig -package_name=interfaces \
  -generate_fakeroot -fakeroot_name=device \
  -typedef_enum_with_defmod \
  -enum_suffix_for_simple_union_enums \
  -generate_simple_unions \
  -ignore_circdeps \
  -ignore_unsupported=true \
  -structs_split_files_count=5 \
  -path=public/release/models,public/third_party \
  "${YANG_FILES[@]}" \
  #-ygot_path=google3/third_party/golang/ygot/ygot/ygot \
  #-ytypes_path=google3/third_party/golang/ygot/ytypes/ytypes \
  #-goyang_path=google3/third_party/golang/goyang/pkg/yang/yang \

gofmt -w arista/aristainterface/yang/openconfig/*.go
goimports -w arista/aristainterface/yang/openconfig/*.go


# go run github.com/openconfig/ygnmi/app/ygnmi generator \
#   --trim_module_prefix=openconfig \
#   --typedef_enum_with_defmod=false \
#   --exclude_modules="${EXCLUDE_MODULES}" \
#   --base_package_path=github.com/openconfig/ondatra/gnmi/oc \
#   --output_dir=gnmi/oc \
#   --paths=public/release/models/...,public/third_party/ietf/... \
#   --split_package_paths="/network-instances/network-instance/protocols/protocol/isis=netinstisis,/network-instances/network-instance/protocols/protocol/bgp=netinstbgp" \
#   --ignore_deviate_notsupported \
#   "${YANG_FILES[@]}"

