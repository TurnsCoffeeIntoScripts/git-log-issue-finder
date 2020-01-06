package find

import "gopkg.in/src-d/go-git.v4/plumbing/object"

func CommitMatching(commits []*object.Commit, ticketRegex string) []string {
	var TicketSlice []string
	for _, c := range commits {
		if presentInMessage, ticket := Tickets(c.Message, ticketRegex); presentInMessage {
			TicketSlice = append(TicketSlice, ticket...)
		}
	}

	return TicketSlice
}
