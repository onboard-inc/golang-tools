// Copyright 2018 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package loopclosure_test

import (
	"path/filepath"
	"testing"

	"github.com/onboard-inc/golang-tools/go/analysis/analysistest"
	"github.com/onboard-inc/golang-tools/go/analysis/passes/loopclosure"
	"github.com/onboard-inc/golang-tools/internal/testenv"
	"github.com/onboard-inc/golang-tools/internal/testfiles"
)

func Test(t *testing.T) {
	// legacy loopclosure test expectations are incorrect > 1.21.
	testenv.SkipAfterGo1Point(t, 21)

	testdata := analysistest.TestData()
	analysistest.Run(t, testdata, loopclosure.Analyzer,
		"a", "golang.org/...", "subtests", "typeparams")
}

func TestVersions22(t *testing.T) {
	t.Skip("Disabled for golang.org/cl/603895. Fix and re-enable.")
	testenv.NeedsGo1Point(t, 22)

	dir := testfiles.ExtractTxtarFileToTmp(t, filepath.Join(analysistest.TestData(), "src", "versions", "go22.txtar"))
	analysistest.Run(t, dir, loopclosure.Analyzer, "golang.org/fake/versions")
}

func TestVersions18(t *testing.T) {
	t.Skip("Disabled for golang.org/cl/603895. Fix and re-enable.")
	dir := testfiles.ExtractTxtarFileToTmp(t, filepath.Join(analysistest.TestData(), "src", "versions", "go18.txtar"))
	analysistest.Run(t, dir, loopclosure.Analyzer, "golang.org/fake/versions")
}
