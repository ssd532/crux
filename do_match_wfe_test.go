/*
Some of the definitions of rulesets below deliberately use a lot of whitespace to keep the code consistent
and to make it easier to understand, add to, and edit these tests
*/

package main

const (
	step       = "step"
	stepFailed = "stepfailed"
	nextStep   = "nextstep"
	start      = "START"
	endFlow    = "ENDFLOW"

	trueStr  = "true"
	falseStr = "false"

	uccCreationClass = "ucccreation"
	complexWFClass   = "complexwf"
)

func testUCCCreation(tests *[]doMatchTest) {
	ruleSchemas = append(ruleSchemas, RuleSchema{
		uccCreationClass,
		[]AttrSchema{
			{step, "str"},
			{stepFailed, "bool"},
			{"mode", "enum"},
		},
	})
	setupUCCCreationRuleSet()

	testUCCStart(tests)
	testUCCGetCustDetailsDemat(tests)
	testUCCGetCustDetailsDematFail(tests)
	testUCCGetCustDetailsPhysical(tests)
	testUCCGetCustDetailsPhysicalFail(tests)
	testUCCReadyForAuthLink(tests)
	testUCCReadyForAuthLinkFail(tests)
	testUCCEndSuccess(tests)
	testUCCEndFailure(tests)
}

func setupUCCCreationRuleSet() {
	rule1 := Rule{
		[]RulePatternTerm{
			{step, "eq", start},
			{stepFailed, "eq", false},
		},
		RuleActions{
			tasks:      []string{"getcustdetails"},
			properties: []Property{{nextStep, "getcustdetails"}},
		},
	}
	rule2 := Rule{
		[]RulePatternTerm{
			{step, "eq", "getcustdetails"},
			{stepFailed, "eq", false},
		},
		RuleActions{
			tasks:      []string{"aof", "kycvalid", "nomauth"},
			properties: []Property{{nextStep, "readyforauthlink"}},
		},
	}
	rule3 := Rule{
		[]RulePatternTerm{
			{step, "eq", "getcustdetails"},
			{stepFailed, "eq", false},
			{"mode", "eq", "physical"},
		},
		RuleActions{
			tasks:      []string{"bankaccvalid"},
			properties: []Property{{nextStep, "readyforauthlink"}},
		},
	}
	rule4 := Rule{
		[]RulePatternTerm{
			{step, "eq", "getcustdetails"},
			{stepFailed, "eq", false},
			{"mode", "eq", "demat"},
		},
		RuleActions{
			tasks:      []string{"dpandbankaccvalid"},
			properties: []Property{{nextStep, "readyforauthlink"}},
		},
	}
	rule5 := Rule{
		[]RulePatternTerm{
			{step, "eq", "readyforauthlink"},
			{stepFailed, "eq", false},
		},
		RuleActions{
			tasks:      []string{"sendauthlinktoclient"},
			properties: []Property{{nextStep, "sendauthlinktoclient"}},
		},
	}
	rule6 := Rule{
		[]RulePatternTerm{
			{step, "eq", "sendauthlinktoclient"},
			{stepFailed, "eq", false},
		},
		RuleActions{
			properties: []Property{{endFlow, trueStr}},
		},
	}
	rule7 := Rule{
		[]RulePatternTerm{
			{stepFailed, "eq", true},
		},
		RuleActions{
			properties: []Property{{endFlow, trueStr}},
		},
	}
	ruleSets["ucccreation"] = RuleSet{1, uccCreationClass, "ucccreation",
		[]Rule{rule1, rule2, rule3, rule4, rule5, rule6, rule7},
	}
}

func testUCCStart(tests *[]doMatchTest) {
	entity := Entity{uccCreationClass, []Attr{
		{step, start},
		{stepFailed, falseStr},
		{"mode", "demat"},
	}}
	want := ActionSet{
		tasks:      []string{"getcustdetails"},
		properties: []Property{{nextStep, "getcustdetails"}},
	}
	*tests = append(*tests, doMatchTest{"ucc start", entity, ruleSets["ucccreation"], ActionSet{}, want})
}

func testUCCGetCustDetailsDemat(tests *[]doMatchTest) {
	entity := Entity{uccCreationClass, []Attr{
		{step, "getcustdetails"},
		{stepFailed, falseStr},
		{"mode", "demat"},
	}}
	want := ActionSet{
		tasks:      []string{"aof", "kycvalid", "nomauth", "dpandbankaccvalid"},
		properties: []Property{{nextStep, "readyforauthlink"}},
	}
	*tests = append(*tests, doMatchTest{"ucc getcustdetails demat", entity, ruleSets["ucccreation"], ActionSet{}, want})
}

