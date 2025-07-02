load("//tools/build_defs/go:go_library.bzl", "go_library")
load("//tools/build_defs/go:go_test.bzl", "go_test")

package(
    default_visibility = ["//third_party/openconfig/functional_translators:users"],
    licenses = ["notice"],
)

go_library(
    name = "ftconsts",
    srcs = ["consts.go"],
)

go_library(
    name = "ftutilities",
    srcs = ["utilities.go"],
)

go_test(
    name = "ftutilities_test",
    size = "small",
    srcs = ["utilities_test.go"],
    library = ":ftutilities",
)

go_library(
    name = "simplemapper",
    srcs = ["simplemapper.go"],
    deps = [":ftutilities"],
)

go_test(
    name = "simplemapper_test",
    size = "small",
    srcs = ["simplemapper_test.go"],
    library = ":simplemapper",
    deps = ["//third_party/openconfig/functional_translators/arista/interfaces/yang/openconfig:interfaces"],
)

go_library(
    name = "translator",
    srcs = ["translator.go"],
    deps = [":ftutilities"],
)

go_test(
    name = "translator_test",
    size = "small",
    srcs = ["translator_test.go"],
    library = ":translator",
)
