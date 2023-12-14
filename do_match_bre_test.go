/*
This file contains the functions that are BRE tests for doMatch(). These functions are called
inside TestDoMatch() in do_match_test.go.

Some of the definitions of rulesets below deliberately use a lot of whitespace to keep the code consistent
and to make it easier to understand, add to, and edit these tests
*/

package main

import "testing"

const (
	// The "main" ruleset that may contain "thenCall"s/"elseCall"s to other rulesets
	mainRS = "main"

	inventoryItemClass = "inventoryitem"
	transactionClass   = "transaction"
	purchaseClass      = "purchase"
	orderClass         = "order"
)

var sampleEntity = Entity{inventoryItemClass, []Attr{
	{"cat", "textbook"},
	{"fullname", "Advanced Physics"},
	{"ageinstock", "5"},
	{"mrp", "50.80"},
	{"received", "2018-06-01T15:04:05Z"},
	{"bulkorder", trueStr},
}}

func testBasic(tests *[]doMatchTest) {
	ruleSet := RuleSet{1, inventoryItemClass, mainRS,
		[]Rule{{
			[]RulePatternTerm{{"cat", opEQ, "textbook"}},
			RuleActions{
				Tasks:      []string{"yearendsale", "summersale"},
				Properties: []Property{{"cashback", "10"}, {"discount", "9"}},
			},
		}},
	}
	*tests = append(*tests, doMatchTest{
		"basic test", sampleEntity, ruleSet, ActionSet{},
		ActionSet{
			tasks:      []string{"yearendsale", "summersale"},
			properties: []Property{{"cashback", "10"}, {"discount", "9"}},
		},
	})
}

func testExit(tests *[]doMatchTest) {
	rA1 := RuleActions{
		Tasks:      []string{"springsale"},
		Properties: []Property{{"cashback", "15"}},
	}
	rA2 := RuleActions{
		Tasks:      []string{"yearendsale", "summersale"},
		Properties: []Property{{"discount", "10"}, {"freegift", "mug"}},
	}
	rA3 := RuleActions{
		Tasks:      []string{"wintersale"},
		Properties: []Property{{"discount", "15"}},
		WillExit:   true,
	}
	rA4 := RuleActions{
		Tasks: []string{"autumnsale"},
	}
	ruleSet := RuleSet{1, inventoryItemClass, mainRS, []Rule{
		{[]RulePatternTerm{{"cat", opEQ, "refbook"}}, rA1},                           // no match
		{[]RulePatternTerm{{"ageinstock", opLT, 7}, {"cat", opEQ, "textbook"}}, rA2}, // match
		{[]RulePatternTerm{{"summersale", opEQ, true}}, rA3},                         // match then exit
		{[]RulePatternTerm{{"ageinstock", opLT, 7}}, rA4},                            // ignored
	}}
	want := ActionSet{
		tasks:      []string{"yearendsale", "summersale", "wintersale"},
		properties: []Property{{"discount", "15"}, {"freegift", "mug"}},
	}
	*tests = append(*tests, doMatchTest{"exit", sampleEntity, ruleSet, ActionSet{}, want})
}

func testReturn(tests *[]doMatchTest) {
	rA1 := RuleActions{
		Tasks:      []string{"yearendsale", "summersale"},
		Properties: []Property{{"discount", "10"}, {"freegift", "mug"}},
	}
	rA2 := RuleActions{
		Tasks:      []string{"springsale"},
		Properties: []Property{{"discount", "15"}},
		WillReturn: true,
	}
	rA3 := RuleActions{
		Tasks: []string{"autumnsale"},
	}
	ruleSet := RuleSet{1, inventoryItemClass, mainRS, []Rule{
		{[]RulePatternTerm{{"ageinstock", opLT, 7}, {"cat", opEQ, "textbook"}}, rA1}, // match
		{[]RulePatternTerm{{"summersale", opEQ, true}}, rA2},                         // match then return
		{[]RulePatternTerm{{"ageinstock", opLT, 7}}, rA3},                            // ignored
	}}
	want := ActionSet{
		tasks:      []string{"yearendsale", "summersale", "springsale"},
		properties: []Property{{"discount", "15"}, {"freegift", "mug"}},
	}
	*tests = append(*tests, doMatchTest{"return", sampleEntity, ruleSet, ActionSet{}, want})
}

func testTransactions(tests *[]doMatchTest) {
	ruleSchemas = append(ruleSchemas, RuleSchema{
		class: transactionClass,
		patternSchema: []AttrSchema{
			{name: "productname", valType: typeStr},
			{name: "price", valType: typeInt},
			{name: "inwintersale", valType: typeBool},
			{name: "paymenttype", valType: typeEnum},
			{name: "ismember", valType: typeBool},
		},
	})

	setupRuleSetMainForTransaction()
	setupRuleSetWinterDisc()
	setupRuleSetRegularDisc()
	setupRuleSetMemberDisc()
	setupRuleSetNonMemberDisc()

	// Each test below involves calling doMatch() with a different entity
	testWinterDiscJacket60(tests)
	testWinterDiscJacket40(tests)
	testWinterDiscKettle110Cash(tests)
	testWinterDiscKettle110Card(tests)
	testMemberDiscLamp60(tests)
	testMemberDiscKettle60Card(tests)
	testMemberDiscKettle60Cash(tests)
	testMemberDiscKettle110Card(tests)
	testMemberDiscKettle110Cash(tests)
	testNonMemberDiscLamp30(tests)
	testNonMemberDiscKettle70(tests)
	testNonMemberDiscKettle110Cash(tests)
	testNonMemberDiscKettle110Card(tests)
}

