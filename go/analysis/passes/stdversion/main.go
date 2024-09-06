// Copyright 2024 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

//go:build ignore
// +build ignore

package main

import (
	"github.com/onboard-inc/golang-tools/go/analysis/singlechecker"
	"github.com/onboard-inc/golang-tools/go/analysis/passes/stdversion"
)

func main() { singlechecker.Main(stdversion.Analyzer) }
