package version

import "strings"

var (
	// these version are set at build time, see scripts/go-build.sh, the version values come from config/vars/Makefile
	Version      = "7.10.0"
	PriorVersion = "7.9.0"
)

// Return the major.minor, as 7.8, instead of 7.8.0
func ShortVersion() string {
	idx := strings.LastIndex(Version, ".")
	return Version[:idx]
}