func setupRuleSetMainForTransaction() {
	rule1 := Rule{
		[]RulePatternTerm{
			{"inwintersale", opEQ, true},
		},
		RuleActions{
			ThenCall: "winterdisc",
			ElseCall: "regulardisc",
		},
	}
	rule2 := Rule{
		[]RulePatternTerm{
			{"paymenttype", opEQ, "cash"},
			{"price", opGT, 10},
		},
		RuleActions{
			Tasks: []string{"freepen"},
		},
	}
	rule3 := Rule{
		[]RulePatternTerm{
			{"paymenttype", opEQ, "card"},
			{"price", opGT, 10},
		},
		RuleActions{
			Tasks: []string{"freemug"},
		},
	}
	rule4 := Rule{
		[]RulePatternTerm{
			{"freehat", opEQ, true},
		},
		RuleActions{Tasks: []string{"freebag"}},
	}
	ruleSets[mainRS] = RuleSet{1, transactionClass, mainRS,
		[]Rule{rule1, rule2, rule3, rule4},
	}
}

func setupRuleSetWinterDisc() {
	rule1 := Rule{
		[]RulePatternTerm{
			{"productname", opEQ, "jacket"},
			{"price", opGT, 50},
		},
		RuleActions{
			Tasks:      []string{"freehat"},
			Properties: []Property{{"discount", "50"}},
			WillReturn: true,
		},
	}
	rule2 := Rule{
		[]RulePatternTerm{
			{"price", opLT, 100},
		},
		RuleActions{
			Properties: []Property{{"discount", "40"}, {"pointsmult", "2"}},
		},
	}
	rule3 := Rule{
		[]RulePatternTerm{
			{"price", opGE, 100},
		},
		RuleActions{
			Properties: []Property{{"discount", "45"}, {"pointsmult", "3"}},
		},
	}
	ruleSets["winterdisc"] = RuleSet{1, transactionClass, "winterdisc",
		[]Rule{rule1, rule2, rule3},
	}
}

func setupRuleSetRegularDisc() {
	rule1 := Rule{
		[]RulePatternTerm{
			{"ismember", opEQ, true},
		},
		RuleActions{
			ThenCall: "memberdisc",
			ElseCall: "nonmemberdisc",
		},
	}
	ruleSets["regulardisc"] = RuleSet{1, transactionClass, "regulardisc",
		[]Rule{rule1},
	}
}

func setupRuleSetMemberDisc() {
	rule1 := Rule{
		[]RulePatternTerm{
			{"productname", opEQ, "lamp"},
			{"price", opGT, 50},
		},
		RuleActions{
			Properties: []Property{{"discount", "35"}, {"pointsmult", "2"}},
			WillExit:   true,
		},
	}
	rule2 := Rule{
		[]RulePatternTerm{
			{"price", opLT, 100},
		},
		RuleActions{
			Properties: []Property{{"discount", "20"}},
		},
	}
	rule3 := Rule{
		[]RulePatternTerm{
			{"price", opGE, 100},
		},
		RuleActions{
			Properties: []Property{{"discount", "25"}},
		},
	}
	ruleSets["memberdisc"] = RuleSet{1, transactionClass, "memberdisc",
		[]Rule{rule1, rule2, rule3},
	}
}

func setupRuleSetNonMemberDisc() {
	rule1 := Rule{
		[]RulePatternTerm{
			{"price", opLT, 50},
		},
		RuleActions{
			Properties: []Property{{"discount", "5"}},
		},
	}
	rule2 := Rule{
		[]RulePatternTerm{
			{"price", opGE, 50},
		},
		RuleActions{
			Properties: []Property{{"discount", "10"}},
		},
	}
	rule3 := Rule{
		[]RulePatternTerm{
			{"price", opGE, 100},
		},
		RuleActions{
			Properties: []Property{{"discount", "15"}},
		},
	}
	ruleSets["nonmemberdisc"] = RuleSet{1, transactionClass, "nonmemberdisc",
		[]Rule{rule1, rule2, rule3},
	}
}

func testWinterDiscJacket60(tests *[]doMatchTest) {
	entity := Entity{transactionClass,
		[]Attr{
			{"productname", "jacket"},
			{"price", "60"},
			{"inwintersale", trueStr},
			{"paymenttype", "card"},
			{"ismember", trueStr},
		},
	}
	want := ActionSet{
		tasks:      []string{"freehat", "freemug", "freebag"},
		properties: []Property{{"discount", "50"}},
	}
	*tests = append(*tests, doMatchTest{
		"winterdisc jacket 60",
		entity,
		ruleSets[mainRS],
		ActionSet{},
		want,
	})
}

