//go:build tools

package tools

// Manage tool dependencies via go.mod.
//
// https://github.com/golang/go/wiki/Modules#how-can-i-track-tool-dependencies-for-a-module
// https://github.com/golang/go/issues/25922
import (
	_ "github.com/client9/misspell/cmd/misspell"
	_ "github.com/dkorunic/betteralign/cmd/betteralign"
	_ "github.com/golangci/golangci-lint/cmd/golangci-lint"
	_ "github.com/joe-at-startupmedia/version-bump/v2"
)
