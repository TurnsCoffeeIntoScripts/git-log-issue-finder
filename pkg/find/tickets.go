package find

import (
	"regexp"
	"strings"
)

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