func testWinterDiscJacket40(tests *[]doMatchTest) {
	entity := Entity{transactionClass,
		[]Attr{
			{"productname", "jacket"},
			{"price", "40"},
			{"inwintersale", trueStr},
			{"paymenttype", "card"},
			{"ismember", trueStr},
		},
	}
	want := ActionSet{
		tasks:      []string{"freemug"},
		properties: []Property{{"discount", "40"}, {"pointsmult", "2"}},
	}
	*tests = append(*tests, doMatchTest{
		"winterdisc jacket 40",
		entity,
		ruleSets[mainRS],
		ActionSet{},
		want,
	})
}

func testWinterDiscKettle110Cash(tests *[]doMatchTest) {
	entity := Entity{transactionClass,
		[]Attr{
			{"productname", "kettle"},
			{"price", "110"},
			{"inwintersale", trueStr},
			{"paymenttype", "cash"},
			{"ismember", trueStr},
		},
	}
	want := ActionSet{
		tasks:      []string{"freepen"},
		properties: []Property{{"discount", "45"}, {"pointsmult", "3"}},
	}
	*tests = append(*tests, doMatchTest{
		"winterdisc kettle 110 cash",
		entity,
		ruleSets[mainRS],
		ActionSet{},
		want,
	})
}

func testWinterDiscKettle110Card(tests *[]doMatchTest) {
	entity := Entity{transactionClass,
		[]Attr{
			{"productname", "kettle"},
			{"price", "110"},
			{"inwintersale", trueStr},
			{"paymenttype", "card"},
			{"ismember", trueStr},
		},
	}
	want := ActionSet{
		tasks:      []string{"freemug"},
		properties: []Property{{"discount", "45"}, {"pointsmult", "3"}},
	}
	*tests = append(*tests, doMatchTest{
		"winterdisc kettle 110 card",
		entity,
		ruleSets[mainRS],
		ActionSet{},
		want,
	})
}

func testMemberDiscLamp60(tests *[]doMatchTest) {
	entity := Entity{transactionClass,
		[]Attr{
			{"productname", "lamp"},
			{"price", "60"},
			{"inwintersale", falseStr},
			{"paymenttype", "card"},
			{"ismember", trueStr},
		},
	}
	want := ActionSet{
		properties: []Property{{"discount", "35"}, {"pointsmult", "2"}},
	}
	*tests = append(*tests, doMatchTest{
		"memberdisc lamp 60",
		entity,
		ruleSets[mainRS],
		ActionSet{},
		want,
	})
}

func testMemberDiscKettle60Card(tests *[]doMatchTest) {
	entity := Entity{transactionClass,
		[]Attr{
			{"productname", "kettle"},
			{"price", "60"},
			{"inwintersale", falseStr},
			{"paymenttype", "card"},
			{"ismember", trueStr},
		},
	}
	want := ActionSet{
		tasks:      []string{"freemug"},
		properties: []Property{{"discount", "20"}},
	}
	*tests = append(*tests, doMatchTest{
		"memberdisc kettle 60 card",
		entity,
		ruleSets[mainRS],
		ActionSet{},
		want,
	})
}

func testMemberDiscKettle60Cash(tests *[]doMatchTest) {
	entity := Entity{transactionClass,
		[]Attr{
			{"productname", "kettle"},
			{"price", "60"},
			{"inwintersale", falseStr},
			{"paymenttype", "cash"},
			{"ismember", trueStr},
		},
	}
	want := ActionSet{
		tasks:      []string{"freepen"},
		properties: []Property{{"discount", "20"}},
	}
	*tests = append(*tests, doMatchTest{
		"memberdisc kettle 60 cash",
		entity,
		ruleSets[mainRS],
		ActionSet{},
		want,
	})
}

func testMemberDiscKettle110Card(tests *[]doMatchTest) {
	entity := Entity{transactionClass,
		[]Attr{
			{"productname", "kettle"},
			{"price", "110"},
			{"inwintersale", falseStr},
			{"paymenttype", "card"},
			{"ismember", trueStr},
		},
	}
	want := ActionSet{
		tasks:      []string{"freemug"},
		properties: []Property{{"discount", "25"}},
	}
	*tests = append(*tests, doMatchTest{
		"memberdisc kettle 110 card",
		entity,
		ruleSets[mainRS],
		ActionSet{},
		want,
	})
}

func testMemberDiscKettle110Cash(tests *[]doMatchTest) {
	entity := Entity{transactionClass,
		[]Attr{
			{"productname", "kettle"},
			{"price", "110"},
			{"inwintersale", falseStr},
			{"paymenttype", "cash"},
			{"ismember", trueStr},
		},
	}
	want := ActionSet{
		tasks:      []string{"freepen"},
		properties: []Property{{"discount", "25"}},
	}
	*tests = append(*tests, doMatchTest{
		"memberdisc kettle 110 cash",
		entity,
		ruleSets[mainRS],
		ActionSet{},
		want,
	})
}

