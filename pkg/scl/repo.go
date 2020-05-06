// Package scl is an abstraction layer on top of the go-git library.
//
// The 'GlifRepo' object is the struct that encloses the *git.Repository.
// It also contains maps and slices of tags object used to perform the diffs.
//
// Some of these function are called from the builtins of the interpreter, tying the .gs scripts to this struct.
package scl

import (
	"fmt"
	"github.com/TurnsCoffeeIntoScripts/git-log-issue-finder/pkg/configuration"
	"github.com/TurnsCoffeeIntoScripts/git-log-issue-finder/pkg/helpers"
	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing"
	"github.com/go-git/go-git/v5/plumbing/object"
	"log"
	"regexp"
	"sort"
)

// GlifRepo is an abstraction on top of the *git.Repository to allow extended functionnality.
type GlifRepo struct {
	Params configuration.GlifParameters

	HeadRef *plumbing.Reference
	GitRepo *git.Repository

	matchingTags         map[string]*object.Tag // Map that take the tag name as index to match the tag object (*object.Tag)
	tagsLatestToEarliest []*object.Tag          // This slice is ordered from latest to earliest. Meaning that the tag at index 0 is the very latest.
}

// Open does a 'PlainOpen' on the *git.Repository and will create the map and slice of *object.Tag.
func (glifRepo *GlifRepo) Open(repoLoc string) {
	repo, err := git.PlainOpen(repoLoc)
	if err != nil {
		fmt.Println(err.Error())
		log.Fatal("an error occured while instantiating a new repository object")
	}

	glifRepo.GitRepo = repo
	glifRepo.matchingTags = make(map[string]*object.Tag)
	glifRepo.tagsLatestToEarliest = make([]*object.Tag, 0)
}

// Fetch does a 'Fetch' on the *git.Repository.
func (glifRepo *GlifRepo) Fetch() {
	if fetchErr := glifRepo.GitRepo.Fetch(&git.FetchOptions{}); fetchErr != nil {
		fmt.Println(fetchErr.Error())
		log.Fatal("an error occured while trying to fetch repo. Continuing without fetching...")
	}
}

// InitHeadRef does a 'Head' on the *git.Repository. This will return a *plumbing.Reference corresponding to the
// git 'HEAD' reference.
func (glifRepo *GlifRepo) InitHeadRef() {
	ref, err := glifRepo.GitRepo.Head()
	if err != nil {
		fmt.Print(err)
		log.Fatal("an error occured while retrieving the HEAD reference")
	}

	glifRepo.HeadRef = ref
}

// FetchAllMatchingTags uses a regex to match tags (by name) in the git repository.
// Used in the following builtin(s):
//	- extractTags
func (glifRepo *GlifRepo) FetchAllMatchingTags(regexString string) bool {
	r, _ := regexp.Compile(regexString)

	//iter, _ := glifRepo.GitRepo.TagObjects()
	iter, err := glifRepo.GitRepo.Tags()

	if err != nil {
		return false
	}

	_ = iter.ForEach(func(reference *plumbing.Reference) error {
		tagObj, err := glifRepo.GitRepo.TagObject(reference.Hash())
		if err == nil && tagObj != nil {
			if r.MatchString(reference.Name().String()) {
				fmt.Printf("Matched tag name (tagged object): %s (%s)\n", reference.Name(), tagObj.Name)
				glifRepo.matchingTags[reference.Name().String()] = tagObj
			} /* else {
				fmt.Printf("Not matched tag name (tagged object): %s (%s)\n", reference.Name(), tagObj.Name)
			}*/
		} else {
			fmt.Printf("No tag object found for: %s (%s)\n", reference.Name(), reference.Target())
		}

		return nil
	})

	for name, reference := range glifRepo.matchingTags {
		reference.Name = name
		glifRepo.tagsLatestToEarliest = append(glifRepo.tagsLatestToEarliest, reference)
	}

	sort.SliceStable(glifRepo.tagsLatestToEarliest, func(i, j int) bool {
		return glifRepo.tagsLatestToEarliest[i].Tagger.When.After(
			glifRepo.tagsLatestToEarliest[j].Tagger.When)
	})

	return true
}

// GetLatestTag returns the appropriate *object.Tag that correspond to the specified offset.
// GetLatestTag(0) would return the very latest tag because it would return the first index of the slice 'tagsLatestToEarliest'
// Used in the following builtin(s):
//	- getLatestTag
func (glifRepo *GlifRepo) GetLatestTag(offset int64) *object.Tag {
	if int(offset) >= len(glifRepo.tagsLatestToEarliest) {
		return nil
	}

	return glifRepo.tagsLatestToEarliest[offset]
}

// GetSpecificTag returns the appropriate *object.Tag that correspond to the specified name
func (glifRepo *GlifRepo) GetSpecificTag(tagName string) *object.Tag {
	if elem, ok := helpers.Contains(glifRepo.tagsLatestToEarliest, tagName); ok {
		return elem
	}

	return nil
}
