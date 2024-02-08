package main

const (
	gobin = "go"
)

const (
	flagNameRace                      = "race"
	flagNameMSAN                      = "msan"
	flagNameASAN                      = "asan"
	flagNameCover                     = "cover"
	flagNameCoverPKG                  = "coverpkg"
	flagNameASMFlags                  = "asmflags"
	flagNameBuildMode                 = "buildmode"
	flagNameBuildVCS                  = "buildvcs"
	flagNameCompiler                  = "compiler"
	flagNameGcCGOFlags                = "gccgoflags"
	flagNameGcFlags                   = "gcflags"
	flagNameInstallSuffix             = "installsuffix"
	flagNameLDFlags                   = "ldflags"
	flagNameLinkShared                = "linkshared"
	flagNameMod                       = "mod"
	flagNameModCacheRW                = "modcacherw"
	flagNameModFile                   = "modfile"
	flagNameOverlay                   = "overlay"
	flagNameProfileGuidedOptimization = "pgo"
	flagNamePackageDir                = "pkgdir"
	flagNameTags                      = "tags"
	flagNameTrimPath                  = "trimpath"
)

// Go Operating System values for GOOS.
const (
	goosWindows    = "windows"
	goosLinux      = "linux"
	goosDarwin     = "darwin"
	goosIOS        = "ios"
	goosFreeBSD    = "freebsd"
	goosNetBSD     = "netbsd"
	goosOpenBSD    = "openbsd"
	goosAndroid    = "android"
	goosPlan9      = "plan9"
	goosDragonfly  = "dragonfly"
	goosNACL       = "nacl"
	goosSolaris    = "solaris"
	goosJavaScript = "js"
	goosAIX        = "aix"
	goosIllumos    = "illumos"
	goosWASIP1     = "wasip1"
)

// Go Architecture values for GOARCH.
const (
	goarch386         = "386"
	goarchAMD64       = "amd64"
	goarchAMD64P32    = "amd64p32"
	goarchARM         = "arm"
	goarchARM64       = "arm64"
	goarchMIPS        = "mips"
	goarchMIPSLE      = "mipsle"
	goarchMIPS64      = "mips64"
	goarchMIPS64LE    = "mips64le"
	goarchS390X       = "s390x"
	goarchPowerPC64   = "ppc64"
	goarchPowerPC64LE = "ppc64le"
	goarchRISCV64     = "riscv64"
	goarchLoong64     = "loong64"
	goarchWebAssembly = "wasm"
)
