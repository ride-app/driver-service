load("@io_bazel_rules_go//go:def.bzl", "go_library")
load("@io_bazel_rules_go//proto:def.bzl", "go_proto_library")
load("@rules_buf//buf:defs.bzl", "buf_lint_test")
load("@rules_proto//proto:defs.bzl", "proto_library")

proto_library(
    name = "ride_driver_v1alpha1_proto",
    srcs = ["driver_service.proto"],
    visibility = ["//visibility:public"],
    deps = [
        "@buf_deps_pkg_protos//google/api:api_proto",
        "@buf_deps_pkg_protos//google/type:type_proto",
        "@buf_deps_pkg_protos//validate:validate_proto",
        "@com_google_protobuf//:field_mask_proto",
        "@com_google_protobuf//:timestamp_proto",
    ],
)

go_proto_library(
    name = "ride_driver_v1alpha1_go_proto",
    compilers = ["@io_bazel_rules_go//proto:go_grpc"],
    importpath = "github.com/ride-app/driver-service/pkg/protos/ride/driver/v1alpha1",
    proto = ":ride_driver_v1alpha1_proto",
    visibility = ["//visibility:public"],
    deps = [
        "//google/api:annotations_proto",
        "//google/type:date_proto",
        "//google/type:phone_number_proto",
        "//validate:validate_proto",
    ],
)

go_library(
    name = "v1alpha1",
    srcs = ["driver_service.pb.validate.go"],
    embed = [":ride_driver_v1alpha1_go_proto"],
    importpath = "github.com/ride-app/driver-service/pkg/protos/ride/driver/v1alpha1",
    visibility = ["//visibility:public"],
    deps = ["@org_golang_google_protobuf//types/known/anypb"],
)

buf_lint_test(
    name = "ride_driver_v1alpha1_proto_lint",
    config = "//pkg/protos:buf.yaml",
    targets = [":ride_driver_v1alpha1_proto"],
)