func testNonMemberDiscLamp30(tests *[]doMatchTest) {
	entity := Entity{transactionClass,
		[]Attr{
			{"productname", "lamp"},
			{"price", "30"},
			{"inwintersale", falseStr},
			{"paymenttype", "cash"},
			{"ismember", falseStr},
		},
	}
	want := ActionSet{
		tasks:      []string{"freepen"},
		properties: []Property{{"discount", "5"}},
	}
	*tests = append(*tests, doMatchTest{
		"nonmemberdisc lamp 30",
		entity,
		ruleSets[mainRS],
		ActionSet{},
		want,
	})
}

func testNonMemberDiscKettle70(tests *[]doMatchTest) {
	entity := Entity{transactionClass,
		[]Attr{
			{"productname", "kettle"},
			{"price", "70"},
			{"inwintersale", falseStr},
			{"paymenttype", "cash"},
			{"ismember", falseStr},
		},
	}
	want := ActionSet{
		tasks:      []string{"freepen"},
		properties: []Property{{"discount", "10"}},
	}
	*tests = append(*tests, doMatchTest{
		"nonmemberdisc kettle 70",
		entity,
		ruleSets[mainRS],
		ActionSet{},
		want,
	})
}

func testNonMemberDiscKettle110Cash(tests *[]doMatchTest) {
	entity := Entity{transactionClass,
		[]Attr{
			{"productname", "kettle"},
			{"price", "110"},
			{"inwintersale", falseStr},
			{"paymenttype", "cash"},
			{"ismember", falseStr},
		},
	}
	want := ActionSet{
		tasks:      []string{"freepen"},
		properties: []Property{{"discount", "15"}},
	}
	*tests = append(*tests, doMatchTest{
		"nonmemberdisc kettle 110 cash",
		entity,
		ruleSets[mainRS],
		ActionSet{},
		want,
	})
}

func testNonMemberDiscKettle110Card(tests *[]doMatchTest) {
	entity := Entity{transactionClass,
		[]Attr{
			{"productname", "kettle"},
			{"price", "110"},
			{"inwintersale", falseStr},
			{"paymenttype", "card"},
			{"ismember", falseStr},
		},
	}
	want := ActionSet{
		tasks:      []string{"freemug"},
		properties: []Property{{"discount", "15"}},
	}
	*tests = append(*tests, doMatchTest{
		"nonmemberdisc kettle 110 card",
		entity,
		ruleSets[mainRS],
		ActionSet{},
		want,
	})
}

func testPurchases(tests *[]doMatchTest) {
	setupPurchaseRuleSchema()
	setupRuleSetForPurchases()

	// Each test below involves calling doMatch() with a different entity
	testJacket35(tests)
	testJacket55ForMember(tests)
	testJacket55ForNonMember(tests)
	testJacket75ForMember(tests)
	testJacket75ForNonMember(tests)
	testLamp35(tests)
	testLamp55(tests)
	testLamp75ForMember(tests)
	testLamp75ForNonMember(tests)
	testKettle35(tests)
	testKettle55(tests)
	testKettle75ForMember(tests)
	testKettle75ForNonMember(tests)
	testOven35(tests)
	testOven55(tests)
}

func setupPurchaseRuleSchema() {
	ruleSchemas = append(ruleSchemas, RuleSchema{
		class: purchaseClass,
		patternSchema: []AttrSchema{
			{name: "product", valType: typeStr},
			{name: "price", valType: typeFloat},
			{name: "ismember", valType: typeBool},
		},
		actionSchema: ActionSchema{
			tasks: []string{"freepen", "freebottle", "freepencil", "freemug", "freejar", "freeplant",
				"freebag", "freenotebook"},
			properties: []string{"discount", "pointsmult"},
		},
	})
}

func testJacket35(tests *[]doMatchTest) {
	entity := Entity{purchaseClass,
		[]Attr{
			{"product", "jacket"},
			{"price", "35"},
			{"ismember", trueStr},
		},
	}
	want := ActionSet{
		tasks:      []string{"freepen", "freebottle", "freepencil"},
		properties: []Property{{"discount", "5"}},
	}
	*tests = append(*tests, doMatchTest{
		"jacket price 35",
		entity,
		ruleSets[mainRS],
		ActionSet{},
		want,
	})
}

func testJacket55ForMember(tests *[]doMatchTest) {
	entity := Entity{purchaseClass,
		[]Attr{
			{"product", "jacket"},
			{"price", "55"},
			{"ismember", trueStr},
		},
	}
	want := ActionSet{
		tasks:      []string{"freepen", "freebottle", "freepencil", "freenotebook"},
		properties: []Property{{"discount", "10"}},
	}
	*tests = append(*tests, doMatchTest{
		"jacket price 55 for member",
		entity,
		ruleSets[mainRS],
		ActionSet{},
		want,
	})
}

