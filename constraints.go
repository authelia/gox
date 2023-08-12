package main

import "github.com/hashicorp/go-version"

var (
	flagConstraints       map[string]version.Constraints
	platformConstraints   []PlatformConstraint
	experimentConstraints map[string]version.Constraints
)

// Parse all of the constraints as efficiently as possible.
func init() {
	var (
		ok         bool
		constraint version.Constraints
		err        error
	)

	parsed := map[string]version.Constraints{}

	flagConstraints = map[string]version.Constraints{}

	flagConstraintsInputs := []struct {
		name       string
		constraint string
	}{
		{"C", ">= 1.20"},
		{flagNameASAN, ">= 1.18"},
		{flagNameCover, ">= 1.20"},
		{flagNameCoverPKG, ">= 1.20"},
		{flagNameBuildVCS, ">= 1.18"},
		{flagNameModCacheRW, ">= 1.14"},
		{flagNameMod, ">= 1.11"},
		{flagNameModFile, ">= 1.14"},
		{flagNameOverlay, ">= 1.16"},
		{flagNameProfileGuidedOptimization, ">= 1.20"},
		{flagNameTrimPath, ">= 1.13"},
	}

	for _, input := range flagConstraintsInputs {
		if constraint, ok = parsed[input.constraint]; ok {
			flagConstraints[input.name] = constraint

			continue
		}

		if constraint, err = version.NewConstraint(input.constraint); err != nil {
			panic(err)
		}

		flagConstraints[input.name] = constraint
		parsed[input.constraint] = constraint
	}

	platformConstraintsInputs := []struct {
		constraint string
		platforms  []Platform
	}{
		{"<= 1.0", Platforms_1_0},
		{">= 1.1, < 1.2", Platforms_1_1},
		{">= 1.2, < 1.3", Platforms_1_2},
		{">= 1.3, < 1.4", Platforms_1_3},
		{">= 1.4, < 1.5", Platforms_1_4},
		{">= 1.5, < 1.6", Platforms_1_5},
		{">= 1.6, < 1.7", Platforms_1_6},
		{">= 1.7, < 1.8", Platforms_1_7},
		{">= 1.8, < 1.9", Platforms_1_8},
		{">= 1.9, < 1.10", Platforms_1_9},
		{">= 1.10, < 1.11", Platforms_1_10},
		{">= 1.11, < 1.12", Platforms_1_11},
		{">= 1.12, < 1.13", Platforms_1_12},
		{">= 1.13, < 1.14", Platforms_1_13},
		{">= 1.14, < 1.15", Platforms_1_14},
		{">= 1.15, < 1.16", Platforms_1_15},
		{">= 1.16, < 1.17", Platforms_1_16},
		{">= 1.17, < 1.18", Platforms_1_17},
		{">= 1.18, < 1.19", Platforms_1_18},
		{">= 1.19, < 1.20", Platforms_1_19},
		{">= 1.20, < 1.21", Platforms_1_20},
		{">= 1.21, < 1.22", Platforms_1_21},
	}

	platformConstraints = make([]PlatformConstraint, len(platformConstraintsInputs))

	for i, input := range platformConstraintsInputs {
		if constraint, ok = parsed[input.constraint]; ok {
			platformConstraints[i] = PlatformConstraint{Constraints: constraint, Platforms: input.platforms}

			continue
		}

		if constraint, err = version.NewConstraint(input.constraint); err != nil {
			panic(err)
		}

		platformConstraints[i] = PlatformConstraint{Constraints: constraint, Platforms: input.platforms}
		parsed[input.constraint] = constraint
	}
}
