package xcontext

import "strings"

type Metadata map[string]string

func New() Metadata {
	return Metadata{}
}

//
func Pairs(kv ...string) Metadata {
	if len(kv)%2 != 0 {
		panic("xcontext.Pairs: length of kv must be even")
	}
	m := make(Metadata)
	for i := 0; i < len(kv); i += 2 {
		m[kv[i]] = kv[i+1]
	}
	return m
}

// Join joins any number of mds into a single MD.
//
// The order of values for each key is determined by the order in which the mds
// containing those values are presented to Join.
func Join(mds ...Metadata) Metadata {
	out := Metadata{}
	for _, md := range mds {
		for k, v := range md {
			out[k] = v
		}
	}
	return out
}

// Len returns the number of items in metadata.
func (m Metadata) Len() int {
	return len(m)
}

// Get obtains the values for a given key.
//
// k is converted to lowercase before searching in md.
func (m Metadata) Get(k string) string {
	k = strings.ToLower(k)
	return m[k]
}

// Set sets the value of a given key with a slice of values.
//
// k is converted to lowercase before storing in md.
func (m Metadata) Set(k string, val string) {
	k = strings.ToLower(k)
	m[k] = val
}

// Copy returns a copy of metadata.
func (m Metadata) Copy() Metadata {
	return Join(m)
}

// Append adds the values to key k, not overwriting what was already stored at
// that key.
//
// k is converted to lowercase before storing in md.
func (m Metadata) Append(k string, val string) {
	k = strings.ToLower(k)
	m[k] = val
}

// Delete removes the values for a given key k which is converted to lowercase
// before removing it from md.
func (m Metadata) Delete(k string) {
	k = strings.ToLower(k)
	delete(m, k)
}