func testJacket55ForNonMember(tests *[]doMatchTest) {
	entity := Entity{purchaseClass,
		[]Attr{
			{"product", "jacket"},
			{"price", "55"},
			{"ismember", falseStr},
		},
	}
	want := ActionSet{
		tasks:      []string{"freepen", "freebottle", "freepencil", "freenotebook"},
		properties: []Property{{"discount", "10"}},
	}
	*tests = append(*tests, doMatchTest{
		"jacket price 55 for non-member",
		entity,
		ruleSets[mainRS],
		ActionSet{},
		want,
	})
}

func testJacket75ForMember(tests *[]doMatchTest) {
	entity := Entity{purchaseClass,
		[]Attr{
			{"product", "jacket"},
			{"price", "75"},
			{"ismember", trueStr},
		},
	}
	want := ActionSet{
		tasks:      []string{"freepen", "freebottle", "freepencil", "freenotebook"},
		properties: []Property{{"discount", "15"}, {"pointsmult", "2"}},
	}
	*tests = append(*tests, doMatchTest{
		"jacket price 75 for member",
		entity,
		ruleSets[mainRS],
		ActionSet{},
		want,
	})
}

func testJacket75ForNonMember(tests *[]doMatchTest) {
	entity := Entity{purchaseClass,
		[]Attr{
			{"product", "jacket"},
			{"price", "75"},
			{"ismember", falseStr},
		},
	}
	want := ActionSet{
		tasks:      []string{"freepen", "freebottle", "freepencil", "freenotebook"},
		properties: []Property{{"discount", "10"}},
	}
	*tests = append(*tests, doMatchTest{
		"jacket price 75 for non-member",
		entity,
		ruleSets[mainRS],
		ActionSet{},
		want,
	})
}

func testLamp35(tests *[]doMatchTest) {
	entity := Entity{purchaseClass,
		[]Attr{
			{"product", "lamp"},
			{"price", "35"},
			{"ismember", falseStr},
		},
	}
	want := ActionSet{
		tasks:      []string{"freemug", "freejar", "freeplant", "freebag"},
		properties: []Property{{"discount", "20"}},
	}
	*tests = append(*tests, doMatchTest{
		"lamp price 35",
		entity,
		ruleSets[mainRS],
		ActionSet{},
		want,
	})
}

func testLamp55(tests *[]doMatchTest) {
	entity := Entity{purchaseClass,
		[]Attr{
			{"product", "lamp"},
			{"price", "55"},
			{"ismember", falseStr},
		},
	}
	want := ActionSet{
		tasks:      []string{"freemug", "freejar", "freeplant", "freebag", "freenotebook"},
		properties: []Property{{"discount", "25"}},
	}
	*tests = append(*tests, doMatchTest{
		"lamp price 55",
		entity,
		ruleSets[mainRS],
		ActionSet{},
		want,
	})
}

func testLamp75ForMember(tests *[]doMatchTest) {
	entity := Entity{purchaseClass,
		[]Attr{
			{"product", "lamp"},
			{"price", "75"},
			{"ismember", trueStr},
		},
	}
	want := ActionSet{
		tasks:      []string{"freemug", "freejar", "freeplant"},
		properties: []Property{{"discount", "30"}, {"pointsmult", "3"}},
	}
	*tests = append(*tests, doMatchTest{
		"lamp price 75 for member",
		entity,
		ruleSets[mainRS],
		ActionSet{},
		want,
	})
}

func testLamp75ForNonMember(tests *[]doMatchTest) {
	entity := Entity{purchaseClass,
		[]Attr{
			{"product", "lamp"},
			{"price", "75"},
			{"ismember", falseStr},
		},
	}
	want := ActionSet{
		tasks:      []string{"freemug", "freejar", "freeplant", "freebag", "freenotebook"},
		properties: []Property{{"discount", "25"}},
	}
	*tests = append(*tests, doMatchTest{
		"lamp price 75 for non-member",
		entity,
		ruleSets[mainRS],
		ActionSet{},
		want,
	})
}

func testKettle35(tests *[]doMatchTest) {
	entity := Entity{purchaseClass,
		[]Attr{
			{"product", "kettle"},
			{"price", "35"},
			{"ismember", trueStr},
		},
	}
	want := ActionSet{
		properties: []Property{{"discount", "35"}},
	}
	*tests = append(*tests, doMatchTest{
		"kettle price 35",
		entity,
		ruleSets[mainRS],
		ActionSet{},
		want,
	})
}

func testKettle55(tests *[]doMatchTest) {
	entity := Entity{purchaseClass,
		[]Attr{
			{"product", "kettle"},
			{"price", "55"},
			{"ismember", falseStr},
		},
	}
	want := ActionSet{
		tasks:      []string{"freenotebook"},
		properties: []Property{{"discount", "40"}},
	}
	*tests = append(*tests, doMatchTest{
		"kettle price 55",
		entity,
		ruleSets[mainRS],
		ActionSet{},
		want,
	})
}

