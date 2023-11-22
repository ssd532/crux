/*
This file contains BRE tests for doMatch()

Some of the definitions of rulesets below deliberately use a lot of whitespace to keep the code consistent
and to make it easier to understand, add to, and edit these tests
*/

package main

import "testing"

const (
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
	ruleSet := RuleSet{1, inventoryItemClass, "main",
		[]Rule{{
			[]RulePatternTerm{{"cat", opEQ, "textbook"}},
			RuleActions{
				tasks:      []string{"yearendsale", "summersale"},
				properties: []Property{{"cashback", "10"}, {"discount", "9"}},
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
		tasks:      []string{"springsale"},
		properties: []Property{{"cashback", "15"}},
	}
	rA2 := RuleActions{
		tasks:      []string{"yearendsale", "summersale"},
		properties: []Property{{"discount", "10"}, {"freegift", "mug"}},
	}
	rA3 := RuleActions{
		tasks:      []string{"wintersale"},
		properties: []Property{{"discount", "15"}},
		willExit:   true,
	}
	rA4 := RuleActions{
		tasks: []string{"autumnsale"},
	}
	ruleSet := RuleSet{1, inventoryItemClass, "main", []Rule{
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
		tasks:      []string{"yearendsale", "summersale"},
		properties: []Property{{"discount", "10"}, {"freegift", "mug"}},
	}
	rA2 := RuleActions{
		tasks:      []string{"springsale"},
		properties: []Property{{"discount", "15"}},
		willReturn: true,
	}
	rA3 := RuleActions{
		tasks: []string{"autumnsale"},
	}
	ruleSet := RuleSet{1, inventoryItemClass, "main", []Rule{
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
		transactionClass,
		[]AttrSchema{
			{"productname", typeStr},
			{"price", typeInt},
			{"inwintersale", typeBool},
			{"paymenttype", typeEnum},
			{"ismember", typeBool},
		},
	})

	setupRuleSetMainForTransaction()
	setupRuleSetWinterDisc()
	setupRuleSetRegularDisc()
	setupRuleSetMemberDisc()
	setupRuleSetNonMemberDisc()

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
			thenCall: "winterdisc",
			elseCall: "regulardisc",
		},
	}
	rule2 := Rule{
		[]RulePatternTerm{
			{"paymenttype", opEQ, "cash"},
			{"price", opGT, 10},
		},
		RuleActions{
			tasks: []string{"freepen"},
		},
	}
	rule3 := Rule{
		[]RulePatternTerm{
			{"paymenttype", opEQ, "card"},
			{"price", opGT, 10},
		},
		RuleActions{
			tasks: []string{"freemug"},
		},
	}
	rule4 := Rule{
		[]RulePatternTerm{
			{"freehat", opEQ, true},
		},
		RuleActions{tasks: []string{"freebag"}},
	}
	ruleSets["main"] = RuleSet{1, transactionClass, "main",
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
			tasks:      []string{"freehat"},
			properties: []Property{{"discount", "50"}},
			willReturn: true,
		},
	}
	rule2 := Rule{
		[]RulePatternTerm{
			{"price", opLT, 100},
		},
		RuleActions{
			properties: []Property{{"discount", "40"}, {"pointsmult", "2"}},
		},
	}
	rule3 := Rule{
		[]RulePatternTerm{
			{"price", opGE, 100},
		},
		RuleActions{
			properties: []Property{{"discount", "45"}, {"pointsmult", "3"}},
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
			thenCall: "memberdisc",
			elseCall: "nonmemberdisc",
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
			properties: []Property{{"discount", "35"}, {"pointsmult", "2"}},
			willExit:   true,
		},
	}
	rule2 := Rule{
		[]RulePatternTerm{
			{"price", opLT, 100},
		},
		RuleActions{
			properties: []Property{{"discount", "20"}},
		},
	}
	rule3 := Rule{
		[]RulePatternTerm{
			{"price", opGE, 100},
		},
		RuleActions{
			properties: []Property{{"discount", "25"}},
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
			properties: []Property{{"discount", "5"}},
		},
	}
	rule2 := Rule{
		[]RulePatternTerm{
			{"price", opGE, 50},
		},
		RuleActions{
			properties: []Property{{"discount", "10"}},
		},
	}
	rule3 := Rule{
		[]RulePatternTerm{
			{"price", opGE, 100},
		},
		RuleActions{
			properties: []Property{{"discount", "15"}},
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
		ruleSets["main"],
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
		ruleSets["main"],
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
		ruleSets["main"],
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
		ruleSets["main"],
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
		ruleSets["main"],
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
		ruleSets["main"],
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
		ruleSets["main"],
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
		ruleSets["main"],
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
		ruleSets["main"],
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
		ruleSets["main"],
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
		ruleSets["main"],
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
		ruleSets["main"],
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
		ruleSets["main"],
		ActionSet{},
		want,
	})
}

func testPurchases(tests *[]doMatchTest) {
	ruleSchemas = append(ruleSchemas, RuleSchema{purchaseClass,
		[]AttrSchema{
			{"product", typeStr},
			{"price", typeFloat},
			{"ismember", typeBool},
		},
	})

	setupRuleSetForPurchases()

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
		ruleSets["main"],
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
		ruleSets["main"],
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
		ruleSets["main"],
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
		ruleSets["main"],
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
		ruleSets["main"],
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
		ruleSets["main"],
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
		ruleSets["main"],
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
		ruleSets["main"],
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
		ruleSets["main"],
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
		ruleSets["main"],
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
		ruleSets["main"],
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
		ruleSets["main"],
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
		ruleSets["main"],
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
		ruleSets["main"],
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
		ruleSets["main"],
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
			tasks:      []string{"freepen", "freebottle", "freepencil"},
			properties: []Property{{"discount", "5"}},
		},
	}
	rule2 := Rule{
		[]RulePatternTerm{
			{"product", opEQ, "jacket"},
			{"price", opGT, 50.0},
		},
		RuleActions{
			properties: []Property{{"discount", "10"}},
		},
	}
	rule3 := Rule{
		[]RulePatternTerm{
			{"product", opEQ, "jacket"},
			{"price", opGT, 70.0},
			{"ismember", opEQ, true},
		},
		RuleActions{
			properties: []Property{{"discount", "15"}, {"pointsmult", "2"}},
		},
	}
	rule4 := Rule{
		[]RulePatternTerm{
			{"product", opEQ, "lamp"},
			{"price", opGT, 30.0},
		},
		RuleActions{
			tasks:      []string{"freemug", "freejar", "freeplant"},
			properties: []Property{{"discount", "20"}},
		},
	}
	rule5 := Rule{
		[]RulePatternTerm{
			{"product", opEQ, "lamp"},
			{"price", opGT, 50.0},
		},
		RuleActions{
			properties: []Property{{"discount", "25"}},
		},
	}
	rule6 := Rule{
		[]RulePatternTerm{
			{"product", opEQ, "lamp"},
			{"price", opGT, 70.0},
			{"ismember", opEQ, true},
		},
		RuleActions{
			properties: []Property{{"discount", "30"}, {"pointsmult", "3"}},
			willExit:   true,
		},
	}
	rule7 := Rule{
		[]RulePatternTerm{
			{"product", opEQ, "kettle"},
			{"price", opGT, 30.0},
		},
		RuleActions{
			properties: []Property{{"discount", "35"}},
		},
	}
	rule8 := Rule{
		[]RulePatternTerm{
			{"product", opEQ, "kettle"},
			{"price", opGT, 50.0},
		},
		RuleActions{
			properties: []Property{{"discount", "40"}},
		},
	}
	rule9 := Rule{
		[]RulePatternTerm{
			{"product", opEQ, "kettle"},
			{"price", opGT, 70.0},
			{"ismember", opEQ, true},
		},
		RuleActions{
			properties: []Property{{"discount", "45"}, {"pointsmult", "4"}},
			willReturn: true,
		},
	}
	rule10 := Rule{
		[]RulePatternTerm{
			{"freemug", opEQ, true},
		},
		RuleActions{
			tasks: []string{"freebag"},
		},
	}
	rule11 := Rule{
		[]RulePatternTerm{
			{"price", opGT, 50.0},
		},
		RuleActions{
			tasks: []string{"freenotebook"},
		},
	}
	ruleSets["main"] = RuleSet{1, purchaseClass, "main",
		[]Rule{rule1, rule2, rule3, rule4, rule5, rule6, rule7, rule8, rule9, rule10, rule11},
	}
}

