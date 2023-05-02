package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"os/exec"
	"runtime"
	"strings"
	"sync"

	"github.com/hashicorp/go-version"
)

func main() {
	// Call realMain so that defers work properly, since os.Exit won't
	// call defers.
	os.Exit(realMain())
}

func realMain() int {
	var (
		flagBuildToolchain, flagVerbose, flagCgo, flagTrimPath, flagListOSArch bool
		flagGoCmd, flagLDFlags, flagTags                                       string
		flagOutput                                                             string
		flagParallel                                                           int
		flagPlatform                                                           PlatformFlag

		flagChangeDir, flagCoverPkg, flagASMFlags, flagBuildMode, flagBuildVCS, flagCompiler, flagGcCGOFlags, flagGcFlags string
		flagInstallSuffix, flagMod, flagModFile, flagOverlay, flagProfileGuidedOptimization, flagPackageDir               string
		flagRebuild, flagRace, flagMSAN, flagASAN, flagCover, flagLinkShared, flagModCacheRW                              bool
	)

	flags := flag.NewFlagSet("gox", flag.ExitOnError)

	flags.Usage = func() { printUsage() }

	// Gox Flags.
	flags.BoolVar(&flagListOSArch, "osarch-list", false, "")
	flags.StringVar(&flagOutput, "output", "{{.Dir}}_{{.OS}}_{{.Arch}}", "")
	flags.IntVar(&flagParallel, "parallel", -1, "")
	flags.BoolVar(&flagVerbose, "verbose", false, "")
	flags.StringVar(&flagGoCmd, "gocmd", "go", "")

	// Misc.
	flags.Var(flagPlatform.ArchFlagValue(), "arch", "")
	flags.BoolVar(&flagBuildToolchain, "build-toolchain", false, "")
	flags.Var(flagPlatform.OSFlagValue(), "os", "")
	flags.Var(flagPlatform.OSArchFlagValue(), "osarch", "")

	// Env Flags.
	flags.BoolVar(&flagCgo, "cgo", false, "")

	// Build Flags.
	flags.StringVar(&flagChangeDir, "change-dir", "", "")
	flags.BoolVar(&flagRebuild, "rebuild", false, "")
	flags.BoolVar(&flagRace, flagNameRace, false, "")
	flags.BoolVar(&flagMSAN, flagNameMSAN, false, "")
	flags.BoolVar(&flagASAN, flagNameASAN, false, "")
	flags.BoolVar(&flagCover, flagNameCover, false, "")
	flags.StringVar(&flagCoverPkg, flagNameCoverPKG, "", "")
	flags.StringVar(&flagASMFlags, flagNameASMFlags, "", "")
	flags.StringVar(&flagBuildMode, flagNameBuildMode, "", "")
	flags.StringVar(&flagBuildVCS, flagNameBuildVCS, "", "")
	flags.StringVar(&flagCompiler, flagNameCompiler, "", "")
	flags.StringVar(&flagGcCGOFlags, flagNameGcCGOFlags, "", "")
	flags.StringVar(&flagGcFlags, flagNameGcFlags, "", "")
	flags.StringVar(&flagInstallSuffix, flagNameInstallSuffix, "", "")
	flags.StringVar(&flagLDFlags, flagNameLDFlags, "", "")
	flags.BoolVar(&flagLinkShared, flagNameLinkShared, false, "")
	flags.StringVar(&flagMod, flagNameMod, "", "")
	flags.BoolVar(&flagModCacheRW, flagNameModCacheRW, false, "")
	flags.StringVar(&flagModFile, flagNameModFile, "", "")
	flags.StringVar(&flagOverlay, flagNameOverlay, "", "")
	flags.StringVar(&flagProfileGuidedOptimization, flagNameProfileGuidedOptimization, "", "")
	flags.StringVar(&flagPackageDir, flagNamePackageDir, "", "")
	flags.StringVar(&flagTags, flagNameTags, "", "")
	flags.BoolVar(&flagTrimPath, flagNameTrimPath, false, "")

	if err := flags.Parse(os.Args[1:]); err != nil {
		flags.Usage()
		return 1
	}

	// Determine what amount of parallelism we want Default to the current
	// number of CPUs-1 is <= 0 is specified.
	if flagParallel <= 0 {
		cpus := runtime.NumCPU()
		if cpus < 2 {
			flagParallel = 1
		} else {
			flagParallel = cpus - 1
		}

		// Joyent containers report 48 cores via runtime.NumCPU(), and a
		// default of 47 parallel builds causes a panic. Default to 3 on
		// Solaris-derived operating systems unless overridden with the
		// -parallel flag.
		if runtime.GOOS == goosSolaris {
			flagParallel = 3
		}
	}

	versionStr, err := GoVersion()
	if err != nil {
		fmt.Fprintf(os.Stderr, "error reading Go version: %s", err)
		return 1
	}

	var v *version.Version

	// Use latest if we get an unexpected version string
	if strings.HasPrefix(versionStr, "go") {
		if v, err = version.NewVersion(versionStr[2:]); err != nil {
			log.Printf("Unable to parse current go version: %s\n%s", versionStr, err.Error())
		}
	}

	if flagBuildToolchain {
		return mainBuildToolchain(v, flagParallel, flagPlatform, flagVerbose)
	}

	if _, err := exec.LookPath(flagGoCmd); err != nil {
		fmt.Fprintf(os.Stderr, "%s executable must be on the PATH\n",
			flagGoCmd)
		return 1
	}

	if flagListOSArch {
		return mainListOSArch(v)
	}

	// Determine the packages that we want to compile. Default to the
	// current directory if none are specified.
	packages := flags.Args()
	if len(packages) == 0 {
		packages = []string{"."}
	}

	// Get the packages that are in the given paths
	mainDirs, err := GoMainDirs(packages, flagGoCmd)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error reading packages: %s", err)
		return 1
	}

	// Determine the platforms we're building for
	platforms := flagPlatform.Platforms(SupportedPlatforms(v))
	if len(platforms) == 0 {
		fmt.Println("No valid platforms to build for. If you specified a value")
		fmt.Println("for the 'os', 'arch', or 'osarch' flags, make sure you're")
		fmt.Println("using a valid value.")
		return 1
	}

	// Assume -mod is supported when no version prefix is found
	if flagMod != "" {
		constraint, err := version.NewConstraint(">= 1.11")
		if err != nil {
			panic(err)
		}

		if !constraint.Check(v) {
			fmt.Printf("Go compiler version %s does not support the -mod flag\n", versionStr)
			flagMod = ""
		}
	}

	// Build in parallel!
	fmt.Printf("Number of parallel builds: %d\n\n", flagParallel)

	var (
		errorLock sync.Mutex
		wg        sync.WaitGroup
	)

	errors := make([]string, 0)
	semaphore := make(chan int, flagParallel)

	for _, platform := range platforms {
		for _, path := range mainDirs {
			// Start the goroutine that will do the actual build
			wg.Add(1)
			go func(path string, platform Platform) {
				defer wg.Done()
				semaphore <- 1
				fmt.Printf("--> %15s: %s\n", platform.String(), path)

				opts := &CompileOpts{
					GoVersion:   v,
					PackagePath: path,
					Platform:    platform,
					OutputTpl:   flagOutput,
					Cgo:         flagCgo,
					GoCmd:       flagGoCmd,

					ChangeDir:                 flagChangeDir,
					Rebuild:                   flagRebuild,
					Race:                      flagRace,
					MemorySanitizer:           flagMSAN,
					AddressSanitizer:          flagASAN,
					Cover:                     flagCover,
					CoverPackage:              flagCoverPkg,
					ASMFlags:                  flagASMFlags,
					BuildMode:                 flagBuildMode,
					BuildVCS:                  flagBuildVCS,
					Compiler:                  flagCompiler,
					GcCGOFlags:                flagGcCGOFlags,
					GcFlags:                   flagGcFlags,
					InstallSuffix:             flagInstallSuffix,
					LDFlags:                   flagLDFlags,
					LinkShared:                flagLinkShared,
					ModMode:                   flagMod,
					ModCacheRW:                flagModCacheRW,
					ModFile:                   flagModFile,
					Overlay:                   flagOverlay,
					ProfileGuidedOptimization: flagProfileGuidedOptimization,
					PackageDir:                flagPackageDir,
					Tags:                      flagTags,
					TrimPath:                  flagTrimPath,
				}

				// Determine if we have specific CFLAGS or LDFLAGS for this
				// GOOS/GOARCH combo and override the defaults if so.
				envOverride(&opts.LDFlags, platform, "LDFLAGS")
				envOverride(&opts.GcFlags, platform, "GCFLAGS")
				envOverride(&opts.ASMFlags, platform, "ASMFLAGS")
				envOverride(&opts.Cc, platform, "CC")
				envOverride(&opts.Cxx, platform, "CXX")

				if err := GoCrossCompile(opts); err != nil {
					errorLock.Lock()
					defer errorLock.Unlock()
					errors = append(errors,
						fmt.Sprintf("%s error: %s", platform.String(), err))
				}
				<-semaphore
			}(path, platform)
		}
	}
	wg.Wait()

	if len(errors) > 0 {
		fmt.Fprintf(os.Stderr, "\n%d errors occurred:\n", len(errors))
		for _, err := range errors {
			fmt.Fprintf(os.Stderr, "--> %s\n", err)
		}
		return 1
	}

	return 0
}

