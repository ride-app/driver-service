load("@rules_buf//buf:defs.bzl", "buf_dependencies")

def buf_deps():
    buf_dependencies(
        name = "buf_deps_pkg_protos",
        modules = ["buf.build/googleapis/googleapis:28151c0d0a1641bf938a7672c500e01d", "buf.build/bufbuild/protovalidate:63dfe56cc2c44cffa4815366ba7a99c0"],
    )
