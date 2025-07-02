load("//tools/build_defs/go:go_library.bzl", "go_library")
load("//tools/build_defs/go:go_test.bzl", "go_test")

package(
    default_visibility = ["//third_party/openconfig/functional_translators:users"],
    licenses = ["notice"],
)

go_library(
    name = "ciscoxrlaser",
    srcs = ["cisco_laser.go"],
    deps = [
        "//third_party/openconfig/functional_translators:ftconsts",
        "//third_party/openconfig/functional_translators:ftutilities",
        "//third_party/openconfig/functional_translators:translator",
        "//third_party/openconfig/functional_translators/ciscoxr/laser/yang/native:gostructs",
        "//third_party/openconfig/functional_translators/ciscoxr/laser/yang/openconfig:gostructs",
    ],
)

go_test(
    name = "ciscoxrlaser_test",
    size = "small",
    srcs = ["cisco_laser_test.go"],
    library = ":ciscoxrlaser",
    deps = ["//third_party/openconfig/functional_translators:ftutilities"],
)
