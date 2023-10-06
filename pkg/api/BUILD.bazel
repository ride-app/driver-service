load("@io_bazel_rules_go//go:def.bzl", "go_library", "go_test")

go_library(
    name = "api",
    srcs = [
        "create_driver.go",
        "delete_driver.go",
        "get_driver.go",
        "get_location.go",
        "get_status.go",
        "get_vehicle.go",
        "go_offline.go",
        "go_online.go",
        "service.go",
        "update_driver.go",
        "update_location.go",
        "update_vehicle.go",
    ],
    importpath = "github.com/ride-app/driver-service/pkg/api",
    visibility = ["//visibility:public"],
    deps = [
        "//pkg/protos/ride/driver/v1alpha1",
        "//pkg/repositories/driver",
        "//pkg/repositories/vehicle",
        "//pkg/repositories/wallet",
        "//pkg/utils/logger",
        "@com_connectrpc_connect//:connect",
        "@org_golang_google_protobuf//types/known/timestamppb",
    ],
)

go_test(
    name = "api_test",
    srcs = [
        "delete_driver_test.go",
        "get_vehicle_test.go",
        "service_suite_test.go",
        "test_helpers_test.go",
        "update_driver_test.go",
        "update_vehicle_test.go",
    ],
    deps = [
        ":api",
        "//pkg/protos/ride/driver/v1alpha1",
        "//pkg/testing/mocks",
        "@com_connectrpc_connect//:connect",
        "@com_github_onsi_ginkgo_v2//:ginkgo",
        "@com_github_onsi_gomega//:gomega",
        "@org_golang_google_genproto//googleapis/type/date",
        "@org_golang_google_protobuf//reflect/protoreflect",
        "@org_golang_google_protobuf//types/known/timestamppb",
        "@org_uber_go_mock//gomock",
    ],
)
