package main

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestCompileOpts_Arguments(t *testing.T) {
	testCases := []struct {
		name     string
		have     *CompileOpts
		expected []string
	}{
		{"ShouldNotReturnAnyArgs", &CompileOpts{}, nil},
		{"ShouldReturnBooleanTrimPath", &CompileOpts{TrimPath: true}, []string{"-trimpath"}},
		{"ShouldReturnBooleanRace", &CompileOpts{Race: true}, []string{"-race"}},
		{"ShouldReturnBooleanCover", &CompileOpts{Cover: true}, []string{"-cover"}},
		{"ShouldReturnBooleanCoverWithVersion", &CompileOpts{GoVersion: MustParseVersion("go1.20"), Cover: true}, []string{"-cover"}},
		{"ShouldReturnBooleanCoverWithBadVersion", &CompileOpts{GoVersion: MustParseVersion("go1.19"), Cover: true}, nil},
		{"ShouldReturnBooleanRebuild", &CompileOpts{Rebuild: true}, []string{"-a"}},
		{"ShouldReturnBooleanMSAN", &CompileOpts{MemorySanitizer: true}, []string{"-msan"}},
		{"ShouldReturnBooleanASAN", &CompileOpts{AddressSanitizer: true}, []string{"-asan"}},
		{"ShouldReturnBooleanLinkShared", &CompileOpts{LinkShared: true}, []string{"-linkshared"}},
		{"ShouldReturnBooleanLinkModCacheRW", &CompileOpts{ModCacheRW: true}, []string{"-modcacherw"}},
		{"ShouldReturnBooleanLinkModCacheRWWithVersion", &CompileOpts{GoVersion: MustParseVersion("go1.19"), ModCacheRW: true}, []string{"-modcacherw"}},
		{"ShouldReturnBooleanLinkModCacheRWWithBadVersion", &CompileOpts{GoVersion: MustParseVersion("go1.13"), ModCacheRW: true}, nil},
		{"ShouldReturnBooleanMiscMixed", &CompileOpts{GoVersion: MustParseVersion("go1.13"), ModCacheRW: true, TrimPath: true, Rebuild: true}, []string{"-a", "-trimpath"}},
		{"ShouldReturnStringC", &CompileOpts{GoVersion: MustParseVersion("go1.20"), ChangeDir: "foo"}, []string{"-C=foo"}},
		{"ShouldReturnStringCoverPackage", &CompileOpts{GoVersion: MustParseVersion("go1.20"), CoverPackage: "foo"}, []string{"-coverpkg=foo"}},
		{"ShouldReturnStringsASMFlags", &CompileOpts{GoVersion: MustParseVersion("go1.20"), ASMFlags: "foo"}, []string{`-asmflags="foo"`}},
		{"ShouldReturnStringBuildMode", &CompileOpts{GoVersion: MustParseVersion("go1.20"), BuildMode: "foo"}, []string{`-buildmode=foo`}},
		{"ShouldReturnStringBuildVCS", &CompileOpts{GoVersion: MustParseVersion("go1.20"), BuildVCS: "foo"}, []string{`-buildvcs=foo`}},
		{"ShouldReturnStringCompiler", &CompileOpts{GoVersion: MustParseVersion("go1.20"), Compiler: "foo"}, []string{`-compiler=foo`}},
		{"ShouldReturnStringsGcCGOFlags", &CompileOpts{GoVersion: MustParseVersion("go1.20"), GcCGOFlags: "foo"}, []string{`-gccgoflags="foo"`}},
		{"ShouldReturnStringsGcFlags", &CompileOpts{GoVersion: MustParseVersion("go1.20"), GcFlags: "foo"}, []string{`-gcflags="foo"`}},
		{"ShouldReturnStringInstallSuffix", &CompileOpts{GoVersion: MustParseVersion("go1.20"), InstallSuffix: "foo"}, []string{`-installsuffix=foo`}},
		{"ShouldReturnStringsLDFlags", &CompileOpts{GoVersion: MustParseVersion("go1.20"), LDFlags: "foo"}, []string{`-ldflags="foo"`}},
		{"ShouldReturnStringMod", &CompileOpts{GoVersion: MustParseVersion("go1.20"), Mod: "foo"}, []string{`-mod=foo`}},
		{"ShouldReturnStringModFile", &CompileOpts{GoVersion: MustParseVersion("go1.20"), ModFile: "foo"}, []string{`-modfile=foo`}},
		{"ShouldReturnStringOverlay", &CompileOpts{GoVersion: MustParseVersion("go1.20"), Overlay: "foo"}, []string{`-overlay=foo`}},
		{"ShouldReturnStringPGO", &CompileOpts{GoVersion: MustParseVersion("go1.20"), ProfileGuidedOptimization: "foo"}, []string{`-pgo=foo`}},
		{"ShouldReturnStringPkgDir", &CompileOpts{GoVersion: MustParseVersion("go1.20"), PackageDir: "foo"}, []string{`-pkgdir=foo`}},
		{"ShouldReturnStringsTags", &CompileOpts{GoVersion: MustParseVersion("go1.20"), Tags: "foo"}, []string{`-tags="foo"`}},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			assert.Equal(t, tc.expected, tc.have.Arguments())
		})
	}
}
