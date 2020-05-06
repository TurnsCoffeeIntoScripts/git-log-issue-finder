package helpers

import (
	"github.com/go-git/go-git/v5/plumbing/object"
	"strings"
)

func Contains(arr []*object.Tag, elemName string) (*object.Tag, bool) {
	for idx := range arr {
		if strings.Contains(arr[idx].Name, elemName) {
			return arr[idx], true
		}
	}

	return nil, false
}
