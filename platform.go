package main

import (
	"fmt"

	"github.com/hashicorp/go-version"
)

// Platform is a combination of OS/arch that can be built against.
type Platform struct {
	OS   string
	Arch string

	// Default, if true, will be included as a default build target
	// if no OS/arch is specified. We try to only set as a default popular
	// targets or targets that are generally useful. For example, Android
	// is not a default because it is quite rare that you're cross-compiling
	// something to Android AND something like Linux.
	Default bool
}

func (p *Platform) String() string {
	return fmt.Sprintf("%s/%s", p.OS, p.Arch)
}

// addDrop appends all of the "add" entries and drops the "drop" entries, ignoring
// the "Default" parameter.
func addDrop(base []Platform, add []Platform, drop []Platform) []Platform {
	newPlatforms := make([]Platform, len(base)+len(add))
	copy(newPlatforms, base)
	copy(newPlatforms[len(base):], add)

	// slow, but we only do this during initialization at most once per version
	for _, platform := range drop {
		found := -1
		for i := range newPlatforms {
			if newPlatforms[i].Arch == platform.Arch && newPlatforms[i].OS == platform.OS {
				found = i
				break
			}
		}
		if found < 0 {
			panic(fmt.Sprintf("Expected to remove %+v but not found in list %+v", platform, newPlatforms))
		}

		switch found {
		case 0:
			newPlatforms = newPlatforms[1:]
		case len(newPlatforms) - 1:
			newPlatforms = newPlatforms[:found]
		default:
			newPlatforms = append(newPlatforms[:found], newPlatforms[found+1:]...)
		}
	}
	return newPlatforms
}

