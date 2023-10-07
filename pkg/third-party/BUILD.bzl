load("@io_bazel_rules_go//go:def.bzl", "go_library")

go_library(
    name = "third-party",
    srcs = ["firebase.go"],
    importpath = "github.com/ride-app/driver-service/pkg/third-party",
    visibility = ["//visibility:public"],
    deps = [
        "//pkg/config",
        "//pkg/utils/logger",
        "@com_google_firebase_go_v4//:go",
    ],
)
