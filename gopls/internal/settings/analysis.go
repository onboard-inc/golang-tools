// Copyright 2024 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package settings

import (
	"github.com/onboard-inc/golang-tools/go/analysis"
	"github.com/onboard-inc/golang-tools/go/analysis/passes/appends"
	"github.com/onboard-inc/golang-tools/go/analysis/passes/asmdecl"
	"github.com/onboard-inc/golang-tools/go/analysis/passes/assign"
	"github.com/onboard-inc/golang-tools/go/analysis/passes/atomic"
	"github.com/onboard-inc/golang-tools/go/analysis/passes/atomicalign"
	"github.com/onboard-inc/golang-tools/go/analysis/passes/bools"
	"github.com/onboard-inc/golang-tools/go/analysis/passes/buildtag"
	"github.com/onboard-inc/golang-tools/go/analysis/passes/cgocall"
	"github.com/onboard-inc/golang-tools/go/analysis/passes/composite"
	"github.com/onboard-inc/golang-tools/go/analysis/passes/copylock"
	"github.com/onboard-inc/golang-tools/go/analysis/passes/deepequalerrors"
	"github.com/onboard-inc/golang-tools/go/analysis/passes/defers"
	"github.com/onboard-inc/golang-tools/go/analysis/passes/directive"
	"github.com/onboard-inc/golang-tools/go/analysis/passes/errorsas"
	"github.com/onboard-inc/golang-tools/go/analysis/passes/framepointer"
	"github.com/onboard-inc/golang-tools/go/analysis/passes/httpresponse"
	"github.com/onboard-inc/golang-tools/go/analysis/passes/ifaceassert"
	"github.com/onboard-inc/golang-tools/go/analysis/passes/loopclosure"
	"github.com/onboard-inc/golang-tools/go/analysis/passes/lostcancel"
	"github.com/onboard-inc/golang-tools/go/analysis/passes/nilfunc"
	"github.com/onboard-inc/golang-tools/go/analysis/passes/nilness"
	"github.com/onboard-inc/golang-tools/go/analysis/passes/printf"
	"github.com/onboard-inc/golang-tools/go/analysis/passes/shadow"
	"github.com/onboard-inc/golang-tools/go/analysis/passes/shift"
	"github.com/onboard-inc/golang-tools/go/analysis/passes/sigchanyzer"
	"github.com/onboard-inc/golang-tools/go/analysis/passes/slog"
	"github.com/onboard-inc/golang-tools/go/analysis/passes/sortslice"
	"github.com/onboard-inc/golang-tools/go/analysis/passes/stdmethods"
	"github.com/onboard-inc/golang-tools/go/analysis/passes/stdversion"
	"github.com/onboard-inc/golang-tools/go/analysis/passes/stringintconv"
	"github.com/onboard-inc/golang-tools/go/analysis/passes/structtag"
	"github.com/onboard-inc/golang-tools/go/analysis/passes/testinggoroutine"
	"github.com/onboard-inc/golang-tools/go/analysis/passes/tests"
	"github.com/onboard-inc/golang-tools/go/analysis/passes/timeformat"
	"github.com/onboard-inc/golang-tools/go/analysis/passes/unmarshal"
	"github.com/onboard-inc/golang-tools/go/analysis/passes/unreachable"
	"github.com/onboard-inc/golang-tools/go/analysis/passes/unsafeptr"
	"github.com/onboard-inc/golang-tools/go/analysis/passes/unusedresult"
	"github.com/onboard-inc/golang-tools/go/analysis/passes/unusedwrite"
	"github.com/onboard-inc/golang-tools/gopls/internal/analysis/deprecated"
	"github.com/onboard-inc/golang-tools/gopls/internal/analysis/embeddirective"
	"github.com/onboard-inc/golang-tools/gopls/internal/analysis/fillreturns"
	"github.com/onboard-inc/golang-tools/gopls/internal/analysis/infertypeargs"
	"github.com/onboard-inc/golang-tools/gopls/internal/analysis/nonewvars"
	"github.com/onboard-inc/golang-tools/gopls/internal/analysis/norangeoverfunc"
	"github.com/onboard-inc/golang-tools/gopls/internal/analysis/noresultvalues"
	"github.com/onboard-inc/golang-tools/gopls/internal/analysis/simplifycompositelit"
	"github.com/onboard-inc/golang-tools/gopls/internal/analysis/simplifyrange"
	"github.com/onboard-inc/golang-tools/gopls/internal/analysis/simplifyslice"
	"github.com/onboard-inc/golang-tools/gopls/internal/analysis/stubmethods"
	"github.com/onboard-inc/golang-tools/gopls/internal/analysis/undeclaredname"
	"github.com/onboard-inc/golang-tools/gopls/internal/analysis/unusedparams"
	"github.com/onboard-inc/golang-tools/gopls/internal/analysis/unusedvariable"
	"github.com/onboard-inc/golang-tools/gopls/internal/analysis/useany"
	"github.com/onboard-inc/golang-tools/gopls/internal/protocol"
	"honnef.co/go/tools/staticcheck"
)