func testUCCGetCustDetailsDematFail(tests *[]doMatchTest) {
	entity := Entity{uccCreationClass, []Attr{
		{step, "getcustdetails"},
		{stepFailed, trueStr},
		{"mode", "demat"},
	}}
	want := ActionSet{
		properties: []Property{{endFlow, trueStr}},
	}
	*tests = append(*tests, doMatchTest{"ucc getcustdetails demat fail", entity, ruleSets["ucccreation"], ActionSet{}, want})
}

func testUCCGetCustDetailsPhysical(tests *[]doMatchTest) {
	entity := Entity{uccCreationClass, []Attr{
		{step, "getcustdetails"},
		{stepFailed, falseStr},
		{"mode", "physical"},
	}}
	want := ActionSet{
		tasks:      []string{"aof", "kycvalid", "nomauth", "bankaccvalid"},
		properties: []Property{{nextStep, "readyforauthlink"}},
	}
	*tests = append(*tests, doMatchTest{"ucc getcustdetails physical", entity, ruleSets["ucccreation"], ActionSet{}, want})
}

func testUCCGetCustDetailsPhysicalFail(tests *[]doMatchTest) {
	entity := Entity{uccCreationClass, []Attr{
		{step, "getcustdetails"},
		{stepFailed, trueStr},
		{"mode", "physical"},
	}}
	want := ActionSet{
		properties: []Property{{endFlow, trueStr}},
	}
	*tests = append(*tests, doMatchTest{"ucc getcustdetails physical fail", entity, ruleSets["ucccreation"], ActionSet{}, want})
}

func testUCCReadyForAuthLink(tests *[]doMatchTest) {
	entity := Entity{uccCreationClass, []Attr{
		{step, "readyforauthlink"},
		{stepFailed, falseStr},
		{"mode", "demat"},
	}}
	want := ActionSet{
		tasks:      []string{"sendauthlinktoclient"},
		properties: []Property{{nextStep, "sendauthlinktoclient"}},
	}
	*tests = append(*tests, doMatchTest{"ucc readyforauthlink", entity, ruleSets["ucccreation"], ActionSet{}, want})
}

func testUCCReadyForAuthLinkFail(tests *[]doMatchTest) {
	entity := Entity{uccCreationClass, []Attr{
		{step, "readyforauthlink"},
		{stepFailed, trueStr},
		{"mode", "demat"},
	}}
	want := ActionSet{
		properties: []Property{{endFlow, trueStr}},
	}
	*tests = append(*tests, doMatchTest{"ucc readyforauthlink fail", entity, ruleSets["ucccreation"], ActionSet{}, want})
}

func testUCCEndSuccess(tests *[]doMatchTest) {
	entity := Entity{uccCreationClass, []Attr{
		{step, "sendauthlinktoclient"},
		{stepFailed, falseStr},
		{"mode", "demat"},
	}}
	want := ActionSet{
		properties: []Property{{endFlow, trueStr}},
	}
	*tests = append(*tests, doMatchTest{"ucc end-success", entity, ruleSets["ucccreation"], ActionSet{}, want})
}

func testUCCEndFailure(tests *[]doMatchTest) {
	entity := Entity{uccCreationClass, []Attr{
		{step, "sendauthlinktoclient"},
		{stepFailed, trueStr},
		{"mode", "demat"},
	}}
	want := ActionSet{
		properties: []Property{{endFlow, trueStr}},
	}
	*tests = append(*tests, doMatchTest{"ucc end-failure", entity, ruleSets["ucccreation"], ActionSet{}, want})
}

func testComplexWF(tests *[]doMatchTest) {
	ruleSchemas = append(ruleSchemas, RuleSchema{
		complexWFClass,
		[]AttrSchema{
			{step, "str"},
			{stepFailed, "bool"},
			{"type", "enum"},
			{"loc", "enum"},
		},
	})

	setupRuleSetMainForComplexWF()
	setupRuleSet2ForComplexWF()
	setupRuleSet3ForComplexWF()

	testWFBasic(tests)
	testWFThen(tests)
	testWFFail1_1(tests)
	testWFElseChecking(tests)
	testWFElsePPF(tests)
	testWFFail1_3(tests)
	testWFElseAndReturn(tests)
	testWFElseAndExit(tests)
	testWFFail3_2(tests)
	testWFSucc1_3(tests)
	testWFSucc3_2(tests)
}

