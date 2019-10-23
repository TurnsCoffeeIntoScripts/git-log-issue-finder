package main

import (
	"bufio"
	"flag"
	"fmt"
	"log"
	"regexp"
	"strings"
)

// TicketSlice slice containing the extracted tickets from the content of the log that was provided
var TicketSlice []string

func main() {
	ticketRegex := flag.String("tickets", "", "")
	content := flag.String("content", "", "")

	flag.Parse()

	if ticketRegex == nil || *ticketRegex == "" {
		log.Fatal("Missing parameter '--tickets'")
	}

	if content == nil || *content == "" {
		log.Fatal("Missing parameter '--content'")
	}

	scanner := bufio.NewScanner(strings.NewReader(*content))
	for scanner.Scan() {
		if present, ticket := findTickets(scanner.Text(), *ticketRegex); present {
			TicketSlice = append(TicketSlice, ticket...)
		}
	}
	if err := scanner.Err(); err != nil {
		log.Fatalf("reading standard input: %v", err)
	}

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