func printUsage() {
	fmt.Fprintf(os.Stderr, helpText, metaVersion)
}

const helpText = `Usage: gox [options] [packages]

  Gox cross-compiles Go applications in parallel.

  Version: %s

  If no specific operating systems or architectures are specified, Gox
  will build for all pairs supported by your version of Go.

Options:

  -arch=""            Space-separated list of architectures to build for
  -build-toolchain    Build cross-compilation toolchain
  -cgo                Sets CGO_ENABLED=1, requires proper C toolchain (advanced)
  -os=""              Space-separated list of operating systems to build for
  -osarch=""          Space-separated list of os/arch pairs to build for
  -osarch-list        List supported os/arch pairs for your Go version
  -output="foo"       Output path template. See below for more info
  -parallel=-1        Amount of parallelism, defaults to number of CPUs
  -gocmd="go"         Build command, defaults to Go
  -verbose            Verbose mode

  -change-dir=""      Passed to go build flag '-C'. Change to dir before running the command
  -rebuild            Passed to go build flag '-a'. Force rebuilding of package that were up to date
  -race               Passed to go build flag '-race'. Build with the go race detector enabled, requires CGO
  -msan               Passed to go build flag '-msan'. Enable interoperation with memory sanitizer
  -asan               Passed to go build flag '-asan'. Enable interoperation with address sanitizer
  -cover              Passed to go build flag '-cover'. Enable code coverage instrumentation
  -coverpkg=""        Passed to go build flag '-coverpkg'. Apply coverage analysis to each package matching the patterns
  -asmflags=""        Passed to go build flag '-asmflags'.
  -buildmode=""       Passed to go build flag '-buildmode'.
  -buildvcs=""        Passed to go build flag '-buildvcs'.
  -compiler=""        Passed to go build flag '-compiler'.
  -gccgoflags=""      Passed to go build flag '-gccgoflags'.
  -gcflags=""         Passed to go build flag '-gcflags'.
  -installsuffix=""   Passed to go build flag '-installsuffix'.
  -ldflags=""         Passed to go build flag '-ldflags'.
  -linkshared         Passed to go build flag '-linkshared'.
  -mod=""             Passed to go build flag '-mod'.
  -modcacherw         Passed to go build flag '-modcacherw'.
  -modfile=""         Passed to go build flag '-modfile'.
  -overlay=""         Passed to go build flag '-overlay'.
  -pgo=""             Passed to go build flag '-pgo'.
  -pkgdir=""          Passed to go build flag '-pkgdir'.
  -tags=""            Passed to go build flag '-tags'.
  -trimpath           Passed to go build flag '-trimpath'. Remove all file system paths from the resulting executable

Output path template:

  The output path for the compiled binaries is specified with the
  "-output" flag. The value is a string that is a Go text template.
  The default value is "{{.Dir}}_{{.OS}}_{{.Arch}}". The variables and
  their values should be self-explanatory.

Platforms (OS/Arch):

  The operating systems and architectures to cross-compile for may be
  specified with the "-arch" and "-os" flags. These are space separated lists
  of valid GOOS/GOARCH values to build for, respectively. You may prefix an
  OS or Arch with "!" to negate and not build for that platform. If the list
  is made up of only negations, then the negations will come from the default
  list.

  Additionally, the "-osarch" flag may be used to specify complete os/arch
  pairs that should be built or ignored. The syntax for this is what you would
  expect: "darwin/amd64" would be a valid osarch value. Multiple can be space
  separated. An os/arch pair can begin with "!" to not build for that platform.

  The "-osarch" flag has the highest precedent when determining whether to
  build for a platform. If it is included in the "-osarch" list, it will be
  built even if the specific os and arch is negated in "-os" and "-arch",
  respectively.

Platform Overrides:

  The "-gcflags", "-ldflags" and "-asmflags" options and "CC"/"CXX" environment
  variables for cross-compilation can be overridden per-platform by using
  environment variables. Gox will look for environment variables in the
  following format and use those to override values if they exist:

    GOX_[OS]_[ARCH]_GCFLAGS
    GOX_[OS]_[ARCH]_LDFLAGS
    GOX_[OS]_[ARCH]_ASMFLAGS
    GOX_[OS]_[ARCH]_CC
    GOX_[OS]_[ARCH]_CXX
`
