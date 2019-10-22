package main

import (
	"fmt"
	"strings"
	"testing"
)

//==============================================================================
// Tests
//==============================================================================

// Test empty Suite
func Test1(t *testing.T) {
	suites := []Suite{Suite{}}
	expected := "alias k='kubectl'"
	test(t, suites, expected)
}

// Test sequence of normal Groups of Segments
func Test2(t *testing.T) {
	segmentA1 := Segment{{Short: "a1", Long: "a1"}}
	segmentA2 := Segment{{Short: "a2", Long: "a2"}}
	segmentB1 := Segment{{Short: "b1", Long: "b1"}}
	segmentB2 := Segment{{Short: "b2", Long: "b2"}}
	groupA := Group{
		[]Segment{segmentA1, segmentA2},
		false,
	}
	groupB := Group{
		[]Segment{segmentB1, segmentB2},
		false,
	}
	suites := []Suite{Suite{groupA, groupB}}
	expected := `
alias ka1='kubectl a1'
alias ka1b1='kubectl a1 b1'
alias ka1b2='kubectl a1 b2'
alias ka2='kubectl a2'
alias ka2b1='kubectl a2 b1'
alias ka2b2='kubectl a2 b2'
`
	test(t, suites, expected)
}

// Test sequence of normal Groups of Segments with mutually exlusive Tokens.
// (This is doesn't make sense in practice. Multiple mutually exlusive Tokens
// per Segment make only sense in Groups with combinable Segments.)
func Test3(t *testing.T) {
	segmentA1 := Segment{
		{Short: "a1", Long: "a1"},
		{Short: "a2", Long: "a2"},
	}
	segmentB1 := Segment{
		{Short: "b1", Long: "b1"},
		{Short: "b2", Long: "b2"},
	}
	groupA := Group{
		[]Segment{segmentA1},
		false,
	}
	groupB := Group{
		[]Segment{segmentB1},
		false,
	}
	suites := []Suite{Suite{groupA, groupB}}
	expected := `
alias ka1='kubectl a1'
alias ka2='kubectl a2'
alias ka1b1='kubectl a1 b1'
alias ka2b1='kubectl a2 b1'
alias ka1b2='kubectl a1 b2'
alias ka2b2='kubectl a2 b2'
`
	test(t, suites, expected)
}

// Test Group of combinable Segments
func Test4(t *testing.T) {
	segmentA1 := Segment{{Short: "a1", Long: "a1"}}
	segmentA2 := Segment{{Short: "a2", Long: "a2"}}
	segmentA3 := Segment{{Short: "a3", Long: "a3"}}
	group := Group{
		[]Segment{segmentA1, segmentA2, segmentA3},
		true,
	}
	suites := []Suite{Suite{group}}
	expected := `
alias ka1='kubectl a1'
alias ka2='kubectl a2'
alias ka3='kubectl a3'
alias ka1a2='kubectl a1 a2'
alias ka2a1='kubectl a2 a1'
alias ka1a3='kubectl a1 a3'
alias ka3a1='kubectl a3 a1'
alias ka2a3='kubectl a2 a3'
alias ka3a2='kubectl a3 a2'
alias ka1a2a3='kubectl a1 a2 a3'
alias ka1a3a2='kubectl a1 a3 a2'
alias ka2a1a3='kubectl a2 a1 a3'
alias ka2a3a1='kubectl a2 a3 a1'
alias ka3a1a2='kubectl a3 a1 a2'
alias ka3a2a1='kubectl a3 a2 a1'
`
	test(t, suites, expected)
}

// Test Group of combinable Segments with mutually exclusive Tokens
func Test5(t *testing.T) {
	segmentA1 := Segment{
		{Short: "a1A", Long: "a1A"},
		{Short: "a1B", Long: "a1B"},
	}
	segmentA2 := Segment{{Short: "a2", Long: "a2"}}
	group := Group{
		[]Segment{segmentA1, segmentA2},
		true,
	}
	suites := []Suite{Suite{group}}
	expected := `
alias ka1A='kubectl a1A'
alias ka1B='kubectl a1B'
alias ka2='kubectl a2'
alias ka1Aa2='kubectl a1A a2'
alias ka1Ba2='kubectl a1B a2'
alias ka2a1A='kubectl a2 a1A'
alias ka2a1B='kubectl a2 a1B'
`
	test(t, suites, expected)
}

