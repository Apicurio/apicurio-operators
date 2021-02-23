package version

import "strings"

var (
	Version      = "7.9.0"
	PriorVersion = "0.2.0"
)

// Return the major.minor, as 7.8, instead of 7.8.0
func ShortVersion() string {
	idx := strings.LastIndex(Version, ".")
	return Version[:idx]
}