func testOrders(tests *[]doMatchTest) {
	ruleSchemas = append(ruleSchemas, RuleSchema{
		orderClass,
		[]AttrSchema{
			{"ordertype", typeEnum},
			{"mode", typeEnum},
			{"liquidscheme", typeBool},
			{"overnightscheme", typeBool},
			{"extendedhours", typeBool},
		},
	})

	setupRuleSetMainForOrder()
	setupRuleSetPurchaseOrSIPForOrder()
	setupRuleSetOtherOrderTypesForOrder()

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
			thenCall: "purchaseorsip",
		},
	}
	rule2 := Rule{
		[]RulePatternTerm{
			{"ordertype", opEQ, "sip"},
		},
		RuleActions{
			thenCall: "purchaseorsip",
		},
	}
	rule3 := Rule{
		[]RulePatternTerm{
			{"ordertype", opNE, "purchase"},
			{"ordertype", opNE, "sip"},
		},
		RuleActions{
			properties: []Property{{"amfiordercutoff", "1500"}, {"bseordercutoff", "1500"}},
			thenCall:   "otherordertypes",
		},
	}
	ruleSets["main"] = RuleSet{1, orderClass, "main",
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
			properties: []Property{{"amfiordercutoff", "1500"}, {"bseordercutoff", "1430"},
				{"fundscutoff", "1430"}},
			willReturn: true,
		},
	}
	rule2 := Rule{
		[]RulePatternTerm{},
		RuleActions{
			properties: []Property{{"amfiordercutoff", "1330"}, {"bseordercutoff", "1300"},
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
			tasks: []string{"unitstoamc", "unitstorta"},
		},
	}
	rule2 := Rule{
		[]RulePatternTerm{
			{"mode", opEQ, "demat"},
			{"extendedhours", opEQ, false},
		},
		RuleActions{
			properties: []Property{{"unitscutoff", "1630"}},
		},
	}
	rule3 := Rule{
		[]RulePatternTerm{
			{"mode", opEQ, "demat"},
			{"extendedhours", opEQ, true},
		},
		RuleActions{
			properties: []Property{{"unitscutoff", "1730"}},
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
		ruleSets["main"],
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
		ruleSets["main"],
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
		ruleSets["main"],
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
		ruleSets["main"],
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
		ruleSets["main"],
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
		ruleSets["main"],
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
		ruleSets["main"],
		ActionSet{},
		want,
	})
}

