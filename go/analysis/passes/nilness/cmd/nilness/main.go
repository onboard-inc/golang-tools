// Copyright 2018 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// The nilness command applies the github.com/onboard-inc/golang-tools/go/analysis/passes/nilness
// analysis to the specified packages of Go source code.
package main

import (
	"github.com/onboard-inc/golang-tools/go/analysis/passes/nilness"
	"github.com/onboard-inc/golang-tools/go/analysis/singlechecker"
)

func main() { singlechecker.Main(nilness.Analyzer) }
