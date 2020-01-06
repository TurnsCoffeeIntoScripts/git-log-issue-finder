package find

import (
	"regexp"
	"strings"
)

// Tickets returns a slice of every ticket (issue) found in the text based on the provided regex
// Returns false if nothing was found
func Tickets(text, ticketRegex string) (bool, []string) {
	regex := "((?:"
	if ticketRegex == "*" {
		regex += "[a-zA-Z0-9]+"
	} else {
		regex += strings.ReplaceAll(ticketRegex, ",", "|")
	}
	regex += ")-[0-9]+)"

	r, _ := regexp.Compile(regex)

	out := r.FindAllString(text, -1)

	if len(out) == 0 {
		return false, []string{}
	}

	return true, out
}
