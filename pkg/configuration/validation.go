package configuration

import (
	"fmt"
	"gopkg.in/src-d/go-git.v4"
	"gopkg.in/src-d/go-git.v4/plumbing"
	"gopkg.in/src-d/go-git.v4/plumbing/object"
	"log"
	"regexp"
	"strconv"
	"strings"
)

// Constants for the regex processing of the 'from' and 'to' tags
const (
	regex                 = `([a-zA-Z0-9.\-_@()/]+)==>([a-zA-Z0-9.\-_@()/]+)`
	variableStarter       = "@("
	variableStarterOffset = 2
	variableEnder         = ")"
)

// Constants for possible values of the variables
const (
	first  = "FIRST"
	latest = "LATEST"
)

// Constant for math operations
const (
	minusSign = "-"
	plusSign  = "+"
)

// Constants for configuration table in 'ExtractFromToHash'
const (
	fromLatest      = "fromLatest"
	fromLatestMinus = "fromLatestMinus"
	fromFirst       = "fromFirst"
	fromFirstPlus   = "fromFirstPlus"

	toLatest      = "toLatest"
	toLatestMinus = "toLatestMinus"
	toFirst       = "toFirst"
	toFirstPlus   = "toFirstPlus"
)

// ValidateParam simply validate that the input flag (which a *string) isn't nil or
// is not set the the string zero value. In either case log.Fatal is used
func ValidateParam(param *string, msg string) {
	if param == nil || *param == "" {
		log.Fatal(msg)
	}
}

// ExtractFromToHash will extract both hash of the 'from' commit to the 'to' commit
// <FROM_TAG>==><TO_TAG>
// 'TO' must always be after 'FROM'
// The diffRegex can contain variable or not.
//
// ## Without variable
// Example:
//		1.0.0-rc.10==>1.0.0-rc.11
// This will return the hash of the commits associated with the tags rc.10 and rc.11
//
// ## With variable
// Examples:
// 		- 1.0.0-rc.$(LATEST-1)/==>1.0.0-rc.$(LATEST)
// If the latest 'rc' tag for version 1.0.0 is 15 then this will return the hash of
// the commits associated with the tags rc.14 and rc.15
//
func ExtractFromToHash(repo *git.Repository, tags []string, diffRegex string) (plumbing.Hash, plumbing.Hash) {
	if diffRegex == "" {
		return plumbing.ZeroHash, plumbing.ZeroHash
	}

	configs := make(map[string]bool, 0)

	r := regexp.MustCompile(regex)

	values := r.FindStringSubmatch(diffRegex)
	if len(values) != 3 {
		// panic TODO
	}

	from := values[1]
	to := values[2]
	var variableFrom, variableTo string

	// -----------------------------------
	// Handle variable in 'from' specifier
	// -----------------------------------
	if strings.Contains(from, variableStarter) {
		idxStart := strings.Index(from, variableStarter)
		offset := strings.Index(from[idxStart:], variableEnder)

		// If the variable is NOT at the end of the specifier then we panic
		// because it's not a supported format and we don't know how to
		// properly handle that
		// The '+1' is to take into account the closing parenthesis after the variable offset
		if idxStart+offset+1 != len(from) {
			panic(fmt.Sprintf("Variable was not specified at the end of the 'from' specifier: %s", from))
		}

		variable := from[idxStart+variableStarterOffset : idxStart+offset]
		variableFrom = strings.TrimSpace(variable)

		if variableFrom == latest {
			configs[fromLatest] = true
		} else if variableFrom == first {
			configs[fromFirst] = true
		} else if strings.HasPrefix(variableFrom, latest) && strings.Contains(variableFrom, minusSign) {
			configs[fromLatest] = false
			configs[fromLatestMinus] = true
		} else {
			panic(fmt.Sprintf("Unkown variable %s", variableFrom))
		}
	}

	// ----------------------------------
	// Handle variable in 'to' specififer
	// ----------------------------------
	if strings.Contains(to, variableStarter) {
		idxStart := strings.Index(to, variableStarter)
		offset := strings.Index(to[idxStart:], variableEnder)

		// If the variable is NOT at the end of the specifier then we panic
		// because it's not a supported format and we don't know how to
		// properly handle that
		// The '+1' is to take into account the closing parenthesis after the variable offset
		if idxStart+offset+1 != len(to) {
			panic(fmt.Sprintf("Variable was not specified at the end of the 'from' specifier: %s", to))
		}

		variable := to[idxStart+variableStarterOffset : idxStart+offset]
		variableTo = strings.TrimSpace(variable)

		if variableTo == latest {
			configs[toLatest] = true
		} else if strings.HasPrefix(variableTo, latest) && strings.Contains(variableTo, minusSign) {
			configs[toLatest] = false
			configs[toLatestMinus] = true
		} else {
			panic(fmt.Sprintf("Unkown variable %s", variableTo))
		}
	}

	// --------------------------------------------
	// Validate integrity of from v.s. to specifier
	// --------------------------------------------
	// TODO add validation for first and firstPlus
	if configs[fromLatest] && configs[toLatestMinus] {
		// This is an invalid range. What's been specified is we want the range starting from the latest value
		// to a latest minus an 'offset' value.
		log.Panic("Invalid diff range")
	}

	if configs[fromLatest] && configs[toLatest] {
		// This is technically not an invalid range, but will return nothing since we'd be returning the diff between
		// two commits (tags) that are actually the same
		log.Panic("Invalid diff range")
	}

	if configs[fromLatestMinus] && configs[toLatestMinus] {
		// This *MIGHT* be an invalid range. To be sure, we need to validate both offset values that were specified
		offsetFrom := extractOffset(variableFrom)
		offsetTo := extractOffset(variableTo)

		if offsetFrom <= offsetTo {
			// The 'from' specified is after the 'to' specified, therefore creating an invalid range
			log.Panic("Invalid diff range")
		}
	}

	// -----------------------------------------------
	// Apply variables to 'from' and 'to' if necessary
	// -----------------------------------------------
	// TODO apply first and firstPlus
	to = stripVariableMetaChar(to, variableTo)
	if configs[toLatest] {
		to = getLatest(repo, tags, to)
	} else if configs[toLatestMinus] {
		offset := extractOffset(variableTo)
		to = getTagWithOffset(repo, tags, to, offset, false)
	}

	from = stripVariableMetaChar(from, variableFrom)
	if configs[fromLatest] {
		from = getLatest(repo, tags, from)
	} else if configs[fromLatestMinus] {
		offset := extractOffset(variableFrom)
		from = getTagWithOffset(repo, tags, from, offset, false)
	} else if configs[fromFirst] {
		from = getFirst(repo, tags, from)
	}

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

	return hashFrom, hashTo
}

