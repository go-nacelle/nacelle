package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"strings"

	"github.com/coreos/go-semver/semver"
	"golang.org/x/mod/modfile"
)

const importPathPrefix = "github.com/go-nacelle"

func main() {
	if err := mainErr(); err != nil {
		fmt.Fprint(os.Stderr, fmt.Sprintf("error: %s\n", err.Error()))
		os.Exit(1)
	}
}

func mainErr() error {
	// TODO - get from user/git
	currentVersion := semver.New("1.0.0")

	depsBefore, err := parseGoMod()
	if err != nil {
		return err
	}

	var paths []string
	for path := range depsBefore {
		paths = append(paths, path)
	}

	_ = exec.Command("go", append([]string{"get", "-u"}, paths...)...).Run()

	depsAfter, err := parseGoMod()
	if err != nil {
		return err
	}

	return compareVersions(currentVersion, depsBefore, depsAfter)
}

func parseGoMod() (map[string]string, error) {
	contents, err := ioutil.ReadFile("go.mod")
	if err != nil {
		return nil, err
	}

	f, err := modfile.Parse("go.mod", contents, nil)
	if err != nil {
		return nil, err
	}

	modules := map[string]string{}
	for _, require := range f.Require {
		if strings.HasPrefix(require.Mod.Path, importPathPrefix) {
			version := require.Mod.Version
			if strings.HasPrefix(version, "v") {
				version = version[1:]
			}
			modules[require.Mod.Path] = version
		}
	}

	return modules, nil
}

const (
	Major = 0
	Minor = 1
	Patch = 2
)

func compareVersions(currentVersion *semver.Version, depsBefore, depsAfter map[string]string) error {
	changes := map[int]int{}
	for path, beforeVersion := range depsBefore {
		b, err := semver.NewVersion(beforeVersion)
		if err != nil {
			return err
		}

		a, err := semver.NewVersion(depsAfter[path])
		if err != nil {
			return err
		}

		if a.Major != b.Major {
			changes[Major]++
		} else if a.Minor != b.Minor {
			changes[Minor]++
		} else if a.Patch != b.Patch {
			changes[Patch]++
		} else {
			continue
		}

		fmt.Printf("Updating %s %s -> %s\n", path, b, a)
	}

	if changes[Major] == 0 && changes[Minor] == 0 && changes[Patch] == 0 {
		return nil
	}

	newVersion := semver.New(currentVersion.String())

	if changes[Major] > 0 {
		newVersion.BumpMajor()
	} else if changes[Minor] > 0 {
		newVersion.BumpMinor()
	} else if changes[Patch] > 0 {
		newVersion.BumpPatch()
	}

	fmt.Printf("Bumping version %s -> %s\n", currentVersion, newVersion)
	return nil
}
