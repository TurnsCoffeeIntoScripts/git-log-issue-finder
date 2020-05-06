/*Package configuration is where all program's input flags are specified. It is also where their validation happens.

Usage

Instantiate a new GlifParameters object at the beginning and use it like so:

GO CODE   --------------------------------------------------
glifParam := GlifParameters{}
if ok := glifParam.Parse(); !ok {
	panic("failed to properly parse input flags")
}
// Program can run using valide input flags
------------------------------------------------------------
The Parse() method will call the flag.String(), flag.Bool(), etc... to set all necessary values
Finally a call to the methode validate() is made and it's where application-related validation happen.
*/
package configuration

import (
	"flag"
	"github.com/TurnsCoffeeIntoScripts/git-log-issue-finder/pkg/helpers"
	"reflect"
)

// Definition of constants that are use for the 'flag' setup
const (
	// Parameters
	script  = "script"
	tickets = "tickets"

	// Flags
	repl       = "repl"
	forceFetch = "force-fetch"

	// Pre configured scripts
	diffLatestSemverWithLatestBuilds = "semver-latest-builds"
	diffLatestSemverWithLatestRCs    = "semver-latest-rcs"
	diffLatestSemver                 = "semver-latest"

	// Default values and descriptions for both paramaters and flags
	scriptDefault         = ""
	scriptDescription     = "The glif script file to execute"
	ticketsDefault        = "*"
	ticketsDescription    = "The Jira tickets regex used to search the repo's log"
	forceFetchDefault     = false
	forceFetchDescription = "Force a 'git fetch' operation on the specified repository"
)

// GlifParameters contains the various flags that were given via the program's input paramters
// It also contains an instance of GlifFlags and GlifPreConfiguredScripts
//   - see: configuration.GlifFlags
//	 - see: configuration.GlifPreConfiguredScripts
type GlifParameters struct {
	Script  *string
	Tickets *string

	Flags   GlifFlags
	Scripts GlifPreConfiguredScripts
}

// GlifFlags contains the various boolean flags (actual command line flags and not parameters) used by glif.
type GlifFlags struct {
	REPL       *bool
	ForceFetch *bool
}

// GlifPreConfiguredScripts contains only boolean flags that specify if a "preconfigured" script should be used.
//   - see script package
type GlifPreConfiguredScripts struct {
	UseDiffLatestSemverWithLatestBuilds *bool
	UseDiffLatestSemverWithLatestRCs    *bool
	UseDiffLatestSemver                 *bool
}

// Parse encapsulate the function calls to the Go flag package. It also internally runs an application-related validation.
func (params *GlifParameters) Parse(forceRepl bool) bool {
	params.Script = flag.String(script, scriptDefault, scriptDescription)
	params.Tickets = flag.String(tickets, ticketsDefault, ticketsDescription)

	params.Flags.REPL = flag.Bool(repl, forceRepl, "")
	params.Flags.ForceFetch = flag.Bool(forceFetch, forceFetchDefault, forceFetchDescription)

	params.Scripts.UseDiffLatestSemverWithLatestBuilds = flag.Bool(diffLatestSemverWithLatestBuilds, false, "script.DiffLatestSemverWithLatestBuilds")
	params.Scripts.UseDiffLatestSemverWithLatestRCs = flag.Bool(diffLatestSemverWithLatestRCs, false, "script.DiffLatestSemverWithLatestRCs")
	params.Scripts.UseDiffLatestSemver = flag.Bool(diffLatestSemver, false, "script.DiffLatestSemver")

	flag.Parse()

	var ok bool
	if ok = params.validate(); !ok {
		return false
	}

	return ok
}

func (params *GlifParameters) validate() bool {
	// If REPL was specified we skip the rest of the validation because glif will be set in interactive mode
	if helpers.IsBoolPtrTrue(params.Flags.REPL) {
		return true
	}

	if helpers.IsStringPtrNilOrEmtpy(params.Script) {
		// No script was specified, before failing the validation we need to check if any of the preconfigured
		// script were declared

		count := 0
		v := reflect.ValueOf(params.Scripts)

		for i := 0; i < v.NumField(); i++ {
			f := v.Field(i).Interface()
			bPtr := f.(*bool)

			if helpers.IsBoolPtrTrue(bPtr) {
				count++
			}
		}

		return count == 1
	}

	// Process 'from'
	//tags.ProcessDiffTag(*param.FromTags, param.Flags)
	//tags.ProcessDiffTag(to)

	return true
}
