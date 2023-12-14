package main

// Scaffolding for testing the matching engine
var ruleSchemas = []RuleSchema{{
	class: inventoryItemClass,
	patternSchema: []AttrSchema{
		{name: "cat", valType: typeEnum},
		{name: "fullname", valType: typeStr},
		{name: "ageinstock", valType: typeInt},
		{name: "mrp", valType: typeFloat},
		{name: "received", valType: typeTS},
		{name: "bulkorder", valType: typeBool},
	},
}}

type RuleSchema struct {
	class         string
	patternSchema []AttrSchema
	actionSchema  ActionSchema
}

type AttrSchema struct {
	name    string
	valType string
	vals    map[string]bool
	valMin  float64
	valMax  float64
	lenMin  int
	lenMax  int
}

type ActionSchema struct {
	tasks      []string
	properties []string
}
