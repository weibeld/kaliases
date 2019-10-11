package main

import "fmt"

/* Types */

type Token struct {
	Short, Long string
}

type Group struct {
	Tokens            []Token
	AllowCombinations bool
}

type Set []Group

//==============================================================================
// Declarative part (define your desired behaviour here)
//==============================================================================

/* Sets */

var sets = []Set{
	Set{actionsGet, resources, optionsGet},
	Set{actionsDelete, resources, optionsDelete},
	Set{actionsDescribe, resources, optionsDescribe},
	Set{actionsExec, optionsExec},
	Set{actionsLogs, optionsLogs},
	Set{actionsEdit, resources},
	Set{actionsOther},
}

/* Groups */

var resources = Group{
	[]Token{resourcePod, resourceDeployment, resourceService, resourceNode,
		resourceIngress, resourceRole, resourceRoleBinding, resourceClusterRole,
		resourceClusterRoleBinding},
	false,
}

var actionsGet = Group{
	[]Token{actionGet},
	false,
}
var optionsGet = Group{
	[]Token{optionWatch, optionOutput, optionAllNamespaces},
	true,
}

var actionsDelete = Group{
	[]Token{actionDelete},
	false,
}
var optionsDelete = Group{
	[]Token{optionAllDelete, optionAllNamespacesDelete},
	true,
}

var actionsDescribe = Group{
	[]Token{actionDescribe},
	false,
}
var optionsDescribe = Group{
	[]Token{optionAllNamespaces},
	true,
}

var actionsExec = Group{
	[]Token{actionExec},
	false,
}
var optionsExec = Group{
	[]Token{optionInteractive},
	true,
}

var actionsLogs = Group{
	[]Token{actionLogs},
	false,
}
var optionsLogs = Group{
	[]Token{optionFollow},
	true,
}

var actionsEdit = Group{
	[]Token{actionEdit},
	false,
}

var actionsOther = Group{
	[]Token{actionApply, actionPortForward, actionExplain},
	false,
}

/* Action tokens */

var actionGet = Token{
	Short: "g", Long: "get",
}
var actionDelete = Token{
	Short: "d", Long: "delete",
}
var actionDescribe = Token{
	Short: "s", Long: "decribe",
}
var actionEdit = Token{
	Short: "e", Long: "edit",
}
var actionExec = Token{
	Short: "x", Long: "exec",
}
var actionLogs = Token{
	Short: "l", Long: "logs",
}
var actionApply = Token{
	Short: "a", Long: "apply",
}
var actionPortForward = Token{
	Short: "p", Long: "port-forward",
}
var actionExplain = Token{
	Short: "ex", Long: "explain",
}

/* Resource tokens */

var resourcePod = Token{
	Short: "p", Long: "pod",
}
var resourceDeployment = Token{
	Short: "d", Long: "deployment",
}
var resourceService = Token{
	Short: "s", Long: "service",
}
var resourceNode = Token{
	Short: "n", Long: "node",
}
var resourceIngress = Token{
	Short: "i", Long: "ingress",
}
var resourceRole = Token{
	Short: "r", Long: "role",
}
var resourceRoleBinding = Token{
	Short: "rb", Long: "rolebinding",
}
var resourceClusterRole = Token{
	Short: "cr", Long: "clusterrole",
}
var resourceClusterRoleBinding = Token{
	Short: "crb", Long: "clusterrolebinding",
}

/* Option tokens */

var optionWatch = Token{
	Short: "w", Long: "-w",
}
var optionOutput = Token{
	Short: "y", Long: "-o yaml",
}
var optionAllNamespaces = Token{
	Short: "a", Long: "--all-namespaces",
}
var optionAllDelete = Token{
	Short: "a", Long: "--all",
}
var optionAllNamespacesDelete = Token{
	Short: "A", Long: "--all-namespaces",
}
var optionInteractive = Token{
	Short: "i", Long: "-it",
}
var optionFollow = Token{
	Short: "f", Long: "-f",
}

//==============================================================================
// Imperative part
//==============================================================================

func main() {
	for _, s := range sets {
		generate(s)
	}
}

// Generate all aliases of a Set
func generate(set Set) {
	generateImpl(set, 0, []Token{})
}
func generateImpl(set Set, i int, stack []Token) {
	if i == len(set) {
		return
	}
	// Group that can't be combined with each other
	if !set[i].AllowCombinations {
		for _, token := range set[i].Tokens {
			stackNew := append(stack, token)
			printAlias(stackNew)
			generateImpl(set, i+1, stackNew)
		}
	} else {
		// Group that can be combined with each other. Create an alias for each per-
		// mutation of each subset (size > 0) of the Group. Total number of aliases:
		// https://www.wolframalpha.com/input/?i=sum%28n+choose+k%29+*+k%21%2Ck%3D1+to+n
		for _, subset := range getSubsets(set[i].Tokens) {
			for _, permutation := range getPermutations(subset) {
				stackNew := append(stack, permutation...)
				printAlias(stackNew)
				generateImpl(set, i+1, stackNew)
			}
		}
	}
}

// Print a single alias definition given its sequence of Tokens
func printAlias(tokens []Token) {
	alias, command := "k", "kubectl"
	for _, token := range tokens {
		alias += token.Short
		command += " " + token.Long
	}
	line := fmt.Sprintf("alias %s='%s'\n", alias, command)
	fmt.Print(line)
}

// Get all subsets of size > 0 from a set (including the set itself)
func getSubsets(set []Token) [][]Token {
	c := make(chan []Token)
	go getSubsetsImpl(set, c)
	result := [][]Token{}
	for subset := range c {
		result = append(result, subset)
	}
	return result
}
func getSubsetsImpl(set []Token, c chan []Token) {
	n := len(set)
	// (2^n)-1 subsets represented by binary numbers i from 1 to (2^n)-1
	for i := 1; i < 1<<uint(n); i++ {
		subset := []Token{}
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

// Get all permutations of a set
func getPermutations(a []Token) [][]Token {
	c := make(chan []Token)
	go getPermutationsImpl(a, len(a), c)
	result := [][]Token{}
	for permutation := range c {
		result = append(result, permutation)
	}
	return result
}

// Heap's algorithm, see https://en.wikipedia.org/wiki/Heap%27s_algorithm
func getPermutationsImpl(a []Token, k int, c chan []Token) {
	if k == 1 {
		x := make([]Token, len(a))
		copy(x, a)
		c <- x
	}
	for i := 0; i < k; i++ {
		getPermutationsImpl(a, k-1, c)
		if k%2 == 0 {
			a[i], a[k-1] = a[k-1], a[i]
		} else {
			a[0], a[k-1] = a[k-1], a[0]
		}
	}
	if k == len(a) {
		close(c)
	}
}