func stripVariableMetaChar(s, v string) string {
	s = strings.ReplaceAll(s, v, "")
	s = strings.ReplaceAll(s, variableStarter, "")
	s = strings.ReplaceAll(s, variableEnder, "")

	return s
}

func extractOffset(v string) int {
	if strings.Contains(v, minusSign) {
		offset, err := strconv.Atoi(strings.Split(v, minusSign)[1])
		if err != nil {
			return 0
		}

		return offset
	}

	return 0
}

func getFirst(repo *git.Repository, tagsName []string, v string) string {
	return getTagWithOffset(repo, tagsName, v, 0, true)
}

func getLatest(repo *git.Repository, tagsName []string, v string) string {
	return getTagWithOffset(repo, tagsName, v, 0, false)
}

func getTagWithOffset(repo *git.Repository, tagsName []string, v string, offset int, reversed bool) string {
	var commitTimeline = make([]*object.Commit, 0)
	var tagsTimeline = make([]string, 0)

	for _, t := range tagsName {
		if !strings.Contains(t, v) {
			continue
		}

		ref, err := repo.Tag(t)
		if err != nil {
			panic("Unknown tag while fetching latest")
		}

		tag, err := repo.TagObject(ref.Hash())
		if err != nil {
			panic("Unable to find tag object while fetching latest")
		}

		c, err := tag.Commit()
		if err != nil {
			panic("Unable to find commit object while fetching latest")
		}

		if len(commitTimeline) == 0 {
			commitTimeline = append(commitTimeline, c)
			tagsTimeline = append(tagsTimeline, t)
		} else {
			fixedTimelineLength := len(commitTimeline)
			for index := 0; index < fixedTimelineLength; index++ {
				if c.Committer.When.Before(commitTimeline[index].Committer.When) {
					if index == fixedTimelineLength-1 {
						commitTimeline = append(commitTimeline, c)
						tagsTimeline = append(tagsTimeline, t)

						// Break after the append so we don't double add a value in the timeline
						break
					} else {
						continue
					}
				}

				// If we're here, it means that the new commit if after the one at the specified 'index'
				// Update the commit timeline
				commitTimeline = append(commitTimeline, nil)
				copy(commitTimeline[index+1:], commitTimeline[index:]) // Equivalent to a "shift right by one at index"
				commitTimeline[index] = c                              // Then insert new value at index to overwrite the previous value that's now at index+1

				// Update the tags timeline
				tagsTimeline = append(tagsTimeline, "")
				copy(tagsTimeline[index+1:], tagsTimeline[index:]) // Equivalent to a "shift right by one at index"
				tagsTimeline[index] = t                            // Then insert new value at index to overwrite the previous value that's now at index+1

				// Ensure that no more than one insert will be done per loop
				break

			}
		}
	}

	// If the offset specified is too big, then by default we set the offset as 0
	if 0+offset >= len(tagsTimeline) {
		offset = 0
	}

	// If the reversed flag was set to true, that means we are looking at the beginning of the timeline first (e.g. The oldest ones)
	// Used with 'FIRST'
	if reversed {
		return tagsTimeline[(len(tagsTimeline)-1)-offset]
	}

	// Otherwise we look at the end first to get the latest value with the specified offset
	return tagsTimeline[0+offset]
}
