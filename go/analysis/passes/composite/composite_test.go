// Copyright 2018 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package composite_test

import (
	"testing"

	"github.com/onboard-inc/golang-tools/go/analysis/analysistest"
	"github.com/onboard-inc/golang-tools/go/analysis/passes/composite"
)

func Test(t *testing.T) {
	testdata := analysistest.TestData()
	analysistest.RunWithSuggestedFixes(t, testdata, composite.Analyzer, "a", "typeparams")
}
