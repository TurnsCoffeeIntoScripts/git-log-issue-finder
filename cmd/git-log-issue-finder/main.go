package main

import (
	"flag"
	"fmt"
	"gopkg.in/src-d/go-git.v4"
	"gopkg.in/src-d/go-git.v4/plumbing/object"
	"gopkg.in/src-d/go-git.v4/storage/memory"
	"log"
	"regexp"
	"strings"
)

// TicketSlice slice containing the extracted tickets from the content of the log that was provided
var TicketSlice []string

func main() {
	ticketRegex := flag.String("tickets", "", "Comma-separated list of jira project keys")
	repoURL := flag.String("repo", "", "The directory of the git repo")
	username := flag.String("username", "", "Username for the git repo")
	password := flag.String("password", "", "Password for the git repo")

	flag.Parse()

	validateParam(ticketRegex, "Missing parameter '--tickets'")
	validateParam(repoURL, "Missing parameter '--repo'")
	validateParam(username, "Missing parameter '--username'")
	validateParam(password, "Missing parameter '--password'")

	o := &git.CloneOptions{}
	o.URL = addCredentialsToURL(*repoURL, *username, *password)

	repo, err := git.Clone(memory.NewStorage(), nil, o)
	if err != nil {
		fmt.Print(err)
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

func addCredentialsToURL(url, username, password string) string {
	partialURL := strings.Split(url, "://")

	if len(partialURL) >= 2 {
		return partialURL[0] + "://" + username + ":" + password + "@" + partialURL[1]
	} else if len(partialURL) == 1 {
		return "https://" + username + ":" + password + "@" + partialURL[1]
	}

	return url
}