func testKettle75ForMember(tests *[]doMatchTest) {
	entity := Entity{purchaseClass,
		[]Attr{
			{"product", "kettle"},
			{"price", "75"},
			{"ismember", trueStr},
		},
	}
	want := ActionSet{
		properties: []Property{{"discount", "45"}, {"pointsmult", "4"}},
	}
	*tests = append(*tests, doMatchTest{
		"kettle price 75 for member",
		entity,
		ruleSets[mainRS],
		ActionSet{},
		want,
	})
}

func testKettle75ForNonMember(tests *[]doMatchTest) {
	entity := Entity{purchaseClass,
		[]Attr{
			{"product", "kettle"},
			{"price", "75"},
			{"ismember", falseStr},
		},
	}
	want := ActionSet{
		tasks:      []string{"freenotebook"},
		properties: []Property{{"discount", "40"}},
	}
	*tests = append(*tests, doMatchTest{
		"kettle price 75 for non-member",
		entity,
		ruleSets[mainRS],
		ActionSet{},
		want,
	})
}

func testOven35(tests *[]doMatchTest) {
	entity := Entity{purchaseClass,
		[]Attr{
			{"product", "oven"},
			{"price", "35"},
			{"ismember", falseStr},
		},
	}
	want := ActionSet{}
	*tests = append(*tests, doMatchTest{
		"oven price 35",
		entity,
		ruleSets[mainRS],
		ActionSet{},
		want,
	})
}

func testOven55(tests *[]doMatchTest) {
	entity := Entity{purchaseClass,
		[]Attr{
			{"product", "oven"},
			{"price", "55"},
			{"ismember", trueStr},
		},
	}
	want := ActionSet{
		tasks: []string{"freenotebook"},
	}
	*tests = append(*tests, doMatchTest{
		"oven price 55",
		entity,
		ruleSets[mainRS],
		ActionSet{},
		want,
	})
}

func setupRuleSetForPurchases() {
	rule1 := Rule{
		[]RulePatternTerm{
			{"product", opEQ, "jacket"},
			{"price", opGT, 30.0},
		},
		RuleActions{
			Tasks:      []string{"freepen", "freebottle", "freepencil"},
			Properties: []Property{{"discount", "5"}},
		},
	}
	rule2 := Rule{
		[]RulePatternTerm{
			{"product", opEQ, "jacket"},
			{"price", opGT, 50.0},
		},
		RuleActions{
			Properties: []Property{{"discount", "10"}},
		},
	}
	rule3 := Rule{
		[]RulePatternTerm{
			{"product", opEQ, "jacket"},
			{"price", opGT, 70.0},
			{"ismember", opEQ, true},
		},
		RuleActions{
			Properties: []Property{{"discount", "15"}, {"pointsmult", "2"}},
		},
	}
	rule4 := Rule{
		[]RulePatternTerm{
			{"product", opEQ, "lamp"},
			{"price", opGT, 30.0},
		},
		RuleActions{
			Tasks:      []string{"freemug", "freejar", "freeplant"},
			Properties: []Property{{"discount", "20"}},
		},
	}
	rule5 := Rule{
		[]RulePatternTerm{
			{"product", opEQ, "lamp"},
			{"price", opGT, 50.0},
		},
		RuleActions{
			Properties: []Property{{"discount", "25"}},
		},
	}
	rule6 := Rule{
		[]RulePatternTerm{
			{"product", opEQ, "lamp"},
			{"price", opGT, 70.0},
			{"ismember", opEQ, true},
		},
		RuleActions{
			Properties: []Property{{"discount", "30"}, {"pointsmult", "3"}},
			WillExit:   true,
		},
	}
	rule7 := Rule{
		[]RulePatternTerm{
			{"product", opEQ, "kettle"},
			{"price", opGT, 30.0},
		},
		RuleActions{
			Properties: []Property{{"discount", "35"}},
		},
	}
	rule8 := Rule{
		[]RulePatternTerm{
			{"product", opEQ, "kettle"},
			{"price", opGT, 50.0},
		},
		RuleActions{
			Properties: []Property{{"discount", "40"}},
		},
	}
	rule9 := Rule{
		[]RulePatternTerm{
			{"product", opEQ, "kettle"},
			{"price", opGT, 70.0},
			{"ismember", opEQ, true},
		},
		RuleActions{
			Properties: []Property{{"discount", "45"}, {"pointsmult", "4"}},
			WillReturn: true,
		},
	}
	rule10 := Rule{
		[]RulePatternTerm{
			{"freemug", opEQ, true},
		},
		RuleActions{
			Tasks: []string{"freebag"},
		},
	}
	rule11 := Rule{
		[]RulePatternTerm{
			{"price", opGT, 50.0},
		},
		RuleActions{
			Tasks: []string{"freenotebook"},
		},
	}
	ruleSets[mainRS] = RuleSet{1, purchaseClass, mainRS,
		[]Rule{rule1, rule2, rule3, rule4, rule5, rule6, rule7, rule8, rule9, rule10, rule11},
	}
}

