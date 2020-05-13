package object

import (
	"github.com/go-git/go-git/v5/plumbing/object"
)

// Tag is a wrapper for the interpreter of the *object.Tag (go-git) object
type Tag struct {
	Value Object
	Tag   *object.Tag
}

// Type returns TagObj (TAG)
func (t *Tag) Type() Type {
	return TagObj
}

// Inspect the value which is the name of the git tag
func (t *Tag) Inspect() string {
	return t.Value.Inspect()
}
