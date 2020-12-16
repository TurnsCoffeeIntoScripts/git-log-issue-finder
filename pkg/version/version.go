// Package version specifies the version informations
package version

import "fmt"

// Definition of the semver elements (NAME: MAJOR.MINOR.PATCH)
const (
	major = 2
	minor = 0
	patch = 1
	name  = "G.L.I.F."
)

// Get returns the formatted string containing the version informations
func Get() string {
	return fmt.Sprintf("%s: %d.%d.%d", name, major, minor, patch)
}
