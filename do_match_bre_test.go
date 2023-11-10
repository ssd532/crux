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
	{"bulkorder", "true"},
}}

func testBasic(tests *[]doMatchTest) {
	ruleSet := RuleSet{1, inventoryItemClass, "main",
		[]Rule{{
			[]RulePatternTerm{{"cat", "eq", "textbook"}},
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
		{[]RulePatternTerm{{"cat", "eq", "refbook"}}, rA1},                           // no match
		{[]RulePatternTerm{{"ageinstock", "lt", 7}, {"cat", "eq", "textbook"}}, rA2}, // match
		{[]RulePatternTerm{{"summersale", "eq", true}}, rA3},                         // match then exit
		{[]RulePatternTerm{{"ageinstock", "lt", 7}}, rA4},                            // ignored
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
		{[]RulePatternTerm{{"ageinstock", "lt", 7}, {"cat", "eq", "textbook"}}, rA1}, // match
		{[]RulePatternTerm{{"summersale", "eq", true}}, rA2},                         // match then return
		{[]RulePatternTerm{{"ageinstock", "lt", 7}}, rA3},                            // ignored
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
			{"productname", "str"},
			{"price", "int"},
			{"inwintersale", "bool"},
			{"paymenttype", "enum"},
			{"ismember", "bool"},
		},
	})

	setupRuleSetMainForTransaction()
	setupRuleSetWinterDisc()
	setupRuleSetRegularDisc()
	setupRuleSetMemberDisc()
	setupRuleSetNonMemberDisc()

	testThenAndReturn(tests)
	testElseAndThenAndExit(tests)
	testElseAndElse(tests)
}

func setupRuleSetMainForTransaction() {
	rule1 := Rule{
		[]RulePatternTerm{
			{"inwintersale", "eq", true},
		},
		RuleActions{
			thenCall: "winterdisc",
			elseCall: "regulardisc",
		},
	}
	rule2 := Rule{
		[]RulePatternTerm{
			{"paymenttype", "eq", "cash"},
			{"price", "gt", 10},
		},
		RuleActions{
			tasks: []string{"freepen"},
		},
	}
	rule3 := Rule{
		[]RulePatternTerm{
			{"paymenttype", "eq", "card"},
			{"price", "gt", 10},
		},
		RuleActions{
			tasks: []string{"freemug"},
		},
	}
	rule4 := Rule{
		[]RulePatternTerm{
			{"freehat", "eq", true},
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
			{"productname", "eq", "jacket"},
			{"price", "gt", 50},
		},
		RuleActions{
			tasks:      []string{"freehat"},
			properties: []Property{{"discount", "50"}},
			willReturn: true,
		},
	}
	rule2 := Rule{
		[]RulePatternTerm{
			{"price", "lt", 100},
		},
		RuleActions{
			properties: []Property{{"discount", "40"}, {"pointsmult", "2"}},
		},
	}
	rule3 := Rule{
		[]RulePatternTerm{
			{"price", "ge", 100},
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
			{"ismember", "eq", true},
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
			{"productname", "eq", "lamp"},
			{"price", "gt", 50},
		},
		RuleActions{
			properties: []Property{{"discount", "35"}, {"pointsmult", "2"}},
			willExit:   true,
		},
	}
	rule2 := Rule{
		[]RulePatternTerm{
			{"price", "lt", 100},
		},
		RuleActions{
			properties: []Property{{"discount", "20"}},
		},
	}
	rule3 := Rule{
		[]RulePatternTerm{
			{"price", "ge", 100},
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
			{"price", "lt", 50},
		},
		RuleActions{
			properties: []Property{{"discount", "5"}},
		},
	}
	rule2 := Rule{
		[]RulePatternTerm{
			{"price", "lt", 100},
		},
		RuleActions{
			properties: []Property{{"discount", "10"}},
		},
	}
	rule3 := Rule{
		[]RulePatternTerm{
			{"price", "ge", 100},
		},
		RuleActions{
			properties: []Property{{"discount", "15"}},
		},
	}
	ruleSets["nonmemberdisc"] = RuleSet{1, transactionClass, "nonmemberdisc",
		[]Rule{rule1, rule2, rule3},
	}
}

