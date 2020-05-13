package object

// Environment is the construct that holds the variables (and associated values) declared by the user
// It also has a reference to its outer environement (if any); allowing for some scoping
type Environment struct {
	store map[string]Object
	outer *Environment
}

// NewEnvironmentWithParams creates a new instance with some predefined values
func NewEnvironmentWithParams(tickets string) *Environment {
	env := NewEnvironment()
	env.Set("repopath", &String{Value: "."})
	env.Set("tickets", &String{Value: tickets})

	return env
}

// NewEnvironment creates new instance with no outer environement (top-level)
func NewEnvironment() *Environment {
	s := make(map[string]Object)

	env := &Environment{store: s, outer: nil}

	return env
}

// NewEnclosedEnvironment creates new instance that's enclosed in an existing one.
func NewEnclosedEnvironment(outer *Environment) *Environment {
	env := NewEnvironment()
	env.outer = outer

	return env
}

// Get returns the value of the specified variable name
func (e *Environment) Get(name string) (Object, bool) {
	obj, ok := e.store[name]
	if !ok && e.outer != nil {
		obj, ok = e.outer.Get(name)
	}
	return obj, ok
}

// Set adds an entry to the environment internal store for a new or existing variable
func (e *Environment) Set(name string, val Object) Object {
	e.store[name] = val
	return val
}
