package object

import "github.com/TurnsCoffeeIntoScripts/git-log-issue-finder/pkg/scl"

// Repo is a wrapper for the interpreter of scl.GlifRepo object which itself wraps the *git.Repository
type Repo struct {
	Repo scl.GlifRepo
	Path Object
}

// Type returns RepoObj (REPO)
func (r *Repo) Type() Type {
	return RepoObj
}

// Inspect the specified path (string) of the repo
func (r *Repo) Inspect() string {
	return r.Path.Inspect()
}