func testOrders(tests *[]doMatchTest) {
	ruleSchemas = append(ruleSchemas, RuleSchema{
		class: orderClass,
		patternSchema: []AttrSchema{
			{name: "ordertype", valType: typeEnum},
			{name: "mode", valType: typeEnum},
			{name: "liquidscheme", valType: typeBool},
			{name: "overnightscheme", valType: typeBool},
			{name: "extendedhours", valType: typeBool},
		},
		actionSchema: ActionSchema{
			tasks:      []string{"unitstoamc", "unitstorta"},
			properties: []string{"amfiordercutoff", "bseordercutoff", "fundscutoff", "unitscutoff"},
		},
	})

	setupRuleSetMainForOrder()
	setupRuleSetPurchaseOrSIPForOrder()
	setupRuleSetOtherOrderTypesForOrder()

	// Each test below involves calling doMatch() with a different entity
	testSIPOrder(tests)
	testSwitchDematOrder(tests)
	testSwitchDematExtHours(tests)
	testRedemptionDematExtHours(tests)
	testPurchaseOvernightOrder(tests)
	testSIPLiquidOrder(tests)
	testSwitchPhysicalOrder(tests)
}

func setupRuleSetMainForOrder() {
	rule1 := Rule{
		[]RulePatternTerm{
			{"ordertype", opEQ, "purchase"},
		},
		RuleActions{
			ThenCall: "purchaseorsip",
		},
	}
	rule2 := Rule{
		[]RulePatternTerm{
			{"ordertype", opEQ, "sip"},
		},
		RuleActions{
			ThenCall: "purchaseorsip",
		},
	}
	rule3 := Rule{
		[]RulePatternTerm{
			{"ordertype", opNE, "purchase"},
			{"ordertype", opNE, "sip"},
		},
		RuleActions{
			Properties: []Property{{"amfiordercutoff", "1500"}, {"bseordercutoff", "1500"}},
			ThenCall:   "otherordertypes",
		},
	}
	ruleSets[mainRS] = RuleSet{1, orderClass, mainRS,
		[]Rule{rule1, rule2, rule3},
	}
}

func setupRuleSetPurchaseOrSIPForOrder() {
	rule1 := Rule{
		[]RulePatternTerm{
			{"liquidscheme", opEQ, false},
			{"overnightscheme", opEQ, false},
		},
		RuleActions{
			Properties: []Property{{"amfiordercutoff", "1500"}, {"bseordercutoff", "1430"},
				{"fundscutoff", "1430"}},
			WillReturn: true,
		},
	}
	rule2 := Rule{
		[]RulePatternTerm{},
		RuleActions{
			Properties: []Property{{"amfiordercutoff", "1330"}, {"bseordercutoff", "1300"},
				{"fundscutoff", "1230"}},
		},
	}
	ruleSets["purchaseorsip"] = RuleSet{1, orderClass, "purchaseorsip",
		[]Rule{rule1, rule2},
	}
}

func setupRuleSetOtherOrderTypesForOrder() {
	rule1 := Rule{
		[]RulePatternTerm{
			{"mode", opEQ, "physical"},
		},
		RuleActions{
			Tasks: []string{"unitstoamc", "unitstorta"},
		},
	}
	rule2 := Rule{
		[]RulePatternTerm{
			{"mode", opEQ, "demat"},
			{"extendedhours", opEQ, false},
		},
		RuleActions{
			Properties: []Property{{"unitscutoff", "1630"}},
		},
	}
	rule3 := Rule{
		[]RulePatternTerm{
			{"mode", opEQ, "demat"},
			{"extendedhours", opEQ, true},
		},
		RuleActions{
			Properties: []Property{{"unitscutoff", "1730"}},
		},
	}
	ruleSets["otherordertypes"] = RuleSet{1, orderClass, "otherordertypes",
		[]Rule{rule1, rule2, rule3},
	}
}

func testSIPOrder(tests *[]doMatchTest) {
	entity := Entity{orderClass,
		[]Attr{
			{"ordertype", "sip"},
			{"mode", "demat"},
			{"liquidscheme", falseStr},
			{"overnightscheme", falseStr},
			{"extendedhours", falseStr},
		},
	}
	want := ActionSet{
		properties: []Property{{"amfiordercutoff", "1500"}, {"bseordercutoff", "1430"},
			{"fundscutoff", "1430"}},
	}
	*tests = append(*tests, doMatchTest{
		"sip order",
		entity,
		ruleSets[mainRS],
		ActionSet{},
		want,
	})
}

func testSwitchDematOrder(tests *[]doMatchTest) {
	entity := Entity{orderClass,
		[]Attr{
			{"ordertype", "switch"},
			{"mode", "demat"},
			{"liquidscheme", falseStr},
			{"overnightscheme", falseStr},
			{"extendedhours", falseStr},
		},
	}
	want := ActionSet{
		properties: []Property{{"amfiordercutoff", "1500"}, {"bseordercutoff", "1500"},
			{"unitscutoff", "1630"}},
	}
	*tests = append(*tests, doMatchTest{
		"switch demat order",
		entity,
		ruleSets[mainRS],
		ActionSet{},
		want,
	})
}

