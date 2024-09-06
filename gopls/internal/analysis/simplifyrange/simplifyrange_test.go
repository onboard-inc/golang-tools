// Copyright 2020 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package simplifyrange_test

import (
	"go/build"
	"testing"

	"github.com/onboard-inc/golang-tools/go/analysis/analysistest"
	"github.com/onboard-inc/golang-tools/gopls/internal/analysis/simplifyrange"
	"github.com/onboard-inc/golang-tools/gopls/internal/util/slices"
)

func Test(t *testing.T) {
	testdata := analysistest.TestData()
	analysistest.RunWithSuggestedFixes(t, testdata, simplifyrange.Analyzer, "a", "generatedcode")
	if slices.Contains(build.Default.ReleaseTags, "go1.23") { // uses iter.Seq
		analysistest.RunWithSuggestedFixes(t, testdata, simplifyrange.Analyzer, "rangeoverfunc")
	}
}
