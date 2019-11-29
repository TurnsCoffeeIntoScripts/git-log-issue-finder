package main

import (
	"flag"
	"fmt"
	"gopkg.in/src-d/go-git.v4"
	"gopkg.in/src-d/go-git.v4/plumbing/object"
	"log"
	"regexp"
	"strings"
)

// TicketSlice slice containing the extracted tickets from the content of the log that was provided
var TicketSlice []string

func main() {
	ticketRegex := flag.String("tickets", "", "Comma-separated list of jira project keys")
	repoDir := flag.String("directory", "", "The directory of the git repo")

	flag.Parse()

	validateParam(ticketRegex, "Missing parameter '--tickets'")
	validateParam(repoDir, "Missing parameter '--directory'")

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

	iter, err := repo.Log(&git.LogOptions{From: ref.Hash()})
	if err != nil {
		fmt.Print(err)
		log.Fatal("an error occured while retrieving the commit history")
	}

	err = iter.ForEach(func(c *object.Commit) error {
		if presentInMessage, ticket := findTickets(c.Message, *ticketRegex); presentInMessage {
			TicketSlice = append(TicketSlice, ticket...)
		}

		return nil
	})

	TicketSlice = ensureUniqueValues(TicketSlice)

	fmt.Println(TicketSlice)
}

func findTickets(text, ticketRegex string) (bool, []string) {
	regex := "((?:"
	regex += strings.ReplaceAll(ticketRegex, ",", "|")
	regex += ")-[0-9]+)"

	r, _ := regexp.Compile(regex)

	out := r.FindAllString(text, -1)

	if len(out) == 0 {
		return false, []string{}
	}

	return true, out
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