// Test sequence of multiple Groups of combinable Segments
func Test6(t *testing.T) {
	segmentA1 := Segment{{Short: "a1", Long: "a1"}}
	segmentA2 := Segment{{Short: "a2", Long: "a2"}}
	segmentB1 := Segment{{Short: "b1", Long: "b1"}}
	segmentB2 := Segment{{Short: "b2", Long: "b2"}}
	groupA := Group{
		[]Segment{segmentA1, segmentA2},
		true,
	}
	groupB := Group{
		[]Segment{segmentB1, segmentB2},
		true,
	}
	suites := []Suite{Suite{groupA, groupB}}
	expected := `
alias ka1='kubectl a1'
alias ka2='kubectl a2'
alias ka1a2='kubectl a1 a2'
alias ka2a1='kubectl a2 a1'
alias ka1b1='kubectl a1 b1'
alias ka1b2='kubectl a1 b2'
alias ka1b1b2='kubectl a1 b1 b2'
alias ka1b2b1='kubectl a1 b2 b1'
alias ka2b1='kubectl a2 b1'
alias ka2b2='kubectl a2 b2'
alias ka2b1b2='kubectl a2 b1 b2'
alias ka2b2b1='kubectl a2 b2 b1'
alias ka1a2b1='kubectl a1 a2 b1'
alias ka1a2b2='kubectl a1 a2 b2'
alias ka1a2b1b2='kubectl a1 a2 b1 b2'
alias ka1a2b2b1='kubectl a1 a2 b2 b1'
alias ka2a1b1='kubectl a2 a1 b1'
alias ka2a1b2='kubectl a2 a1 b2'
alias ka2a1b1b2='kubectl a2 a1 b1 b2'
alias ka2a1b2b1='kubectl a2 a1 b2 b1'
`
	test(t, suites, expected)
}

//==============================================================================
// Helper functions
//==============================================================================

// Compare the expected list of aliases to the generated list of aliases. The
// order of the aliases doesn't matter. For example, if  expected is "A\nB\nC"
// and generated is "C\nB\nA", the test succeeds. Leading and trailing newlines
// are stripped from both expected and generated.
func test(t *testing.T, suites []Suite, expected string) {
	expectedStr := strings.Trim(expected, "\n")
	expectedArr := strings.Split(expectedStr, "\n")

	var out strings.Builder
	generateAliases(suites, &out)
	actualStr := strings.Trim(out.String(), "\n")
	actualArr := strings.Split(actualStr, "\n")

	diff := getDiff(expectedArr, actualArr)

	// If diff is non-empty, expected and actual contain differing sets of aliases
	if len(diff) != 0 {
		var s strings.Builder
		s.WriteString("\n")
		for alias, v := range diff {
			if v > 0 {
				s.WriteString(fmt.Sprintf("In expected but not in actual: %s\n", alias))
			} else if v < 0 {
				s.WriteString(fmt.Sprintf("In actual but not in expected: %s\n", alias))
			}
		}
		t.Errorf(s.String())
	}
}

// Return the difference between two string arrays by looking at them as sets
// (["a", "b", "c"] is the same as ["c", "b", "a"]). The returned map contains
// an entry for each element that is in setA but not in setB (value > 0) and
// that is in setB but not in setA (value < 0).
func getDiff(setA []string, setB []string) map[string]int {
	diff := make(map[string]int)
	for _, a := range setA {
		diff[a]++
	}
	for _, b := range setB {
		diff[b]--
		if diff[b] == 0 {
			delete(diff, b)
		}
	}
	return diff
}
