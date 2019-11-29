package main

import (
	"flag"
	"fmt"
	"github.com/TurnsCoffeeIntoScripts/git-log-issue-finder/pkg/find"
	"github.com/TurnsCoffeeIntoScripts/git-log-issue-finder/pkg/sort"
	"gopkg.in/src-d/go-git.v4"
	"gopkg.in/src-d/go-git.v4/plumbing/object"
	"log"
)

// TicketSlice slice containing the extracted tickets from the content of the log that was provided
var TicketSlice []string

func main() {
	// Parameters
	ticketRegex := flag.String("tickets", "", "Comma-separated list of jira project keys")
	repoDir := flag.String("directory", "", "The directory of the git repo")

	// Flags
	fullHistory := flag.Bool("full-history", false, "Search the entire git log")
	sinceLatestTag := flag.Bool("since-latest-tag", false, "Search only from HEAD to the most recent tag")

	flag.Parse()

	validateParam(ticketRegex, "Missing parameter '--tickets'")
	validateParam(repoDir, "Missing parameter '--directory'")
	if *fullHistory == *sinceLatestTag { // Only one should be set
		*fullHistory = true
		*sinceLatestTag = false
	}

	repo, err := git.PlainOpen(*repoDir)
	if err != nil {
		fmt.Print(err.Error())
		log.Fatal("an error occured while instantiating a new repository object")
	}

	ref, err := repo.Head()
	if err != nil {
		fmt.Print(err)
		log.Fatal("an error occured while retrieving the HEAD reference")
	}

	if *sinceLatestTag {
		commits := make([]*object.Commit, 0)
		iter, err := repo.TagObjects()
		hashLatest := find.LatestTag(iter, ref.Hash(), err)
		commits = sort.CommitSliceDiff(repo, ref.Hash(), hashLatest)
		for _, c := range commits {
			if presentInMessage, ticket := find.Tickets(c.Message, *ticketRegex); presentInMessage {
				TicketSlice = append(TicketSlice, ticket...)
			}
		}
	} else if *fullHistory {
		iter, err := repo.Log(&git.LogOptions{From: ref.Hash()})
		if err != nil {
			fmt.Print(err.Error())
			log.Fatal("an error occured while retrieving the commit history")
		}

		err = iter.ForEach(func(c *object.Commit) error {
			if presentInMessage, ticket := find.Tickets(c.Message, *ticketRegex); presentInMessage {
				TicketSlice = append(TicketSlice, ticket...)
			}

			return nil
		})
	}

	TicketSlice = ensureUniqueValues(TicketSlice)

	fmt.Println(TicketSlice)
}

func ensureUniqueValues(s []string) []string {
	keys := make(map[string]bool)
	var list []string
	for _, entry := range s {
		if _, value := keys[entry]; !value {
			keys[entry] = true
			list = append(list, entry)
		}
	}

	return list
}

func validateParam(param *string, msg string) {
	if param == nil || *param == "" {
		log.Fatal(msg)
	}
}
