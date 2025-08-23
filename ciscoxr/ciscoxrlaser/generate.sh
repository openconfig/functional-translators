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

# native

YANG_FILES=(
  yang/vendor/cisco/xr/2431/Cisco-IOS-XR-controller-optics-oper.yang
  yang/vendor/cisco/xr/2431/Cisco-IOS-XR-controller-optics-oper-sub1.yang
  yang/vendor/cisco/xr/2431/Cisco-IOS-XR-controller-optics-oper-sub2.yang
  yang/vendor/cisco/xr/2431/cisco-semver.yang
  yang/vendor/cisco/xr/2431/ietf-inet-types.yang
)

go run github.com/openconfig/ygot/generator \
  -compress_paths=false \
  -exclude_modules=ietf-interfaces \
  -output_dir=ciscoxr/ciscoxrlaser/yang/native -package_name=xr2431 \
  -ignore_circdeps \
  -ignore_unsupported=true \
  -generate_fakeroot -fakeroot_name=CiscoDevice \
  -generate_simple_unions \
  -generate_append -generate_getters -generate_rename \
  -structs_split_files_count=3 \
  -path=yang/vendor/cisco/xr/2431 \
  "${YANG_FILES[@]}"

gofmt -w ciscoxr/ciscoxrlaser/yang/native/*.go
goimports -w ciscoxr/ciscoxrlaser/yang/native/*.go

# openconfig

YANG_FILES=(
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
)

go run github.com/openconfig/ygot/generator \
  -annotations \
  -compress_paths=false \
  -exclude_modules=ietf-interfaces \
  -output_dir=ciscoxr/ciscoxrlaser/yang/openconfig -package_name=integratedcircuit \
  -generate_fakeroot -fakeroot_name=device \
  -typedef_enum_with_defmod \
  -enum_suffix_for_simple_union_enums \
  -generate_simple_unions \
  -generate_append -generate_getters -generate_rename \
  -generate_delete \
  -ignore_circdeps \
  -ignore_unsupported=true \
  -structs_split_files_count=5 \
  -path=public/release/models,public/third_party \
  "${YANG_FILES[@]}"

gofmt -w ciscoxr/ciscoxrlaser/yang/openconfig/*.go
goimports -w ciscoxr/ciscoxrlaser/yang/openconfig/*.go


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