func testThenAndReturn(tests *[]doMatchTest) {
	entity := Entity{transactionClass,
		[]Attr{
			{"productname", "jacket"},
			{"price", "60"},
			{"inwintersale", "true"},
			{"paymenttype", "card"},
			{"ismember", "true"},
		},
	}
	want := ActionSet{
		tasks:      []string{"freehat", "freemug", "freebag"},
		properties: []Property{{"discount", "50"}},
	}
	*tests = append(*tests, doMatchTest{
		"then and return",
		entity,
		ruleSets["main"],
		ActionSet{},
		want,
	})
}

func testElseAndThenAndExit(tests *[]doMatchTest) {
	entity := Entity{transactionClass,
		[]Attr{
			{"productname", "lamp"},
			{"price", "60"},
			{"inwintersale", "false"},
			{"paymenttype", "card"},
			{"ismember", "true"},
		},
	}
	want := ActionSet{
		properties: []Property{{"discount", "35"}, {"pointsmult", "2"}},
	}
	*tests = append(*tests, doMatchTest{
		"else and then and exit",
		entity,
		ruleSets["main"],
		ActionSet{},
		want,
	})
}

func testElseAndElse(tests *[]doMatchTest) {
	entity := Entity{transactionClass,
		[]Attr{
			{"productname", "umbrella"},
			{"price", "70"},
			{"inwintersale", "false"},
			{"paymenttype", "cash"},
			{"ismember", "false"},
		},
	}
	want := ActionSet{
		tasks:      []string{"freepen"},
		properties: []Property{{"discount", "10"}},
	}
	*tests = append(*tests, doMatchTest{
		"else and else",
		entity,
		ruleSets["main"],
		ActionSet{},
		want,
	})
}

func testPurchases(tests *[]doMatchTest) {
	ruleSchemas = append(ruleSchemas, RuleSchema{purchaseClass,
		[]AttrSchema{
			{"product", "str"},
			{"price", "float"},
			{"ismember", "bool"},
		},
	})

	setupRuleSetForPurchases()

	testJacket35(tests)
	testJacket55(tests)
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
			{"ismember", "true"},
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

func testJacket55(tests *[]doMatchTest) {
	entity := Entity{purchaseClass,
		[]Attr{
			{"product", "jacket"},
			{"price", "55"},
			{"ismember", "false"},
		},
	}
	want := ActionSet{
		tasks:      []string{"freepen", "freebottle", "freepencil", "freenotebook"},
		properties: []Property{{"discount", "10"}},
	}
	*tests = append(*tests, doMatchTest{
		"jacket price 55",
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
			{"ismember", "true"},
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
			{"ismember", "false"},
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
			{"ismember", "false"},
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
			{"ismember", "false"},
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
			{"ismember", "true"},
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
			{"ismember", "false"},
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
			{"ismember", "true"},
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
			{"ismember", "false"},
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
			{"ismember", "true"},
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
			{"ismember", "false"},
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
			{"ismember", "false"},
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
			{"ismember", "false"},
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
			{"product", "eq", "jacket"},
			{"price", "gt", 30.0},
		},
		RuleActions{
			tasks:      []string{"freepen", "freebottle", "freepencil"},
			properties: []Property{{"discount", "5"}},
		},
	}
	rule2 := Rule{
		[]RulePatternTerm{
			{"product", "eq", "jacket"},
			{"price", "gt", 50.0},
		},
		RuleActions{
			properties: []Property{{"discount", "10"}},
		},
	}
	rule3 := Rule{
		[]RulePatternTerm{
			{"product", "eq", "jacket"},
			{"price", "gt", 70.0},
			{"ismember", "eq", true},
		},
		RuleActions{
			properties: []Property{{"discount", "15"}, {"pointsmult", "2"}},
		},
	}
	rule4 := Rule{
		[]RulePatternTerm{
			{"product", "eq", "lamp"},
			{"price", "gt", 30.0},
		},
		RuleActions{
			tasks:      []string{"freemug", "freejar", "freeplant"},
			properties: []Property{{"discount", "20"}},
		},
	}
	rule5 := Rule{
		[]RulePatternTerm{
			{"product", "eq", "lamp"},
			{"price", "gt", 50.0},
		},
		RuleActions{
			properties: []Property{{"discount", "25"}},
		},
	}
	rule6 := Rule{
		[]RulePatternTerm{
			{"product", "eq", "lamp"},
			{"price", "gt", 70.0},
			{"ismember", "eq", true},
		},
		RuleActions{
			properties: []Property{{"discount", "30"}, {"pointsmult", "3"}},
			willExit:   true,
		},
	}
	rule7 := Rule{
		[]RulePatternTerm{
			{"product", "eq", "kettle"},
			{"price", "gt", 30.0},
		},
		RuleActions{
			properties: []Property{{"discount", "35"}},
		},
	}
	rule8 := Rule{
		[]RulePatternTerm{
			{"product", "eq", "kettle"},
			{"price", "gt", 50.0},
		},
		RuleActions{
			properties: []Property{{"discount", "40"}},
		},
	}
	rule9 := Rule{
		[]RulePatternTerm{
			{"product", "eq", "kettle"},
			{"price", "gt", 70.0},
			{"ismember", "eq", true},
		},
		RuleActions{
			properties: []Property{{"discount", "45"}, {"pointsmult", "4"}},
			willReturn: true,
		},
	}
	rule10 := Rule{
		[]RulePatternTerm{
			{"freemug", "eq", true},
		},
		RuleActions{
			tasks: []string{"freebag"},
		},
	}
	rule11 := Rule{
		[]RulePatternTerm{
			{"price", "gt", 50.0},
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
			{"ordertype", "enum"},
			{"mode", "enum"},
			{"liquidscheme", "bool"},
			{"overnightscheme", "bool"},
			{"extendedhours", "bool"},
		},
	})

	setupRuleSetMainForOrder()
	setupRuleSetPurchaseOrSIPForOrder()
	setupRuleSetOtherOrderTypesForOrder()

	testSIPOrder(tests)
	testSwitchDematOrder(tests)
	testPurchaseOvernightOrder(tests)
	testSwitchPhysicalOrder(tests)
}

