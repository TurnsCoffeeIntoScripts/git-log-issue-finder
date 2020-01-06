package main

import (
	"flag"
	"fmt"
	"github.com/TurnsCoffeeIntoScripts/git-log-issue-finder/pkg/configuration"
	"github.com/TurnsCoffeeIntoScripts/git-log-issue-finder/pkg/diff"
	"github.com/TurnsCoffeeIntoScripts/git-log-issue-finder/pkg/find"
	"gopkg.in/src-d/go-git.v4"
	"gopkg.in/src-d/go-git.v4/plumbing"
	"gopkg.in/src-d/go-git.v4/plumbing/object"
	"log"
	"sort"
)

// TicketSlice slice containing the extracted tickets from the content of the log that was provided
var TicketSlice []string

// Full text for the usage of '--diff-tags'
var diffTagDesc = "This parameter takes a specially formatted string: <FROM_TAG>==><TO_TAG>\n" +
	"\tThe '==>' is litteral\n" +
	"\tThe '<FROM_TAG>' and '<TO_TAG>' can have 2 formats:\n" +
	"\t\t1- A litteral tag name\n" +
	"\t\t2- A string with parameters (see below)\n" +
	"Parameters are values declared between the '@(' and ')' litterals. Possible parameters:\n" +
	"\tLATEST: finds the latest value matching the string\n" +
	"\tLATEST-N: finds the Nth commit behind the latest value matching the string\n" +
	"Examples:\n" +
	"--diff-tags=\"v1.0.0-rc.@(LATEST)==>v1.0.0-rc.@(LATEST-1)\"\n" +
	"\tThis will return the issues found between the latest RC of version 1.0.0 and the RC before that one"

func main() {
	// Parameters
	ticketRegex := flag.String("tickets", "", "Comma-separated list of jira project keys")
	repoDir := flag.String("directory", "./", "The directory of the git repo")
	diffTags := flag.String("diff-tags", "", diffTagDesc)

	// Flags
	fullHistory := flag.Bool("full-history", false, "Search the entire git log")
	sinceLatestTag := flag.Bool("since-latest-tag", false, "Search only from HEAD to the most recent tag")
	forceFetch := flag.Bool("force-fetch", false, "Force a 'git fetch' operation on the specified repository")

	// Flags
	fullHistory := flag.Bool("full-history", false, "Search the entire git log")
	sinceLatestTag := flag.Bool("since-latest-tag", false, "Search only from HEAD to the most recent tag")

	flag.Parse()

	// Validating mandatory params
	configuration.ValidateParam(ticketRegex, "Missing parameter '--tickets'")
	configuration.ValidateParam(repoDir, "Missing parameter '--directory'")

	// Only one should be set
	if *fullHistory == *sinceLatestTag {
		*fullHistory = true
		*sinceLatestTag = false
	}

	repo, err := git.PlainOpen(*repoDir)
	if err != nil {
		fmt.Print(err.Error())
		log.Fatal("an error occured while instantiating a new repository object")
	}

	if forceFetch != nil && *forceFetch {
		if fetchErr := repo.Fetch(&git.FetchOptions{}); fetchErr != nil {
			fmt.Println(fetchErr.Error())
			log.Print("an error occured while trying to fetch repo. Continuing without fetching...")
		}
	}

	ref, err := repo.Head()
	if err != nil {
		fmt.Print(err)
		log.Fatal("an error occured while retrieving the HEAD reference")
	}

	hashFrom, hashTo := configuration.ExtractFromToHash(repo, extractTagList(repo), *diffTags)

	if hashFrom != plumbing.ZeroHash && hashTo != plumbing.ZeroHash {
		commits := diff.FromTagToTag(repo, hashFrom, hashTo)
		TicketSlice = find.CommitMatching(commits, *ticketRegex)
	} else if *sinceLatestTag {
		commits := make([]*object.Commit, 0)
		iter, err := repo.TagObjects()
		hashLatest := find.LatestTag(iter, ref.Hash(), err)
		commits = diff.FromTagToTag(repo, hashLatest, ref.Hash())
		TicketSlice = find.CommitMatching(commits, *ticketRegex)
	} else if *fullHistory {
		iter, err := repo.Log(&git.LogOptions{From: ref.Hash()})
		if err != nil {
			fmt.Print(err.Error())
			log.Fatal("an error occured while retrieving the commit history")
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

func extractTagList(repo *git.Repository) []string {
	tags := make([]string, 0)

	iter, err := repo.TagObjects()
	if err == nil {
		err = iter.ForEach(func(tag *object.Tag) error {
			tags = append(tags, tag.Name)

			return nil
		})
	}

	sort.Strings(tags)
	return tags
}
