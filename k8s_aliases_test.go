package main

import (
	"fmt"
	"testing"
)

// Test empty Suite
func ExampleGenerate_foo1() {
	var suite = Suite{}
	generate(suite)
	// Unordered output:
	// alias k='kubectl'
}

// Test normal sequence of Segments
func Test2(*testing.T) {
	var segmentA1 = Segment{{Short: "a1", Long: "a1"}}
	var segmentA2 = Segment{{Short: "a2", Long: "a2"}}
	var segmentB1 = Segment{{Short: "b1", Long: "b1"}}
	var segmentB2 = Segment{{Short: "b2", Long: "b2"}}
	var groupA = Group{
		[]Segment{segmentA1, segmentA2},
		false,
	}
	var groupB = Group{
		[]Segment{segmentB1, segmentB2},
		false,
	}
	var suite = Suite{groupA, groupB}
	generate(suite)
	fmt.Println(aliases)
	// Unordered output:
	// alias ka1='kubectl a1'
	// alias ka1b1='kubectl a1 b1'
	// alias ka1b2='kubectl a1 b2'
	// alias ka2='kubectl a2'
	// alias ka2b1='kubectl a2 b1'
	// alias ka2b2='kubectl a2 b2'
}

func Test1(*testing.T) {
	suite := Suite{}
	//expected := "alias k='kubectl'"
	generate(suite)
	fmt.Println(aliases)
}

// Test normal sequence of Segments
func Test3(*testing.T) {
	var segmentA1 = Segment{{Short: "a1", Long: "a1"}}
	var segmentA2 = Segment{{Short: "a2", Long: "a2"}}
	var segmentB1 = Segment{{Short: "b1", Long: "b1"}}
	var segmentB2 = Segment{{Short: "b2", Long: "b2"}}
	var groupA = Group{
		[]Segment{segmentA1, segmentA2},
		false,
	}
	var groupB = Group{
		[]Segment{segmentB1, segmentB2},
		false,
	}
	var suite = Suite{groupA, groupB}
	generate(suite)
	fmt.Println(aliases)
	// Unordered output:
	// alias ka1='kubectl a1'
	// alias ka1b1='kubectl a1 b1'
	// alias ka1b2='kubectl a1 b2'
	// alias ka2='kubectl a2'
	// alias ka2b1='kubectl a2 b1'
	// alias ka2b2='kubectl a2 b2'
}

// Test Segment with combinations
func ExampleGenerate_foo3() {
	var segmentA1 = Segment{{Short: "a1", Long: "a1"}}
	var segmentA2 = Segment{{Short: "a2", Long: "a2"}}
	var segmentA3 = Segment{{Short: "a3", Long: "a3"}}
	var groupA = Group{
		[]Segment{segmentA1, segmentA2, segmentA3},
		true,
	}
	var suite = Suite{groupA}
	generate(suite)
	// Unordered output:
	// alias ka1='kubectl a1'
	// alias ka2='kubectl a2'
	// alias ka3='kubectl a3'
	// alias ka1a2='kubectl a1 a2'
	// alias ka2a1='kubectl a2 a1'
	// alias ka1a3='kubectl a1 a3'
	// alias ka3a1='kubectl a3 a1'
	// alias ka2a3='kubectl a2 a3'
	// alias ka3a2='kubectl a3 a2'
	// alias ka1a2a3='kubectl a1 a2 a3'
	// alias ka1a3a2='kubectl a1 a3 a2'
	// alias ka2a1a3='kubectl a2 a1 a3'
	// alias ka2a3a1='kubectl a2 a3 a1'
	// alias ka3a1a2='kubectl a3 a1 a2'
	// alias ka3a2a1='kubectl a3 a2 a1'
}
