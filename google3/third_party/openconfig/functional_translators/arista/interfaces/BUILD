load("//tools/build_defs/go:go_library.bzl", "go_library")
load("//tools/build_defs/go:go_test.bzl", "go_test")

package(
    default_visibility = ["//third_party/openconfig/functional_translators:users"],
    licenses = ["notice"],
)

go_library(
    name = "aristainterfacedesc",
    srcs = ["arista_interface_desc.go"],
    deps = [
        "//third_party/openconfig/functional_translators:ftconsts",
        "//third_party/openconfig/functional_translators:ftutilities",
        "//third_party/openconfig/functional_translators:simplemapper",
        "//third_party/openconfig/functional_translators:translator",
        "//third_party/openconfig/functional_translators/arista/interfaces/yang/openconfig:interfaces",
    ],
)

go_test(
    name = "aristainterfacedesc_test",
    size = "small",
    srcs = ["arista_interface_desc_test.go"],
    library = ":aristainterfacedesc",
    deps = ["//third_party/openconfig/functional_translators:ftutilities"],
)

go_library(
    name = "aristainterfacemac",
    srcs = ["arista_interface_mac.go"],
    deps = [
        "//third_party/openconfig/functional_translators:ftconsts",
        "//third_party/openconfig/functional_translators:ftutilities",
        "//third_party/openconfig/functional_translators:simplemapper",
        "//third_party/openconfig/functional_translators:translator",
        "//third_party/openconfig/functional_translators/arista/interfaces/yang/openconfig:interfaces",
    ],
)

go_test(
    name = "aristainterfacemac_test",
    size = "small",
    srcs = ["arista_interface_mac_test.go"],
    library = ":aristainterfacemac",
)
