/*
This file contains the data structures used by the matching engine
*/

package main

type Entity struct {
	class string
	attrs []Attr
}

type Attr struct {
	name string
	val  string
}

type ActionSet struct {
	tasks      []string
	properties []Property
}

type Property struct {
	Name string
	Val  string
}

type RuleSet struct {
	Ver     int
	Class   string
	SetName string
	Rules   []Rule
}

type Rule struct {
	RulePattern []RulePatternTerm
	RuleActions RuleActions
}

type RulePatternTerm struct {
	AttrName string
	Op       string
	AttrVal  any
}

type RuleActions struct {
	Tasks      []string
	Properties []Property
	ThenCall   string
	ElseCall   string
	WillReturn bool
	WillExit   bool
}
