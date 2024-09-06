// Copyright 2021 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
package main

import (
	"github.com/onboard-inc/golang-tools/go/analysis/passes/fieldalignment"
	"github.com/onboard-inc/golang-tools/go/analysis/singlechecker"
)

func main() { singlechecker.Main(fieldalignment.Analyzer) }
