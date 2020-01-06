package sort

import (
	"fmt"
	"gopkg.in/src-d/go-git.v4"
	"gopkg.in/src-d/go-git.v4/plumbing"
	"gopkg.in/src-d/go-git.v4/plumbing/object"
	"log"
)

func CommitSliceDiff(repo *git.Repository, hashHead, hashLatest plumbing.Hash) []*object.Commit {
	iterFromHead, err := repo.Log(&git.LogOptions{From: hashHead})
	if err != nil {
		fmt.Print(err.Error())
		log.Fatal("an error occured while retrieving the commit history from head")
	}

	iterFromLatest, err := repo.Log(&git.LogOptions{From: hashLatest})
	if err != nil {
		fmt.Print(err.Error())
		log.Fatal("an error occured while retrieving the commit history from latest")
	}

	commitFromHead := make([]*object.Commit, 0)
	commitFromLatest := make([]*object.Commit, 0)
	err = iterFromHead.ForEach(func(commit *object.Commit) error {
		commitFromHead = append(commitFromHead, commit)
		return nil
	})
	if err != nil {
		fmt.Print(err.Error())
		log.Fatal("an error occured while reading commits from head")
	}

	err = iterFromLatest.ForEach(func(commit *object.Commit) error {
		commitFromLatest = append(commitFromLatest, commit)
		return nil
	})
	if err != nil {
		fmt.Print(err.Error())
		log.Fatal("an error occured while reading commits from latest")
	}

	diff := make([]*object.Commit, 0)
	for i := 0; i < 2; i++ {
		for _, s1 := range commitFromHead {
			found := false
			for _, s2 := range commitFromLatest {
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
			commitFromHead, commitFromLatest = commitFromLatest, commitFromHead
		}
	}

	return diff
}
