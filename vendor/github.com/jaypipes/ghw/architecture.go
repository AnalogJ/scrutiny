//
// Use and distribution licensed under the Apache license version 2.
//
// See the COPYING file in the root project directory for full text.
//

package ghw

import "strings"

// Architecture describes the overall hardware architecture. It can be either
// Symmetric Multi-Processor (SMP) or Non-Uniform Memory Access (NUMA)
type Architecture int

const (
	// SMP is a Symmetric Multi-Processor system
	ARCHITECTURE_SMP Architecture = iota
	// NUMA is a Non-Uniform Memory Access system
	ARCHITECTURE_NUMA
)

var (
	architectureString = map[Architecture]string{
		ARCHITECTURE_SMP:  "SMP",
		ARCHITECTURE_NUMA: "NUMA",
	}
)

func (a Architecture) String() string {
	return architectureString[a]
}

// NOTE(jaypipes): since serialized output is as "official" as we're going to
// get, let's lowercase the string output when serializing, in order to
// "normalize" the expected serialized output
func (a Architecture) MarshalJSON() ([]byte, error) {
	return []byte("\"" + strings.ToLower(a.String()) + "\""), nil
}
