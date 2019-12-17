package configuration

import (
	"fmt"
	"gopkg.in/src-d/go-git.v4"
	"gopkg.in/src-d/go-git.v4/plumbing/object"
	"log"
	"regexp"
	"strings"
)

// Constants for the regex processing of the 'from' and 'to' tags
const (
	regex                 = `([a-zA-Z0-9.\-$()]+)==>([a-zA-Z0-9.\-$()]+)`
	variableStarter       = "$("
	variableStarterOffset = 2
	variableEnder         = ")"
)

// Constants for possible values of the variables
const (
	latest = "LATEST"
)

// Constants for configuration table in 'ExtractLogIter'
const (
	fromLatest = "fromLatest"
)

func ValidateParam(param *string, msg string) {
	if param == nil || *param == "" {
		log.Fatal(msg)
	}
}

//
// <FROM_TAG>==><TO_TAG>
//
// Ex.: 1.0.0-rc$(LATEST)/==>1.0.0-rc$(LATEST-1)
//
func ExtractLogIter(repo *git.Repository, diffRegex string) object.CommitIter {
	configs := make(map[string]bool, 0)

	/*
		for _, table := range tables {
				val := Context(table.i).String()
				if val != table.o {
					t.Errorf("String value was incorrect, got: %s, want: %s.", val, table.o)
				}
			}
	*/

	r := regexp.MustCompile(regex)

	values := r.FindStringSubmatch(diffRegex)
	if len(values) != 3 {
		// panic
	}

	from := values[1]
	to := values[2]

	//var processedFrom, processedTo bool

	// Handle variable in 'from' specifier
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

		iter, err := repo.TagObjects()
		if err == nil {
			err = iter.ForEach(func(tag *object.Tag) error {
				fmt.Println(tag.Name)

				return nil
			})
		}
	}

	// Handle variable in 'to' specififer
	if strings.Contains(to, variableStarter) {
		idxStart := strings.Index(to, variableStarter)
		offset := strings.Index(to[idxStart:], variableEnder)

		variable := to[idxStart+variableStarterOffset : idxStart+offset]

		switch variable {
		case latest:
			configs[fromLatest] = true
		default:
			panic(fmt.Sprintf("Unkown variable %s", variable))
		}
	}

	// Validate integrity of from v.s. to specifier

	// Extract iterator of the resulting commit log

	return nil
}