var (
	Platforms_1_0 = []Platform{
		{goosDarwin, goarch386, true},
		{goosDarwin, goarchAMD64, true},
		{goosLinux, goarch386, true},
		{goosLinux, goarchAMD64, true},
		{goosLinux, goarchARM, true},
		{goosFreeBSD, goarch386, true},
		{goosFreeBSD, goarchAMD64, true},
		{goosOpenBSD, goarch386, true},
		{goosOpenBSD, goarchAMD64, true},
		{goosWindows, goarch386, true},
		{goosWindows, goarchAMD64, true},
	}

	Platforms_1_1 = addDrop(Platforms_1_0,
		[]Platform{
			{goosFreeBSD, goarchARM, true},
			{goosNetBSD, goarch386, true},
			{goosNetBSD, goarchAMD64, true},
			{goosNetBSD, goarchARM, true},
			{goosPlan9, goarch386, false},
		},
		nil,
	)

	// no new platforms in 1.2
	Platforms_1_2 = Platforms_1_1

	Platforms_1_3 = addDrop(Platforms_1_2,
		[]Platform{
			{goosDragonfly, goarch386, false},
			{goosDragonfly, goarchAMD64, false},
			{goosNACL, goarchAMD64, false},
			{goosNACL, goarchAMD64P32, false},
			{goosNACL, goarchARM, false},
			{goosSolaris, goarchAMD64, false},
		},
		nil,
	)

	Platforms_1_4 = addDrop(Platforms_1_3,
		[]Platform{
			{goosAndroid, goarchARM, false},
			{goosPlan9, goarchAMD64, false},
		},
		nil,
	)

	Platforms_1_5 = addDrop(Platforms_1_4,
		[]Platform{
			{goosDarwin, goarchARM, false},
			{goosDarwin, goarchARM64, false},
			{goosLinux, goarchARM64, false},
			{goosLinux, goarchPowerPC64, false},
			{goosLinux, goarchPowerPC64LE, false},
		},
		nil,
	)

	Platforms_1_6 = addDrop(Platforms_1_5,
		[]Platform{
			{goosAndroid, goarch386, false},
			{goosAndroid, goarchAMD64, false},
			{goosLinux, goarchMIPS64, false},
			{goosLinux, goarchMIPS64LE, false},
			{goosNACL, goarch386, false},
			{goosOpenBSD, goarchARM, true},
		},
		nil,
	)

	Platforms_1_7 = addDrop(Platforms_1_5,
		[]Platform{
			// Add the 1.6 Platforms, but reflect full support for mips64 and mips64le
			{goosAndroid, goarch386, false},
			{goosAndroid, goarchAMD64, false},
			{goosLinux, goarchMIPS64, true},
			{goosLinux, goarchMIPS64LE, true},
			{goosNACL, goarch386, false},
			{goosOpenBSD, goarchARM, true},

			{goosAndroid, goarchARM64, false},
			{goosLinux, goarchS390X, true},
			{goosPlan9, goarchARM, false},
		},
		[]Platform{
			{goosNACL, goarchAMD64, false},
			{goosDragonfly, goarch386, false},
		},
	)

	Platforms_1_8 = addDrop(Platforms_1_7,
		[]Platform{
			{goosLinux, goarchMIPS, true},
			{goosLinux, goarchMIPSLE, true},
		},
		nil,
	)

	// no new platforms in 1.9
	Platforms_1_9 = Platforms_1_8

	// no new platforms in 1.10
	Platforms_1_10 = Platforms_1_9

	Platforms_1_11 = addDrop(Platforms_1_10,
		[]Platform{
			{goosJavaScript, goarchWebAssembly, true},
			{goosLinux, goarchRISCV64, false},
		},
		nil,
	)

	Platforms_1_12 = addDrop(Platforms_1_11,
		[]Platform{
			{goosAIX, goarchPowerPC64, false},
			{goosWindows, goarchARM, true},
		},
		[]Platform{
			{goosLinux, goarchRISCV64, false},
		},
	)

	Platforms_1_13 = addDrop(Platforms_1_12,
		[]Platform{
			{goosIllumos, goarchAMD64, false},
			{goosNetBSD, goarchARM64, true},
			{goosOpenBSD, goarchARM64, true},
		},
		nil,
	)

	Platforms_1_14 = addDrop(Platforms_1_13,
		[]Platform{
			{goosFreeBSD, goarchARM64, true},
			{goosLinux, goarchRISCV64, true},
		},
		[]Platform{
			// drop nacl
			{goosNACL, goarch386, false},
			{goosNACL, goarchARM, false},
			{goosNACL, goarchAMD64P32, false},
		},
	)

	Platforms_1_15 = addDrop(Platforms_1_14,
		nil,
		[]Platform{
			{goosDarwin, goarch386, false},
			{goosDarwin, goarchARM, false},
		},
	)

	Platforms_1_16 = addDrop(Platforms_1_15,
		[]Platform{
			{goosAndroid, goarchAMD64, false},
			{goosDarwin, goarchARM64, true},
			{goosIOS, goarchARM64, false},
			{goosIOS, goarchAMD64, false},
			{goosOpenBSD, goarchMIPS64, false},
		},
		nil,
	)

	Platforms_1_17 = addDrop(Platforms_1_16,
		[]Platform{
			{goosWindows, goarchARM64, true},
		},
		nil,
	)

	Platforms_1_18 = addDrop(Platforms_1_17,
		[]Platform{
			{goosIllumos, goarchAMD64, false},
		},
		nil,
	)

	Platforms_1_19 = addDrop(Platforms_1_18,
		[]Platform{
			{goosLinux, goarchLoong64, true},
		},
		nil,
	)

	Platforms_1_20 = addDrop(Platforms_1_19,
		[]Platform{
			{goosFreeBSD, goarchRISCV64, false},
		},
		nil,
	)

	Platforms_1_21 = addDrop(Platforms_1_20,
		[]Platform{
			{goosWASIP1, goarchWebAssembly, false},
		},
		[]Platform{
			{goosOpenBSD, goarchMIPS64, false},
		},
	)

	Platforms_1_22 = Platforms_1_21

	PlatformsLatest = Platforms_1_22
)

// SupportedPlatforms returns the full list of supported platforms for
// the version of Go that is
func SupportedPlatforms(v *version.Version) []Platform {
	if v == nil {
		return PlatformsLatest
	}

	for _, p := range platformConstraints {
		if p.Constraints.Check(v) {
			return p.Platforms
		}
	}

	// Assume latest
	return PlatformsLatest
}

// A PlatformConstraint describes a constraint for a list of platforms.
type PlatformConstraint struct {
	Constraints version.Constraints
	Platforms   []Platform
}
