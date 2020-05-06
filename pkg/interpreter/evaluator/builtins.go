package evaluator

import (
	"bytes"
	"fmt"
	"github.com/TurnsCoffeeIntoScripts/git-log-issue-finder/pkg/interpreter/object"
	"github.com/go-git/go-git/v5"
	gitobject "github.com/go-git/go-git/v5/plumbing/object"
	"regexp"
	"strings"
)

var builtins = map[string]*object.Builtin{
	"len": {
		Fn: func(args ...object.Object) object.Object {
			if len(args) != 1 {
				return newError("wrong number of arguments. got=%d, want=1", len(args))
			}

			switch arg := args[0].(type) {
			case *object.Array:
				return &object.Integer{Value: int64(len(arg.Elements))}
			case *object.String:
				return &object.Integer{Value: int64(len(arg.Value))}
			default:
				return newError("argument to `len` not supported, got %s", args[0].Type())
			}
		},
	},
	"first": {
		Fn: func(args ...object.Object) object.Object {
			if len(args) != 1 {
				return newError("wrong number of arguments. got=%d, want=1", len(args))
			}

			if args[0].Type() != object.ArrayObj {
				return newError("argument to `first` must be ARRAY, got %s", args[0].Type())
			}

			arr := args[0].(*object.Array)
			if len(arr.Elements) > 0 {
				return arr.Elements[0]
			}

			return NULL
		},
	},
	"last": {
		Fn: func(args ...object.Object) object.Object {
			if len(args) != 1 {
				return newError("wrong number of arguments. got=%d, want=1", len(args))
			}

			if args[0].Type() != object.ArrayObj {
				return newError("argument to `last` must be ARRAY, got %s",
					args[0].Type())
			}

			arr := args[0].(*object.Array)
			length := len(arr.Elements)
			if length > 0 {
				return arr.Elements[length-1]
			}

			return NULL
		},
	},
	"rest": {
		Fn: func(args ...object.Object) object.Object {
			if len(args) != 1 {
				return newError("wrong number of arguments. got=%d, want=1",
					len(args))
			}

			if args[0].Type() != object.ArrayObj {
				return newError("argument to `rest` must be ARRAY, got %s",
					args[0].Type())
			}

			arr := args[0].(*object.Array)
			length := len(arr.Elements)
			if length > 0 {
				newElements := make([]object.Object, length-1, length-1)
				copy(newElements, arr.Elements[1:length])
				return &object.Array{Elements: newElements}
			}

			return NULL
		},
	},
	"push": {
		Fn: func(args ...object.Object) object.Object {
			if len(args) != 2 {
				return newError("wrong number of arguments. got=%d, want=2",
					len(args))
			}

			if args[0].Type() != object.ArrayObj {
				return newError("argument to `push` must be ARRAY, got %s",
					args[0].Type())
			}

			arr := args[0].(*object.Array)
			length := len(arr.Elements)
			newElements := make([]object.Object, length+1, length+1)
			copy(newElements, arr.Elements)
			newElements[length] = args[1]
			return &object.Array{Elements: newElements}
		},
	},
	"print": {
		Fn: func(args ...object.Object) object.Object {
			for _, arg := range args {
				fmt.Println(arg.Inspect())
			}

			return NULL
		},
	},
	"initRepo": {
		Fn: func(args ...object.Object) object.Object {
			if len(args) != 1 {
				return newError("wrong number of arguments. got=%d, want=0", len(args)-1)
			}

			repoObj := &object.Repo{Path: args[0]}
			repoObj.Repo.Open(args[0].Inspect())
			repoObj.Repo.InitHeadRef()

			return repoObj
		},
		RequireEnv: true,
		EnvName:    "repopath",
	},
	"whichRepo": {
		Fn: func(args ...object.Object) object.Object {
			if len(args) != 1 {
				return newError("wrong number of arguments. got=%d, want=0", len(args)-1)
			}

			return &object.String{Value: args[0].Inspect()}
		},
		RequireEnv: true,
		EnvName:    "repopath",
	},
	"extractTags": {
		Fn: func(args ...object.Object) object.Object {
			if len(args) != 2 {
				return newError("wrong number of arguments. got=%d, want=2", len(args))
			}

			repo, ok := args[0].(*object.Repo)
			if !ok {
				return newError("Unable to convert args[0] to *object.Repo while executing 'extractTags'")
			}

			version, ok := args[1].(*object.String)
			if !ok {
				return newError("Unable to convert args[1] to *object.String while executing 'extractTags'")
			}

			var buffer bytes.Buffer
			for _, c := range version.Inspect() {
				cc := string(c)

				if cc == "$" {
					buffer.WriteString("([0-9]+)")
				} else if cc == "." {
					buffer.WriteString("\\.")
				} else if cc == "*" {
					buffer.WriteString(".*")
				} else if cc == "+" {
					buffer.WriteString(".+")
				} else {
					buffer.WriteString(cc)
				}
			}

			buffer.WriteString("$")

			if ok := repo.Repo.FetchAllMatchingTags(buffer.String()); !ok {
				return newError("Failed to fetch tag (plumbing.Reference) on repo")
			}

			return NULL
		},
	},
	"getTag": {
		Fn: func(args ...object.Object) object.Object {
			if len(args) != 2 {
				return newError("wrong number of arguments. got=%d, want=2", len(args))
			}

			repo, ok := args[0].(*object.Repo)
			if !ok {
				return newError("Unable to convert args[0] to *object.Repo while executing 'getTag'")
			}

			tagName, ok := args[1].(*object.String)
			if !ok {
				return newError("Unable to convert args[1] to *object.String while executing 'getTag'")
			}

			tag := repo.Repo.GetSpecificTag(tagName.Value)
			if tag == nil {
				return NULL
			}

			value := &object.String{Value: tag.Name}
			tagObj := &object.Tag{Value: value, Tag: tag}
			return tagObj
		},
	},
	"getLatestTag": {
		Fn: func(args ...object.Object) object.Object {
			if len(args) != 2 {
				return newError("wrong number of arguments. got=%d, want=2", len(args))
			}

			repo, ok := args[0].(*object.Repo)
			if !ok {
				return newError("Unable to convert args[0] to *object.Repo while executing 'getLatestTag'")
			}

			integer, ok := args[1].(*object.Integer)
			if !ok {
				return newError("Unable to convert args[1] to *object.Integer while executing 'getLatestTag'")
			}

			tag := repo.Repo.GetLatestTag(integer.Value)
			if tag == nil {
				return NULL
			}

			value := &object.String{Value: tag.Name}
			pr := &object.Tag{Value: value, Tag: tag}
			return pr
		},
	},
	"diff": {
		Fn: func(args ...object.Object) object.Object {
			if len(args) != 4 {
				return newError("wrong number of arguments. got=%d, want=3", len(args)-1)
			}

			ticketRegex, ok := args[0].(*object.String)
			if !ok {
				return newError("Unable to convert args[0] to *object.String while executing 'diff'")
			}

			repo, ok := args[1].(*object.Repo)
			if !ok {
				return newError("Unable to convert args[1] to *object.Repo while executing 'diff'")
			}

			from, ok := args[2].(*object.Tag)
			if !ok {
				return newError("Unable to convert args[2] ('from') to *object.Tag while executing 'diff'")
			}

			to, ok := args[3].(*object.Tag)
			if !ok {
				return newError("Unable to convert args[3] ('to') to *object.Tag while executing 'diff'")
			}

			// ---------------------------------------------
			// Extract the log iterator from the 'from' hash
			// ---------------------------------------------
			commitFromTag, _ := from.Tag.Commit()
			iterFrom, errFrom := repo.Repo.GitRepo.Log(&git.LogOptions{From: commitFromTag.Hash})
			if errFrom != nil {
				return newError("an error occured while retrieving the commit history from 'fromHash'")
			}

			// -------------------------------------------
			// Extract the log iterator from the 'to' hash
			// -------------------------------------------
			commitToTag, _ := to.Tag.Commit()
			iterTo, errTo := repo.Repo.GitRepo.Log(&git.LogOptions{From: commitToTag.Hash})
			if errTo != nil {
				return newError("an error occured while retrieving the commit history from 'toHash'")
			}

			fmt.Printf("Performing diff on %s --> %s\n", from.Tag.Name, to.Tag.Name)

			// ----------------------------
			// Initialize the commit slices
			// ----------------------------
			commitFromSlice := make([]*gitobject.Commit, 0)
			commitToSlice := make([]*gitobject.Commit, 0)

			_ = iterFrom.ForEach(func(commit *gitobject.Commit) error {
				commitFromSlice = append(commitFromSlice, commit)
				return nil
			})

			_ = iterTo.ForEach(func(commit *gitobject.Commit) error {
				commitToSlice = append(commitToSlice, commit)
				return nil
			})

			// ------------------------------------------------
			// Perform the actual diff operation on both slices
			// ------------------------------------------------
			diff := make([]*gitobject.Commit, 0)
			for i := 0; i < 2; i++ {
				for _, s1 := range commitToSlice {
					found := false
					for _, s2 := range commitFromSlice {
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
					commitToSlice, commitFromSlice = commitFromSlice, commitToSlice
				}
			}

			// ------------------------------------------------
			// Look in the commits and append when matching
			// ------------------------------------------------
			var TicketSlice []string
			for _, c := range diff {
				if presentInMessage, ticket := tickets(c.Message, ticketRegex.Value); presentInMessage {
					TicketSlice = append(TicketSlice, ticket...)
				}
			}

			TicketSlice = unique(TicketSlice)
			fmt.Println(TicketSlice)

			return NULL
		},
		RequireEnv: true,
		EnvName:    "tickets",
	},
}

func tickets(text, ticketRegex string) (bool, []string) {
	regex := "((?:"
	if ticketRegex == "*" {
		regex += "[a-zA-Z0-9]+"
	} else {
		regex += strings.ReplaceAll(ticketRegex, ",", "|")
	}
	regex += ")-[0-9]+)"

	r, _ := regexp.Compile(regex)

	out := r.FindAllString(text, -1)

	if len(out) == 0 {
		return false, []string{}
	}

	return true, out
}

func unique(s []string) []string {
	u := make([]string, 0, len(s))
	m := make(map[string]bool)

	for _, val := range s {
		if _, ok := m[val]; !ok {
			m[val] = true
			u = append(u, val)
		}
	}

	return u
}