func setupRuleSetMainForComplexWF() {
	rule1 := Rule{
		[]RulePatternTerm{
			{step, "eq", start},
			{stepFailed, "eq", false},
		},
		RuleActions{
			tasks:      []string{"s1.1"},
			properties: []Property{{nextStep, "s1.1"}},
		},
	}
	rule2 := Rule{
		[]RulePatternTerm{
			{step, "eq", "s1.1"},
			{stepFailed, "eq", false},
			{"type", "eq", "saving"},
		},
		RuleActions{
			thenCall:   "rs2",
			elseCall:   "rs3",
			properties: []Property{{nextStep, "l1.2"}},
		},
	}
	rule3 := Rule{
		[]RulePatternTerm{
			{step, "eq", "s1.1"},
			{stepFailed, "eq", false},
		},
		RuleActions{
			tasks:      []string{"s1.2"},
			properties: []Property{{nextStep, "l1.3"}},
		},
	}
	rule4 := Rule{
		[]RulePatternTerm{
			{step, "eq", "l1.3"},
			{stepFailed, "eq", false},
		},
		RuleActions{
			properties: []Property{{endFlow, trueStr}},
		},
	}
	rule5 := Rule{
		[]RulePatternTerm{
			{step, "eq", "l3.2"},
			{stepFailed, "eq", false},
		},
		RuleActions{
			properties: []Property{{endFlow, trueStr}},
		},
	}
	rule6 := Rule{
		[]RulePatternTerm{
			{stepFailed, "eq", true},
		},
		RuleActions{
			properties: []Property{{endFlow, trueStr}},
		},
	}
	ruleSets["main"] = RuleSet{
		1, complexWFClass, "main",
		[]Rule{rule1, rule2, rule3, rule4, rule5, rule6},
	}
}

func setupRuleSet2ForComplexWF() {
	rule1 := Rule{
		[]RulePatternTerm{
			{step, "eq", "s1.1"},
			{stepFailed, "eq", false},
		},
		RuleActions{
			tasks:      []string{"s2.1"},
			properties: []Property{{nextStep, "l2.1"}},
		},
	}
	ruleSets["rs2"] = RuleSet{
		1, complexWFClass, "rs2",
		[]Rule{rule1},
	}
}

func setupRuleSet3ForComplexWF() {
	rule1 := Rule{
		[]RulePatternTerm{
			{step, "eq", "s1.1"},
			{stepFailed, "eq", false},
			{"loc", "eq", "urban"},
		},
		RuleActions{
			tasks:      []string{"s3.1"},
			properties: []Property{{nextStep, "l3.1"}},
			willReturn: true,
		},
	}
	rule2 := Rule{
		[]RulePatternTerm{
			{step, "eq", "s1.1"},
			{stepFailed, "eq", false},
			{"loc", "eq", "rural"},
		},
		RuleActions{
			tasks:      []string{"s3.2"},
			properties: []Property{{nextStep, "l3.2"}},
			willExit:   true,
		},
	}
	rule3 := Rule{
		[]RulePatternTerm{
			{step, "eq", "s1.1"},
			{stepFailed, "eq", false},
		},
		RuleActions{
			tasks:      []string{"s3.3"},
			properties: []Property{{nextStep, "l3.3"}},
		},
	}
	ruleSets["rs3"] = RuleSet{
		1, complexWFClass, "rs3",
		[]Rule{rule1, rule2, rule3},
	}
}

func testWFBasic(tests *[]doMatchTest) {
	entity := Entity{complexWFClass, []Attr{
		{step, start},
		{stepFailed, falseStr},
		{"type", "saving"},
		{"loc", "urban"},
	}}
	want := ActionSet{
		tasks:      []string{"s1.1"},
		properties: []Property{{nextStep, "s1.1"}},
	}
	*tests = append(*tests, doMatchTest{
		"wf basic",
		entity,
		ruleSets["main"],
		ActionSet{},
		want,
	})
}

func testWFThen(tests *[]doMatchTest) {
	entity := Entity{complexWFClass, []Attr{
		{step, "s1.1"},
		{stepFailed, falseStr},
		{"type", "saving"},
		{"loc", "urban"},
	}}
	want := ActionSet{
		tasks:      []string{"s2.1", "s1.2"},
		properties: []Property{{nextStep, "l1.3"}},
	}
	*tests = append(*tests, doMatchTest{
		"wf then",
		entity,
		ruleSets["main"],
		ActionSet{},
		want,
	})
}

