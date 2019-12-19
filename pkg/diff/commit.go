package diff

import (
	"fmt"
	"gopkg.in/src-d/go-git.v4"
	"gopkg.in/src-d/go-git.v4/plumbing"
	"gopkg.in/src-d/go-git.v4/plumbing/object"
	"log"
)

func FromTagToTag(repo *git.Repository, fromHash, toHash plumbing.Hash) []*object.Commit {
	// ---------------------------------------------
	// Extract the log iterator from the 'from' hash
	// ---------------------------------------------
	iterFrom, errFrom := repo.Log(&git.LogOptions{From: fromHash})
	if errFrom != nil {
		fmt.Println(errFrom.Error())
		log.Fatal("an error occured while retrieving the commit history from 'fromHash'")
	}

	// -------------------------------------------
	// Extract the log iterator from the 'to' hash
	// -------------------------------------------
	iterTo, errTo := repo.Log(&git.LogOptions{From: toHash})
	if errTo != nil {
		fmt.Println(errTo.Error())
		log.Fatal("an error occured while retrieving the commit history from 'toHash'")
	}

	// ----------------------------
	// Initialize the commit slices
	// ----------------------------
	commitFrom := make([]*object.Commit, 0)
	commitTo := make([]*object.Commit, 0)

	_ = iterFrom.ForEach(func(commit *object.Commit) error {
		commitFrom = append(commitFrom, commit)
		return nil
	})

	_ = iterTo.ForEach(func(commit *object.Commit) error {
		commitTo = append(commitTo, commit)
		return nil
	})

	// ------------------------------------------------
	// Perform the actual diff operation on both slices
	// ------------------------------------------------
	diff := make([]*object.Commit, 0)
	for i := 0; i < 2; i++ {
		for _, s1 := range commitTo {
			found := false
			for _, s2 := range commitFrom {
				if s1.Hash == s2.Hash {
					found = true
					break
				}
			}

			if !found {
				diff = append(diff, s1)
			}
		}

		if i == 0 {
			commitTo, commitFrom = commitFrom, commitTo
		}
	}

	return diff
}
