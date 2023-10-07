package parser

import (
	"fmt"
	"github.com/tklab-group/forge/util/optional"
	"strings"
)

type aptPackageInfo struct {
	name    string
	version optional.Of[string]
}

func parseAptPackageInfo(s string) *aptPackageInfo {
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

func (a *aptPackageInfo) toString() string {
	if a.version.HasValue() {
		return fmt.Sprintf("%s=%s", a.name, a.version.ValueOrZero())
	} else {
		return a.name
	}
}