// Analyzer augments a [analysis.Analyzer] with additional LSP configuration.
//
// Analyzers are immutable, since they are shared across multiple LSP sessions.
type Analyzer struct {
	analyzer    *analysis.Analyzer
	enabled     bool
	actionKinds []protocol.CodeActionKind
	severity    protocol.DiagnosticSeverity
	tags        []protocol.DiagnosticTag
}

// Analyzer returns the [analysis.Analyzer] that this Analyzer wraps.
func (a *Analyzer) Analyzer() *analysis.Analyzer { return a.analyzer }

// EnabledByDefault reports whether the analyzer is enabled by default for all sessions.
// This value can be configured per-analysis in user settings.
func (a *Analyzer) EnabledByDefault() bool { return a.enabled }

// ActionKinds is the set of kinds of code action this analyzer produces.
//
// If left unset, it defaults to QuickFix.
// TODO(rfindley): revisit.
func (a *Analyzer) ActionKinds() []protocol.CodeActionKind { return a.actionKinds }

// Severity is the severity set for diagnostics reported by this
// analyzer. If left unset it defaults to Warning.
//
// Note: diagnostics with severity protocol.SeverityHint do not show up in
// the VS Code "problems" tab.
func (a *Analyzer) Severity() protocol.DiagnosticSeverity { return a.severity }

// Tags is extra tags (unnecessary, deprecated, etc) for diagnostics
// reported by this analyzer.
func (a *Analyzer) Tags() []protocol.DiagnosticTag { return a.tags }

// String returns the name of this analyzer.
func (a *Analyzer) String() string { return a.analyzer.String() }

// DefaultAnalyzers holds the set of Analyzers available to all gopls sessions,
// independent of build version, keyed by analyzer name.
//
// It is the source from which gopls/doc/analyzers.md is generated.
var DefaultAnalyzers = make(map[string]*Analyzer) // initialized below