func testWFFail1_1(tests *[]doMatchTest) {
	entity := Entity{complexWFClass, []Attr{
		{step, "s1.1"},
		{stepFailed, trueStr},
		{"type", "saving"},
		{"loc", "urban"},
	}}
	want := ActionSet{
		properties: []Property{{endFlow, trueStr}},
	}
	*tests = append(*tests, doMatchTest{
		"wf fail 1.1",
		entity,
		ruleSets["main"],
		ActionSet{},
		want,
	})
}

func testWFElseChecking(tests *[]doMatchTest) {
	entity := Entity{complexWFClass, []Attr{
		{step, "s1.1"},
		{stepFailed, falseStr},
		{"type", "checking"},
		{"loc", "semirural"},
	}}
	want := ActionSet{
		tasks:      []string{"s3.3", "s1.2"},
		properties: []Property{{nextStep, "l1.3"}},
	}
	*tests = append(*tests, doMatchTest{
		"wf else checking",
		entity,
		ruleSets["main"],
		ActionSet{},
		want,
	})
}

func testWFElsePPF(tests *[]doMatchTest) {
	entity := Entity{complexWFClass, []Attr{
		{step, "s1.1"},
		{stepFailed, falseStr},
		{"type", "ppf"},
		{"loc", "semirural"},
	}}
	want := ActionSet{
		tasks:      []string{"s3.3", "s1.2"},
		properties: []Property{{nextStep, "l1.3"}},
	}
	*tests = append(*tests, doMatchTest{
		"wf else ppf",
		entity,
		ruleSets["main"],
		ActionSet{},
		want,
	})
}

func testWFFail1_3(tests *[]doMatchTest) {
	entity := Entity{complexWFClass, []Attr{
		{step, "l1.3"},
		{stepFailed, trueStr},
		{"type", "checking"},
		{"loc", "semirural"},
	}}
	want := ActionSet{
		properties: []Property{{endFlow, trueStr}},
	}
	*tests = append(*tests, doMatchTest{
		"wf fail 1.3",
		entity,
		ruleSets["main"],
		ActionSet{},
		want,
	})
}

func testWFElseAndReturn(tests *[]doMatchTest) {
	entity := Entity{complexWFClass, []Attr{
		{step, "s1.1"},
		{stepFailed, falseStr},
		{"type", "checking"},
		{"loc", "urban"},
	}}
	want := ActionSet{
		tasks:      []string{"s3.1", "s1.2"},
		properties: []Property{{nextStep, "l1.3"}},
	}
	*tests = append(*tests, doMatchTest{
		"wf else and return",
		entity,
		ruleSets["main"],
		ActionSet{},
		want,
	})
}

func testWFElseAndExit(tests *[]doMatchTest) {
	entity := Entity{complexWFClass, []Attr{
		{step, "s1.1"},
		{stepFailed, falseStr},
		{"type", "checking"},
		{"loc", "rural"},
	}}
	want := ActionSet{
		tasks:      []string{"s3.2"},
		properties: []Property{{nextStep, "l3.2"}},
	}
	*tests = append(*tests, doMatchTest{
		"wf else and exit",
		entity,
		ruleSets["main"],
		ActionSet{},
		want,
	})
}

func testWFFail3_2(tests *[]doMatchTest) {
	entity := Entity{complexWFClass, []Attr{
		{step, "l3.2"},
		{stepFailed, trueStr},
		{"type", "checking"},
		{"loc", "rural"},
	}}
	want := ActionSet{
		properties: []Property{{endFlow, trueStr}},
	}
	*tests = append(*tests, doMatchTest{
		"wf fail 3.2",
		entity,
		ruleSets["main"],
		ActionSet{},
		want,
	})
}

func testWFSucc1_3(tests *[]doMatchTest) {
	entity := Entity{complexWFClass, []Attr{
		{step, "l1.3"},
		{stepFailed, falseStr},
		{"type", "saving"},
		{"loc", "urban"},
	}}
	want := ActionSet{
		properties: []Property{{endFlow, trueStr}},
	}
	*tests = append(*tests, doMatchTest{
		"wf succ 1.3",
		entity,
		ruleSets["main"],
		ActionSet{},
		want,
	})
}

func testWFSucc3_2(tests *[]doMatchTest) {
	entity := Entity{complexWFClass, []Attr{
		{step, "l3.2"},
		{stepFailed, falseStr},
		{"type", "checking"},
		{"loc", "rural"},
	}}
	want := ActionSet{
		properties: []Property{{endFlow, trueStr}},
	}
	*tests = append(*tests, doMatchTest{
		"wf succ 3.2",
		entity,
		ruleSets["main"],
		ActionSet{},
		want,
	})
}
