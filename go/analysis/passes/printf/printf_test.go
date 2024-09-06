// Copyright 2018 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package printf_test

import (
	"testing"

	"github.com/onboard-inc/golang-tools/go/analysis/analysistest"
	"github.com/onboard-inc/golang-tools/go/analysis/passes/printf"
)

func Test(t *testing.T) {
	testdata := analysistest.TestData()
	printf.Analyzer.Flags.Set("funcs", "Warn,Warnf")

	analysistest.Run(t, testdata, printf.Analyzer,
		"a", "b", "nofmt", "typeparams", "issue68744")
	analysistest.RunWithSuggestedFixes(t, testdata, printf.Analyzer, "fix")
}
