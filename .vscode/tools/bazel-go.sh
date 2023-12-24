#!/bin/bash

exec .trunk/tools/bazel run -- @io_bazel_rules_go//go/tools/gopackagesdriver "$@"