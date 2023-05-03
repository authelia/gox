package main

import (
	"fmt"
	"github.com/hashicorp/go-version"
	"os"
	"strings"
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

	ChangeDir string // -C
	Rebuild   bool   // -a
	Race      bool   // -race

	MemorySanitizer  bool // -msan
	AddressSanitizer bool // -asan

	Cover         bool   // -cover
	CoverPackage  string // -coverpkg
	ASMFlags      string // -asmflags
	BuildMode     string // -buildmode
	BuildVCS      string // -buildvs
	Compiler      string // -compiler
	GcCGOFlags    string // -gccgoflags
	GcFlags       string // -gcflags
	InstallSuffix string // -installsuffix
	LDFlags       string // -ldflags
	LinkShared    bool   // -linkshared
	Mod           string // -mod
	ModCacheRW    bool   // -modcacherw
	ModFile       string // -modfile
	Overlay       string // -overlay

	ProfileGuidedOptimization string // -pgo

	PackageDir string // -pkgdir
	Tags       string // -tags
	TrimPath   bool   // -trimpath
}

func (opts *CompileOpts) Experiments() (exp []string) {
	exp = strings.Split(os.Getenv("GOEXPERIMENT"), ",")

	if opts.Cover && (opts.GoVersion == nil || flagConstraints[flagNameCover].Check(opts.GoVersion)) {
		if constraints, ok := experimentConstraints[flagNameCover]; !ok || constraints.Check(opts.GoVersion) {
			exp = append(exp, "coverageredesign")
		}
	}

	return exp
}

func (opts *CompileOpts) Env() []string {
	env := append(os.Environ(),
		"GOOS="+opts.Platform.OS,
		"GOARCH="+opts.Platform.Arch)

	if opts.Cc != "" {
		env = append(env, "CC="+opts.Cc)
	}

	if opts.Cxx != "" {
		env = append(env, "CXX="+opts.Cxx)
	}

	if opts.Cgo {
		env = append(env, "CGO_ENABLED=1")
	} else {
		env = append(env, "CGO_ENABLED=0")
	}

	if exp := opts.Experiments(); len(exp) != 0 {
		env = append(env, fmt.Sprintf("GOEXPERIMENT=%s", strings.Join(exp, ",")))
	}

	return env
}

func (opts *CompileOpts) Arguments() []string {
	arguments := &BuildArgFlags{}

	arguments.AddString(opts.GoVersion, "C", opts.ChangeDir)
	arguments.AddBoolean(opts.GoVersion, "a", opts.Rebuild)
	arguments.AddBoolean(opts.GoVersion, flagNameRace, opts.Race)
	arguments.AddBoolean(opts.GoVersion, flagNameMSAN, opts.MemorySanitizer)
	arguments.AddBoolean(opts.GoVersion, flagNameASAN, opts.AddressSanitizer)
	arguments.AddBoolean(opts.GoVersion, flagNameCover, opts.Cover)
	arguments.AddString(opts.GoVersion, flagNameCoverPKG, opts.CoverPackage)
	arguments.AddStrings(opts.GoVersion, flagNameASMFlags, opts.ASMFlags)
	arguments.AddString(opts.GoVersion, flagNameBuildMode, opts.BuildMode)
	arguments.AddString(opts.GoVersion, flagNameBuildVCS, opts.BuildVCS)
	arguments.AddString(opts.GoVersion, flagNameCompiler, opts.Compiler)
	arguments.AddStrings(opts.GoVersion, flagNameGcCGOFlags, opts.GcCGOFlags)
	arguments.AddStrings(opts.GoVersion, flagNameGcFlags, opts.GcFlags)
	arguments.AddString(opts.GoVersion, flagNameInstallSuffix, opts.InstallSuffix)
	arguments.AddStrings(opts.GoVersion, flagNameLDFlags, opts.LDFlags)
	arguments.AddBoolean(opts.GoVersion, flagNameLinkShared, opts.LinkShared)
	arguments.AddString(opts.GoVersion, flagNameMod, opts.Mod)
	arguments.AddBoolean(opts.GoVersion, flagNameModCacheRW, opts.ModCacheRW)
	arguments.AddString(opts.GoVersion, flagNameModFile, opts.ModFile)
	arguments.AddString(opts.GoVersion, flagNameOverlay, opts.Overlay)
	arguments.AddString(opts.GoVersion, flagNameProfileGuidedOptimization, opts.ProfileGuidedOptimization)
	arguments.AddString(opts.GoVersion, flagNamePackageDir, opts.PackageDir)
	arguments.AddStrings(opts.GoVersion, flagNameTags, opts.Tags)
	arguments.AddBoolean(opts.GoVersion, flagNameTrimPath, opts.TrimPath)

	return arguments.Args()
}

type BuildArgFlags struct {
	args []string
}

func (f *BuildArgFlags) Args() []string {
	return f.args
}

func (f *BuildArgFlags) ShouldSkipVersion(v *version.Version, name string) bool {
	if v == nil {
		return false
	}

	if constraints, ok := flagConstraints[name]; ok {
		if !constraints.Check(v) {
			return true
		}
	}

	return false
}

func (f *BuildArgFlags) AddStrings(v *version.Version, name, value string) {
	if value == "" {
		return
	}

	if f.ShouldSkipVersion(v, name) {
		return
	}

	f.args = append(f.args, fmt.Sprintf(`-%s="%s"`, name, value))
}

func (f *BuildArgFlags) AddString(v *version.Version, name, value string) {
	if value == "" {
		return
	}

	if f.ShouldSkipVersion(v, name) {
		return
	}

	f.args = append(f.args, fmt.Sprintf("-%s=%s", name, value))
}

func (f *BuildArgFlags) AddBoolean(v *version.Version, name string, value bool) {
	if !value {
		return
	}

	if f.ShouldSkipVersion(v, name) {
		return
	}

	f.args = append(f.args, "-"+name)
}
