package main

// Scaffolding for testing the matching engine
var ruleSchemas = []RuleSchema{
	{
		"inventoryitem",
		[]AttrSchema{
			{"cat", typeEnum},
			{"fullname", typeStr},
			{"ageinstock", typeInt},
			{"mrp", typeFloat},
			{"received", typeTS},
			{"bulkorder", typeBool},
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