func setupRuleSetMainForOrder() {
	rule1 := Rule{
		[]RulePatternTerm{
			{"ordertype", "eq", "purchase"},
		},
		RuleActions{
			thenCall: "purchaseorsip",
		},
	}
	rule2 := Rule{
		[]RulePatternTerm{
			{"ordertype", "eq", "sip"},
		},
		RuleActions{
			thenCall: "purchaseorsip",
		},
	}
	rule3 := Rule{
		[]RulePatternTerm{
			{"ordertype", "ne", "purchase"},
			{"ordertype", "ne", "sip"},
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
			{"liquidscheme", "eq", false},
			{"overnightscheme", "eq", false},
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
			{"mode", "eq", "physical"},
		},
		RuleActions{
			tasks: []string{"unitstoamc", "unitstorta"},
		},
	}
	rule2 := Rule{
		[]RulePatternTerm{
			{"mode", "eq", "demat"},
			{"extendedhours", "eq", false},
		},
		RuleActions{
			properties: []Property{{"unitscutoff", "1630"}},
		},
	}
	rule3 := Rule{
		[]RulePatternTerm{
			{"mode", "eq", "demat"},
			{"extendedhours", "eq", true},
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
			{"liquidscheme", "false"},
			{"overnightscheme", "false"},
			{"extendedhours", "false"},
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
			{"liquidscheme", "false"},
			{"overnightscheme", "false"},
			{"extendedhours", "false"},
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

func testPurchaseOvernightOrder(tests *[]doMatchTest) {
	entity := Entity{orderClass,
		[]Attr{
			{"ordertype", "purchase"},
			{"mode", "physical"},
			{"liquidscheme", "false"},
			{"overnightscheme", "true"},
			{"extendedhours", "false"},
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

func testSwitchPhysicalOrder(tests *[]doMatchTest) {
	entity := Entity{orderClass,
		[]Attr{
			{"ordertype", "switch"},
			{"mode", "physical"},
			{"liquidscheme", "false"},
			{"overnightscheme", "true"},
			{"extendedhours", "true"},
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
			{"cat", "eq", "textbook"},
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
			{"cat", "eq", "textbook"},
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
			{"cat", "eq", "textbook"},
		},
		RuleActions{
			tasks: []string{"testtask"},
		},
	}
	rule2 := Rule{
		[]RulePatternTerm{
			{"cat", "eq", "textbook"},
		},
		RuleActions{
			thenCall: "second",
		},
	}
	ruleSets["third"] = RuleSet{1, inventoryItemClass, "third",
		[]Rule{rule1, rule2},
	}
}
