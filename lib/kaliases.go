package kaliases

import (
	"fmt"
	"io"
	"os"
)

/* Types */

type Token struct {
	Short, Long string
}

type Segment []Token

type Group struct {
	Segments          []Segment
	AllowCombinations bool
}

type Suite []Group

// Generate aliases from a list of Suites, ensuring that no two aliases have
// the same name, and writing them to the provided writer.
func Generate(suites []Suite, out io.Writer) {
	aliases := map[string]string{}
	for _, suite := range suites {
		generateImpl(suite, 0, []Segment{}, aliases, out)
	}
}

func generateImpl(suite Suite, i int, stack []Segment, aliases map[string]string, out io.Writer) {
	if len(suite) == 0 {
		writeAlias([]Token{}, aliases, out)
		return
	}
	if i == len(suite) {
		return
	}
	// Group of non-combinable Segments
	if !suite[i].AllowCombinations {
		for _, token := range suite[i].Segments {
			stackNew := append(stack, token)
			for _, alternative := range getAlternatives(stackNew) {
				writeAlias(alternative, aliases, out)
			}
			generateImpl(suite, i+1, stackNew, aliases, out)
		}
	} else {
		// Group of combinable Segments. All permutations of all subsets:
		// https://www.wolframalpha.com/input/?i=sum%28n+choose+k%29+*+k%21%2Ck%3D1+to+n
		for _, subset := range getSubsets(suite[i].Segments) {
			for _, permutation := range getPermutations(subset) {
				stackNew := append(stack, permutation...)
				for _, alternative := range getAlternatives(stackNew) {
					writeAlias(alternative, aliases, out)
				}
				generateImpl(suite, i+1, stackNew, aliases, out)
			}
		}
	}
}

// Write an alias definition to the provided writer. If an alias with the same
// name already exists, raise an error.
func writeAlias(tokens []Token, aliases map[string]string, out io.Writer) {
	alias, command := "k", "kubectl"
	for _, token := range tokens {
		alias += token.Short
		command += " " + token.Long
	}
	if _, exists := aliases[alias]; exists {
		fmt.Fprintf(os.Stderr, "\033[31m")
		fmt.Fprintf(os.Stderr, "Error: conflicting aliases:\n")
		fmt.Fprintf(os.Stderr, "  Existing: alias %s='%s'\n", alias, aliases[alias])
		fmt.Fprintf(os.Stderr, "  New:      alias %s='%s'\n", alias, command)
		fmt.Fprintf(os.Stderr, "\033[m")
		os.Exit(1)
	}
	aliases[alias] = command
	fmt.Fprintf(out, "alias %s='%s'\n", alias, command)
}

// Get final sequences of Tokens by expanding any mutually exclusive Tokens
func getAlternatives(tokens []Segment) [][]Token {
	result := [][]Token{}
	c := make(chan []Token)
	go getAlternativesImpl(tokens, 0, []Token{}, c)
	// Number of expansions (for knowing how many values to receive from channel)
	total := 1
	for _, t := range tokens {
		total = total * len(t)
	}
	for i := 0; i < total; i++ {
		result = append(result, <-c)
	}
	return result
}
func getAlternativesImpl(tokens []Segment, i int, stack []Token, c chan []Token) {
	if i == len(tokens) {
		stackCopy := make([]Token, len(stack))
		copy(stackCopy, stack)
		c <- stackCopy
	} else {
		for _, pair := range tokens[i] {
			getAlternativesImpl(tokens, i+1, append(stack, pair), c)
		}
	}
}

// Get all subsets of size > 0 (including the set itself) of a set of Segments
func getSubsets(set []Segment) [][]Segment {
	c := make(chan []Segment)
	go getSubsetsImpl(set, c)
	result := [][]Segment{}
	for subset := range c {
		result = append(result, subset)
	}
	return result
}
func getSubsetsImpl(set []Segment, c chan []Segment) {
	n := len(set)
	// (2^n)-1 subsets represented by binary numbers i from 1 to (2^n)-1
	for i := 1; i < 1<<uint(n); i++ {
		subset := []Segment{}
		// Append elements at indices where the binary number i is 1
		for j := 0; j < n; j++ {
			if 1<<uint(j)&i > 0 {
				subset = append(subset, set[j])
			}
		}
		c <- subset
	}
	close(c)
}

// Get all permutations of a set of Segments
func getPermutations(a []Segment) [][]Segment {
	c := make(chan []Segment)
	go getPermutationsImpl(a, len(a), c)
	result := [][]Segment{}
	for permutation := range c {
		result = append(result, permutation)
	}
	return result
}

// Heap's algorithm, see https://en.wikipedia.org/wiki/Heap%27s_algorithm
func getPermutationsImpl(a []Segment, k int, c chan []Segment) {
	if k == 1 {
		aCopy := make([]Segment, len(a))
		copy(aCopy, a)
		c <- aCopy
	} else {
		getPermutationsImpl(a, k-1, c)
		for i := 0; i < k-1; i++ {
			if k%2 == 0 {
				a[i], a[k-1] = a[k-1], a[i]
			} else {
				a[0], a[k-1] = a[k-1], a[0]
			}
			getPermutationsImpl(a, k-1, c)
		}
	}
	if k == len(a) {
		close(c)
	}
}