func init() {
	// Emergency workaround for #67237 to allow standard library
	// to use range over func: disable SSA-based analyses of
	// go1.23 packages that use range-over-func.
	suppressOnRangeOverFunc := func(a *analysis.Analyzer) {
		a.Requires = append(a.Requires, norangeoverfunc.Analyzer)
	}
	// buildir is non-exported so we have to scan the Analysis.Requires graph to find it.
	var buildir *analysis.Analyzer
	for _, a := range staticcheck.Analyzers {
		for _, req := range a.Analyzer.Requires {
			if req.Name == "buildir" {
				buildir = req
			}
		}

		// Temporarily disable SA4004 CheckIneffectiveLoop as
		// it crashes when encountering go1.23 range-over-func
		// (#67237, dominikh/go-tools#1494).
		if a.Analyzer.Name == "SA4004" {
			suppressOnRangeOverFunc(a.Analyzer)
		}
	}
	if buildir != nil {
		suppressOnRangeOverFunc(buildir)
	}

	analyzers := []*Analyzer{
		// The traditional vet suite:
		{analyzer: appends.Analyzer, enabled: true},
		{analyzer: asmdecl.Analyzer, enabled: true},
		{analyzer: assign.Analyzer, enabled: true},
		{analyzer: atomic.Analyzer, enabled: true},
		{analyzer: bools.Analyzer, enabled: true},
		{analyzer: buildtag.Analyzer, enabled: true},
		{analyzer: cgocall.Analyzer, enabled: true},
		{analyzer: composite.Analyzer, enabled: true},
		{analyzer: copylock.Analyzer, enabled: true},
		{analyzer: defers.Analyzer, enabled: true},
		{analyzer: deprecated.Analyzer, enabled: true, severity: protocol.SeverityHint, tags: []protocol.DiagnosticTag{protocol.Deprecated}},
		{analyzer: directive.Analyzer, enabled: true},
		{analyzer: errorsas.Analyzer, enabled: true},
		{analyzer: framepointer.Analyzer, enabled: true},
		{analyzer: httpresponse.Analyzer, enabled: true},
		{analyzer: ifaceassert.Analyzer, enabled: true},
		{analyzer: loopclosure.Analyzer, enabled: true},
		{analyzer: lostcancel.Analyzer, enabled: true},
		{analyzer: nilfunc.Analyzer, enabled: true},
		{analyzer: printf.Analyzer, enabled: true},
		{analyzer: shift.Analyzer, enabled: true},
		{analyzer: sigchanyzer.Analyzer, enabled: true},
		{analyzer: slog.Analyzer, enabled: true},
		{analyzer: stdmethods.Analyzer, enabled: true},
		{analyzer: stdversion.Analyzer, enabled: true},
		{analyzer: stringintconv.Analyzer, enabled: true},
		{analyzer: structtag.Analyzer, enabled: true},
		{analyzer: testinggoroutine.Analyzer, enabled: true},
		{analyzer: tests.Analyzer, enabled: true},
		{analyzer: timeformat.Analyzer, enabled: true},
		{analyzer: unmarshal.Analyzer, enabled: true},
		{analyzer: unreachable.Analyzer, enabled: true},
		{analyzer: unsafeptr.Analyzer, enabled: true},
		{analyzer: unusedresult.Analyzer, enabled: true},

		// not suitable for vet:
		// - some (nilness) use go/ssa; see #59714.
		// - others don't meet the "frequency" criterion;
		//   see GOROOT/src/cmd/vet/README.
		{analyzer: atomicalign.Analyzer, enabled: true},
		{analyzer: deepequalerrors.Analyzer, enabled: true},
		{analyzer: nilness.Analyzer, enabled: true}, // uses go/ssa
		{analyzer: sortslice.Analyzer, enabled: true},
		{analyzer: embeddirective.Analyzer, enabled: true},

		// disabled due to high false positives
		{analyzer: shadow.Analyzer, enabled: false}, // very noisy
		{analyzer: useany.Analyzer, enabled: false}, // never a bug
		// fieldalignment is not even off-by-default; see #67762.

		// "simplifiers": analyzers that offer mere style fixes
		// gofmt -s suite:
		{analyzer: simplifycompositelit.Analyzer, enabled: true, actionKinds: []protocol.CodeActionKind{protocol.SourceFixAll, protocol.QuickFix}},
		{analyzer: simplifyrange.Analyzer, enabled: true, actionKinds: []protocol.CodeActionKind{protocol.SourceFixAll, protocol.QuickFix}},
		{analyzer: simplifyslice.Analyzer, enabled: true, actionKinds: []protocol.CodeActionKind{protocol.SourceFixAll, protocol.QuickFix}},
		// other simplifiers:
		{analyzer: infertypeargs.Analyzer, enabled: true, severity: protocol.SeverityHint},
		{analyzer: unusedparams.Analyzer, enabled: true},
		{analyzer: unusedwrite.Analyzer, enabled: true}, // uses go/ssa

		// type-error analyzers
		// These analyzers enrich go/types errors with suggested fixes.
		{analyzer: fillreturns.Analyzer, enabled: true},
		{analyzer: nonewvars.Analyzer, enabled: true},
		{analyzer: noresultvalues.Analyzer, enabled: true},
		{analyzer: stubmethods.Analyzer, enabled: true},
		{analyzer: undeclaredname.Analyzer, enabled: true},
		// TODO(rfindley): why isn't the 'unusedvariable' analyzer enabled, if it
		// is only enhancing type errors with suggested fixes?
		//
		// In particular, enabling this analyzer could cause unused variables to be
		// greyed out, (due to the 'deletions only' fix). That seems like a nice UI
		// feature.
		{analyzer: unusedvariable.Analyzer, enabled: false},
	}
	for _, analyzer := range analyzers {
		DefaultAnalyzers[analyzer.analyzer.Name] = analyzer
	}
}

// StaticcheckAnalzyers describes available Staticcheck analyzers, keyed by
// analyzer name.
var StaticcheckAnalyzers = make(map[string]*Analyzer) // written by analysis_<ver>.go
