package find

import (
	"fmt"
	"gopkg.in/src-d/go-git.v4/plumbing"
	"gopkg.in/src-d/go-git.v4/plumbing/object"
	"log"
)

// LatestTag returns the Hash value of the lastest Tag found based on specified iterator
func LatestTag(iter *object.TagIter, hash plumbing.Hash, err error) plumbing.Hash {
	if err != nil {
		fmt.Print(err.Error())
		log.Fatal("an error occured while retrieving the tag list")
		return hash
	}

	tagsName := make(map[string]*object.Tag)
	err = iter.ForEach(func(tag *object.Tag) error {
		tagsName[tag.Name] = tag
		return nil
	})

	var latestTag *object.Tag
	for _, v := range tagsName {
		if latestTag == nil {
			latestTag = v
		}

		if latestTag.Tagger.When.Before(v.Tagger.When) {
			latestTag = v
		}
	}

	return latestTag.Target
}
