package main

// Scaffolding for testing the matching engine
var ruleSchemas = []RuleSchema{
	{
		"inventoryitem",
		[]AttrSchema{
			{"cat", "enum"},
			{"fullname", "str"},
			{"ageinstock", "int"},
			{"mrp", "float"},
			{"received", "ts"},
			{"bulkorder", "bool"},
		},
	},
}

type RuleSchema struct {
	class         string
	patternSchema []AttrSchema
}

type AttrSchema struct {
	name    string
	valType string
}
