load("@io_bazel_rules_go//go:def.bzl", "go_library")

go_library(
    name = "driverv1alpha1connect",
    srcs = ["driver_service.connect.go"],
    importpath = "github.com/ride-app/driver-service/pkg/protos/ride/driver/v1alpha1/driverv1alpha1connect",
    visibility = ["//visibility:public"],
    deps = [
        "//pkg/protos/ride/driver/v1alpha1",
        "@com_connectrpc_connect//:connect",
    ],
)
