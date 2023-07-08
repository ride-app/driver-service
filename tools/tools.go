//go:build tools

package tools

import (
	_ "github.com/google/wire/cmd/wire"
	_ "github.com/onsi/ginkgo/v2/ginkgo"
	_ "go.uber.org/mock/mockgen"
)
