load("@rules_buf//buf:defs.bzl", "buf_dependencies")

def buf_deps():
    buf_dependencies(
        name = "buf_deps_pkg_protos",
        modules = [
            "buf.build/envoyproxy/protoc-gen-validate:eac44469a7af47e7839a7f1f3d7ac004",
            "buf.build/googleapis/googleapis:28151c0d0a1641bf938a7672c500e01d",
        ],
    )
