package parser

import (
	"fmt"
	"github.com/tklab-group/forge/util/optional"
	"log/slog"
	"strings"
)

type aptPackageInfo struct {
	name    string
	version optional.Of[string]
}

func parseAptPackageInfo(s string) packageInfo {
	split := strings.Split(s, "=")
	if len(split) == 2 {
		return &aptPackageInfo{
			name:    split[0],
			version: optional.NewWithValue(split[1]),
		}
	} else {
		return &aptPackageInfo{
			name:    split[0],
			version: optional.Of[string]{},
		}
	}
}

func (a *aptPackageInfo) updatePackageInfo(reference packageVersions) {
	aptPkgVersions, ok := reference["apt"]
	if !ok {
		slog.Warn("apt packages information not found")
		return
	}

	version, ok := aptPkgVersions[a.name]
	if !ok {
		slog.Warn(fmt.Sprintf("version of package %s not found", a.name))
	}

	a.version = optional.NewWithValue(version)
	slog.Debug(fmt.Sprintf("update package information: %s", a.toString()))
}

func (a *aptPackageInfo) toString() string {
	if a.version.HasValue() {
		return fmt.Sprintf("%s=%s", a.name, a.version.ValueOrZero())
	} else {
		return a.name
	}
}
