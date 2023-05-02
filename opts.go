package main

import (
	"github.com/hashicorp/go-version"
)

type CompileOpts struct {
	PackagePath string
	Platform    Platform
	OutputTpl   string
	Cc          string
	Cxx         string
	Cgo         bool
	GoCmd       string
	GoVersion   *version.Version

	ChangeDir string // -C (todo)
	Rebuild   bool   // -a
	Race      bool   // -race

	MemorySanitizer  bool // -msan (todo)
	AddressSanitizer bool // -asan (todo)

	Cover         bool   // -cover (todo) MUST SET GOEXPERIMENT=coverageredesign
	CoverPackage  string // -coverpkg (todo)
	ASMFlags      string // -asmflags
	BuildMode     string // -buildmode
	BuildVCS      string // -buildvs
	Compiler      string // -compiler (todo)
	GcCGOFlags    string // -gccgoflags (todo)
	GcFlags       string // -gcflags
	InstallSuffix string // -installsuffix (todo)
	LDFlags       string // -ldflags
	LinkShared    bool   // -linkshared (todo)
	ModMode       string // -mod
	ModCacheRW    bool   // -modcacherw (todo)
	ModFile       string // -modfile (todo)
	Overlay       string // -overlay (todo)

	ProfileGuidedOptimization string // -pgo (todo)

	PackageDir string // -pkgdir (todo)
	Tags       string // -tags
	TrimPath   bool   // -trimpath
}

type BuildArgFlags struct {
	args []string
}

func (f *BuildArgFlags) Add(v *version.Version, name, value string, formatter func(string, string) string) {

}

var flagConstraints map[string]version.Constraints

func init() {
	inputs := []struct {
		name       string
		constraint string
	}{
		{"C", ">= 1.20"},
		{"asan", ">= 1.18"},
		{"cover", ">= 1.20"},
		{"coverpkg", ">= 1.20"},
		{"buildvcs", ">= 1.18"},
		{"modcacherw", ">= 1.14"},
		{"modfile", ">= 1.14"},
		{"overlay", ">= 1.16"},
		{"pgo", ">= 1.20"},
		{"trimpath", ">= 1.13"},
	}

	var (
		constraint version.Constraints
		err        error
	)

	flagConstraints = map[string]version.Constraints{}

	for _, input := range inputs {
		if constraint, err = version.NewConstraint(input.constraint); err != nil {
			panic(err)
		}

		flagConstraints[input.name] = constraint
	}
}
