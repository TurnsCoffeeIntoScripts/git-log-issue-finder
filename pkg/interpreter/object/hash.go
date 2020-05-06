package object

import (
	"bytes"
	"fmt"
	"hash/fnv"
	"strings"
)

// Hashable is an interface to be implemented by object that can by hashed
type Hashable interface {
	HashKey() HashKey
}

// HashKey is the value (unsigned integer) of the hash key
type HashKey struct {
	Type  Type
	Value uint64
}

// HashPair is an object that can be inserted into the Hash construct
type HashPair struct {
	Key   Object
	Value Object
}

// Hash is a simple map, mapping HashKey to HashPair
type Hash struct {
	Pairs map[HashKey]HashPair
}

// HashKey returns the hashkey of a boolean
// TODO optimize perf by caching their return values
func (b *Boolean) HashKey() HashKey {
	var value uint64

	if b.Value {
		value = 1
	} else {
		value = 0
	}

	return HashKey{Type: b.Type(), Value: value}
}

// HashKey returns the hashkey of an integer
// TODO optimize perf by caching their return values
func (i *Integer) HashKey() HashKey {
	return HashKey{Type: i.Type(), Value: uint64(i.Value)}
}

// HashKey returns the hashkey of a string
// TODO optimize perf by caching their return values
func (s *String) HashKey() HashKey {
	h := fnv.New64a()
	h.Write([]byte(s.Value))

	return HashKey{Type: s.Type(), Value: h.Sum64()}
}

// Type returns HashObj (HASH)
func (h *Hash) Type() Type {
	return HashObj
}

// Inspect returns the string value of the pairs contained in the hash
func (h *Hash) Inspect() string {
	var out bytes.Buffer
	pairs := []string{}

	for _, pair := range h.Pairs {
		pairs = append(pairs, fmt.Sprintf("%s: %s", pair.Key.Inspect(), pair.Value.Inspect()))
	}

	out.WriteString("{")
	out.WriteString(strings.Join(pairs, ", "))
	out.WriteString("}")
	return out.String()

}
