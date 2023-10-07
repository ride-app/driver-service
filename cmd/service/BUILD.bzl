load("@io_bazel_rules_go//go:def.bzl", "go_library")

go_library(
    name = "service",
    srcs = [
        "main.go",
        "wire_gen.go",
    ],
    importpath = "github.com/ride-app/driver-service/cmd/service",
    visibility = ["//visibility:public"],
    deps = [
        "//pkg/api",
        "//pkg/api/interceptors",
        "//pkg/config",
        "//pkg/protos/ride/driver/v1alpha1/driverv1alpha1connect",
        "//pkg/repositories/driver",
        "//pkg/repositories/vehicle",
        "//pkg/repositories/wallet",
        "//pkg/third-party",
        "//pkg/utils/logger",
        "@com_connectrpc_connect//:connect",
        "@org_golang_x_net//http2",
        "@org_golang_x_net//http2/h2c",
    ],
)
