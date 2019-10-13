package main

import (
	"fmt"
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

//==============================================================================
// Declarative part (define your desired behaviour here)
//==============================================================================

/* Suites */

var suites = []Suite{
	Suite{actionsGet, resources, optionsGet},
	Suite{actionsDelete, resources, optionsDelete},
	Suite{actionsDescribe, resources, optionsDescribe},
	Suite{actionsExec, optionsExec},
	Suite{actionsLogs, optionsLogs},
	Suite{actionsEdit, resources},
	Suite{actionsOther},
}

/* Groups */

var resources = Group{
	[]Segment{resourcePod, resourceDeployment, resourceService, resourceNode,
		resourceIngress, resourceRole, resourceRoleBinding, resourceClusterRole,
		resourceClusterRoleBinding},
	false,
}

var actionsGet = Group{
	[]Segment{actionGet},
	false,
}
var optionsGet = Group{
	[]Segment{optionWatch, optionOutput, optionAllNamespaces},
	true,
}

var actionsDelete = Group{
	[]Segment{actionDelete},
	false,
}
var optionsDelete = Group{
	[]Segment{optionAll, optionAllNamespaces},
	true,
}

var actionsDescribe = Group{
	[]Segment{actionDescribe},
	false,
}
var optionsDescribe = Group{
	[]Segment{optionAllNamespaces},
	true,
}

var actionsExec = Group{
	[]Segment{actionExec},
	false,
}
var optionsExec = Group{
	[]Segment{optionInteractive},
	true,
}

var actionsLogs = Group{
	[]Segment{actionLogs},
	false,
}
var optionsLogs = Group{
	[]Segment{optionFollow},
	true,
}

var actionsEdit = Group{
	[]Segment{actionEdit},
	false,
}

var actionsOther = Group{
	[]Segment{actionApply, actionPortForward, actionExplain},
	false,
}

/* Action segments */

var actionGet = Segment{
	{Short: "g", Long: "get"},
}
var actionDelete = Segment{
	{Short: "d", Long: "delete"},
}
var actionDescribe = Segment{
	{Short: "s", Long: "decribe"},
}
var actionEdit = Segment{
	{Short: "e", Long: "edit"},
}
var actionExec = Segment{
	{Short: "x", Long: "exec"},
}
var actionLogs = Segment{
	{Short: "l", Long: "logs"},
}
var actionApply = Segment{
	{Short: "a", Long: "apply"},
}
var actionPortForward = Segment{
	{Short: "p", Long: "port-forward"},
}
var actionExplain = Segment{
	{Short: "ex", Long: "explain"},
}

/* Resource segments */

var resourcePod = Segment{
	{Short: "p", Long: "pod"},
}
var resourceDeployment = Segment{
	{Short: "d", Long: "deployment"},
}
var resourceService = Segment{
	{Short: "s", Long: "service"},
}
var resourceNode = Segment{
	{Short: "n", Long: "node"},
}
var resourceIngress = Segment{
	{Short: "i", Long: "ingress"},
}
var resourceRole = Segment{
	{Short: "r", Long: "role"},
}
var resourceRoleBinding = Segment{
	{Short: "rb", Long: "rolebinding"},
}
var resourceClusterRole = Segment{
	{Short: "cr", Long: "clusterrole"},
}
var resourceClusterRoleBinding = Segment{
	{Short: "crb", Long: "clusterrolebinding"},
}

/* Option segments */

var optionWatch = Segment{
	{Short: "w", Long: "-w"},
}
var optionOutput = Segment{
	{Short: "y", Long: "-o yaml"},
	{Short: "j", Long: "-o json"},
}
var optionAllNamespaces = Segment{
	{Short: "a", Long: "--all-namespaces"},
}
var optionAll = Segment{
	{Short: "A", Long: "--all"},
}
var optionInteractive = Segment{
	{Short: "i", Long: "-it"},
}
var optionFollow = Segment{
	{Short: "f", Long: "-f"},
}

//==============================================================================
// Imperative part
//==============================================================================

func main() {
	for _, suite := range suites {
		generate(suite)
	}
}

// Generate all aliases of a Suite
func generate(suite Suite) {
	generateImpl(suite, 0, []Segment{})
}
func generateImpl(suite Suite, i int, stack []Segment) {
	if i == len(suite) {
		return
	}
	// Group of non-combinable Segments
	if !suite[i].AllowCombinations {
		for _, token := range suite[i].Segments {
			stackNew := append(stack, token)
			for _, alternative := range getAlternatives(stackNew) {
				printAlias(alternative)
			}
			generateImpl(suite, i+1, stackNew)
		}
	} else {
		// Group of combinable Segments. All permutations of all subsets:
		// https://www.wolframalpha.com/input/?i=sum%28n+choose+k%29+*+k%21%2Ck%3D1+to+n
		for _, subset := range getSubsets(suite[i].Segments) {
			for _, permutation := range getPermutations(subset) {
				stackNew := append(stack, permutation...)
				for _, alternative := range getAlternatives(stackNew) {
					printAlias(alternative)
				}
				generateImpl(suite, i+1, stackNew)
			}
		}
	}
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

// Aliases generated so far (for detecting name clashes)
var aliases = map[string]string{}

// Print a single alias definition given its sequence of Segments
func printAlias(pairs []Token) {
	alias, command := "k", "kubectl"
	for _, pair := range pairs {
		alias += pair.Short
		command += " " + pair.Long
	}
	if _, exists := aliases[alias]; exists {
		fmt.Printf("\033[31m")
		fmt.Printf("Error: conflicting aliases:\n")
		fmt.Printf("  Existing: alias %s='%s'\n", alias, aliases[alias])
		fmt.Printf("  New:      alias %s='%s'\n", alias, command)
		fmt.Printf("\033[m")
		os.Exit(1)
	}
	aliases[alias] = command
	line := fmt.Sprintf("alias %s='%s'\n", alias, command)
	fmt.Print(line)
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