func testCycleError(t *testing.T) {
	t.Log("Running cycle test")
	setupRuleSetsForCycleError()
	_, _, err := doMatch(sampleEntity, ruleSets["main"], ActionSet{}, map[string]bool{})
	if err == nil {
		t.Errorf("test cycle: expected but did not get error")
	}
}

func setupRuleSetsForCycleError() {
	// main ruleset
	rule1 := Rule{
		[]RulePatternTerm{
			{"cat", opEQ, "textbook"},
		},
		RuleActions{
			thenCall: "second",
		},
	}
	ruleSets["main"] = RuleSet{1, inventoryItemClass, "main",
		[]Rule{rule1},
	}

	// second ruleset
	rule1 = Rule{
		[]RulePatternTerm{
			{"cat", opEQ, "textbook"},
		},
		RuleActions{
			thenCall: "third",
		},
	}
	ruleSets["second"] = RuleSet{1, inventoryItemClass, "second",
		[]Rule{rule1},
	}

	// third ruleset
	rule1 = Rule{
		[]RulePatternTerm{
			{"cat", opEQ, "textbook"},
		},
		RuleActions{
			tasks: []string{"testtask"},
		},
	}
	rule2 := Rule{
		[]RulePatternTerm{
			{"cat", opEQ, "textbook"},
		},
		RuleActions{
			thenCall: "second",
		},
	}
	ruleSets["third"] = RuleSet{1, inventoryItemClass, "third",
		[]Rule{rule1, rule2},
	}
}