func testSwitchDematExtHours(tests *[]doMatchTest) {
	entity := Entity{orderClass,
		[]Attr{
			{"ordertype", "switch"},
			{"mode", "demat"},
			{"liquidscheme", falseStr},
			{"overnightscheme", falseStr},
			{"extendedhours", trueStr},
		},
	}
	want := ActionSet{
		properties: []Property{{"amfiordercutoff", "1500"}, {"bseordercutoff", "1500"},
			{"unitscutoff", "1730"}},
	}
	*tests = append(*tests, doMatchTest{
		"switch demat ext-hours order",
		entity,
		ruleSets[mainRS],
		ActionSet{},
		want,
	})
}

func testRedemptionDematExtHours(tests *[]doMatchTest) {
	entity := Entity{orderClass,
		[]Attr{
			{"ordertype", "redemption"},
			{"mode", "demat"},
			{"liquidscheme", falseStr},
			{"overnightscheme", falseStr},
			{"extendedhours", trueStr},
		},
	}
	want := ActionSet{
		properties: []Property{{"amfiordercutoff", "1500"}, {"bseordercutoff", "1500"},
			{"unitscutoff", "1730"}},
	}
	*tests = append(*tests, doMatchTest{
		"redemption demat ext-hours order",
		entity,
		ruleSets[mainRS],
		ActionSet{},
		want,
	})
}

func testPurchaseOvernightOrder(tests *[]doMatchTest) {
	entity := Entity{orderClass,
		[]Attr{
			{"ordertype", "purchase"},
			{"mode", "physical"},
			{"liquidscheme", falseStr},
			{"overnightscheme", trueStr},
			{"extendedhours", falseStr},
		},
	}
	want := ActionSet{
		properties: []Property{{"amfiordercutoff", "1330"}, {"bseordercutoff", "1300"},
			{"fundscutoff", "1230"}},
	}
	*tests = append(*tests, doMatchTest{
		"purchase overnight order",
		entity,
		ruleSets[mainRS],
		ActionSet{},
		want,
	})
}

func testSIPLiquidOrder(tests *[]doMatchTest) {
	entity := Entity{orderClass,
		[]Attr{
			{"ordertype", "sip"},
			{"mode", "physical"},
			{"liquidscheme", trueStr},
			{"overnightscheme", falseStr},
			{"extendedhours", falseStr},
		},
	}
	want := ActionSet{
		properties: []Property{{"amfiordercutoff", "1330"}, {"bseordercutoff", "1300"},
			{"fundscutoff", "1230"}},
	}
	*tests = append(*tests, doMatchTest{
		"sip liquid order",
		entity,
		ruleSets[mainRS],
		ActionSet{},
		want,
	})
}

func testSwitchPhysicalOrder(tests *[]doMatchTest) {
	entity := Entity{orderClass,
		[]Attr{
			{"ordertype", "switch"},
			{"mode", "physical"},
			{"liquidscheme", falseStr},
			{"overnightscheme", trueStr},
			{"extendedhours", trueStr},
		},
	}
	want := ActionSet{
		tasks:      []string{"unitstoamc", "unitstorta"},
		properties: []Property{{"amfiordercutoff", "1500"}, {"bseordercutoff", "1500"}},
	}
	*tests = append(*tests, doMatchTest{
		"switch physical order",
		entity,
		ruleSets[mainRS],
		ActionSet{},
		want,
	})
}

func testCycleError(t *testing.T) {
	t.Log("Running cycle test")
	setupRuleSetsForCycleError()
	_, _, err := doMatch(sampleEntity, ruleSets[mainRS], ActionSet{}, map[string]bool{})
	if err == nil {
		t.Errorf("test cycle: expected but did not get error")
	}
}

func setupRuleSetsForCycleError() {
	// main ruleset that contains a ThenCall to ruleset "second"
	rule1 := Rule{
		[]RulePatternTerm{
			{"cat", opEQ, "textbook"},
		},
		RuleActions{
			ThenCall: "second",
		},
	}
	ruleSets[mainRS] = RuleSet{1, inventoryItemClass, mainRS,
		[]Rule{rule1},
	}

	// "second" ruleset that contains a ThenCall to ruleset "third"
	rule1 = Rule{
		[]RulePatternTerm{
			{"cat", opEQ, "textbook"},
		},
		RuleActions{
			ThenCall: "third",
		},
	}
	ruleSets["second"] = RuleSet{1, inventoryItemClass, "second",
		[]Rule{rule1},
	}

	// "third" ruleset that contains a ThenCall back to ruleset "second"
	rule1 = Rule{
		[]RulePatternTerm{
			{"cat", opEQ, "textbook"},
		},
		RuleActions{
			Tasks: []string{"testtask"},
		},
	}
	rule2 := Rule{
		[]RulePatternTerm{
			{"cat", opEQ, "textbook"},
		},
		RuleActions{
			ThenCall: "second",
		},
	}
	ruleSets["third"] = RuleSet{1, inventoryItemClass, "third",
		[]Rule{rule1, rule2},
	}
}
