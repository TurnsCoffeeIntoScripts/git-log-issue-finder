package find

import "gopkg.in/src-d/go-git.v4/plumbing/object"

// CommitMatching returns a slice of string containing every ticket (issue) matching the specified regex
func CommitMatching(commits []*object.Commit, ticketRegex string) []string {
	var TicketSlice []string
	for _, c := range commits {
		if presentInMessage, ticket := Tickets(c.Message, ticketRegex); presentInMessage {
			TicketSlice = append(TicketSlice, ticket...)
		}
	}

	return TicketSlice
}
