package configuration

import (
	"fmt"
	"gopkg.in/src-d/go-git.v4"
	"gopkg.in/src-d/go-git.v4/plumbing"
	"gopkg.in/src-d/go-git.v4/plumbing/object"
	"log"
	"regexp"
	"strings"
)

// Constants for the regex processing of the 'from' and 'to' tags
const (
	regex                 = `([a-zA-Z0-9.\-_$()/]+)==>([a-zA-Z0-9.\-_$()/]+)`
	variableStarter       = "$("
	variableStarterOffset = 2
	variableEnder         = ")"
)

// Constants for possible values of the variables
const (
	latest = "LATEST"
)

// Constants for configuration table in 'ExtractFromToHash'
const (
	fromLatest = "fromLatest"

	toLatest = "toLatest"
)

func ValidateParam(param *string, msg string) {
	if param == nil || *param == "" {
		log.Fatal(msg)
	}
}

//
// <FROM_TAG>==><TO_TAG>
// 'TO' must always be after 'FROM'
//
// Examples:
// 		- 1.0.0-rc.$(LATEST-1)/==>1.0.0-rc.$(LATEST)
//
func ExtractFromToHash(repo *git.Repository, tags []string, diffRegex string) (plumbing.Hash, plumbing.Hash) {
	if diffRegex == "" {
		return plumbing.ZeroHash, plumbing.ZeroHash
	}

	configs := make(map[string]bool, 0)

	r := regexp.MustCompile(regex)

	values := r.FindStringSubmatch(diffRegex)
	if len(values) != 3 {
		// panic
	}

	from := values[1]
	to := values[2]

	//var processedFrom, processedTo bool

	// -----------------------------------
	// Handle variable in 'from' specifier
	// -----------------------------------
	if strings.Contains(from, variableStarter) {
		idxStart := strings.Index(from, variableStarter)
		offset := strings.Index(from[idxStart:], variableEnder)

		variable := from[idxStart+variableStarterOffset : idxStart+offset]

		switch variable {
		case latest:
			configs[fromLatest] = true
		default:
			panic(fmt.Sprintf("Unkown variable %s", variable))
		}
	}

	// ----------------------------------
	// Handle variable in 'to' specififer
	// ----------------------------------
	if strings.Contains(to, variableStarter) {
		idxStart := strings.Index(to, variableStarter)
		offset := strings.Index(to[idxStart:], variableEnder)

		variable := to[idxStart+variableStarterOffset : idxStart+offset]

		switch variable {
		case latest:
			configs[toLatest] = true
		default:
			panic(fmt.Sprintf("Unkown variable %s", variable))
		}
	}

	// --------------------------------------------
	// Validate integrity of from v.s. to specifier
	// --------------------------------------------

	// --------------------------------------------
	// Extract iterator of the resulting commit log
	// --------------------------------------------
	refFrom, errFrom := repo.Tag(from)
	refTo, errTo := repo.Tag(to)

	if errFrom != nil || errTo != nil {
		panic("Unknown tag(s)")
	}

	tagFrom, errFrom := repo.TagObject(refFrom.Hash())
	tagTo, errTo := repo.TagObject(refTo.Hash())

	if errFrom != nil || errTo != nil {
		panic("Unable to find tag(s) object(s)")
	}

	var hashFrom, hashTo plumbing.Hash

	//commits := sort.CommitSliceDiff(repo, refFrom.Hash(), refTo.Hash())
	headRef, _ := repo.Head()

	i, _ := repo.Log(&git.LogOptions{From: headRef.Hash()})

	_ = i.ForEach(func(commit *object.Commit) error {
		if commitFrom, err := tagFrom.Commit(); err == nil {
			if commit.Hash == commitFrom.Hash && hashFrom == plumbing.ZeroHash {
				hashFrom = commit.Hash
			}
		}

		if commitTo, err := tagTo.Commit(); err == nil {
			if commit.Hash == commitTo.Hash && hashTo == plumbing.ZeroHash {
				hashTo = commit.Hash
			}
		}

		return nil
	})

	//fmt.Print(commits)

	return hashFrom, hashTo
}
