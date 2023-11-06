/*
Some of the definitions of rulesets below deliberately use a lot of whitespace to keep the code consistent
and to make it easier to understand, add to, and edit these tests
*/

package main

const (
	lastStep        = "laststep"
	lastStepSucc    = "laststepsucc"
	completionLabel = "completionlabel"
	start           = "START"
	endSuccess      = "END-SUCCESS"
	endFailure      = "END-FAILURE"

	uccCreationClass = "ucccreation"
	complexWFClass   = "complexwf"
)

func testUCCCreation(tests *[]doMatchTest) {
	ruleSchemas = append(ruleSchemas, RuleSchema{
		uccCreationClass,
		[]AttrSchema{
			{lastStep, "str"},
			{lastStepSucc, "bool"},
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
			{lastStep, "eq", start},
			{lastStepSucc, "eq", true},
		},
		RuleActions{
			tasks:      []string{"getcustdetails"},
			properties: []Property{{completionLabel, "getcustdetails"}},
		},
	}
	rule2 := Rule{
		[]RulePatternTerm{
			{lastStep, "eq", "getcustdetails"},
			{lastStepSucc, "eq", true},
		},
		RuleActions{
			workflows:  []string{"aof", "kycvalid", "nomauth"},
			properties: []Property{{completionLabel, "readyforauthlink"}},
		},
	}
	rule3 := Rule{
		[]RulePatternTerm{
			{lastStep, "eq", "getcustdetails"},
			{lastStepSucc, "eq", true},
			{"mode", "eq", "physical"},
		},
		RuleActions{
			workflows:  []string{"bankaccvalid"},
			properties: []Property{{completionLabel, "readyforauthlink"}},
		},
	}
	rule4 := Rule{
		[]RulePatternTerm{
			{lastStep, "eq", "getcustdetails"},
			{lastStepSucc, "eq", true},
			{"mode", "eq", "demat"},
		},
		RuleActions{
			workflows:  []string{"dpandbankaccvalid"},
			properties: []Property{{completionLabel, "readyforauthlink"}},
		},
	}
	rule5 := Rule{
		[]RulePatternTerm{
			{lastStep, "eq", "readyforauthlink"},
			{lastStepSucc, "eq", true},
		},
		RuleActions{
			tasks:      []string{"sendauthlinktoclient"},
			properties: []Property{{completionLabel, "sendauthlinktoclient"}},
		},
	}
	rule6 := Rule{
		[]RulePatternTerm{
			{lastStep, "eq", "sendauthlinktoclient"},
			{lastStepSucc, "eq", true},
		},
		RuleActions{
			properties: []Property{{completionLabel, endSuccess}},
		},
	}
	rule7 := Rule{
		[]RulePatternTerm{
			{lastStepSucc, "eq", false},
		},
		RuleActions{
			properties: []Property{{completionLabel, endFailure}},
		},
	}
	ruleSets["ucccreation"] = RuleSet{1, uccCreationClass, "ucccreation",
		[]Rule{rule1, rule2, rule3, rule4, rule5, rule6, rule7},
	}
}

func testUCCStart(tests *[]doMatchTest) {
	entity := Entity{uccCreationClass, []Attr{
		{lastStep, start},
		{lastStepSucc, "true"},
		{"mode", "demat"},
	}}
	want := ActionSet{
		tasks:      []string{"getcustdetails"},
		properties: []Property{{completionLabel, "getcustdetails"}},
	}
	*tests = append(*tests, doMatchTest{"ucc start", entity, ruleSets["ucccreation"], ActionSet{}, want})
}

func testUCCGetCustDetailsDemat(tests *[]doMatchTest) {
	entity := Entity{uccCreationClass, []Attr{
		{lastStep, "getcustdetails"},
		{lastStepSucc, "true"},
		{"mode", "demat"},
	}}
	want := ActionSet{
		workflows:  []string{"aof", "kycvalid", "nomauth", "dpandbankaccvalid"},
		properties: []Property{{completionLabel, "readyforauthlink"}},
	}
	*tests = append(*tests, doMatchTest{"ucc getcustdetails demat", entity, ruleSets["ucccreation"], ActionSet{}, want})
}

func testUCCGetCustDetailsDematFail(tests *[]doMatchTest) {
	entity := Entity{uccCreationClass, []Attr{
		{lastStep, "getcustdetails"},
		{lastStepSucc, "false"},
		{"mode", "demat"},
	}}
	want := ActionSet{
		properties: []Property{{completionLabel, endFailure}},
	}
	*tests = append(*tests, doMatchTest{"ucc getcustdetails demat fail", entity, ruleSets["ucccreation"], ActionSet{}, want})
}

func testUCCGetCustDetailsPhysical(tests *[]doMatchTest) {
	entity := Entity{uccCreationClass, []Attr{
		{lastStep, "getcustdetails"},
		{lastStepSucc, "true"},
		{"mode", "physical"},
	}}
	want := ActionSet{
		workflows:  []string{"aof", "kycvalid", "nomauth", "bankaccvalid"},
		properties: []Property{{completionLabel, "readyforauthlink"}},
	}
	*tests = append(*tests, doMatchTest{"ucc getcustdetails physical", entity, ruleSets["ucccreation"], ActionSet{}, want})
}

func testUCCGetCustDetailsPhysicalFail(tests *[]doMatchTest) {
	entity := Entity{uccCreationClass, []Attr{
		{lastStep, "getcustdetails"},
		{lastStepSucc, "false"},
		{"mode", "physical"},
	}}
	want := ActionSet{
		properties: []Property{{completionLabel, endFailure}},
	}
	*tests = append(*tests, doMatchTest{"ucc getcustdetails physical fail", entity, ruleSets["ucccreation"], ActionSet{}, want})
}

func testUCCReadyForAuthLink(tests *[]doMatchTest) {
	entity := Entity{uccCreationClass, []Attr{
		{lastStep, "readyforauthlink"},
		{lastStepSucc, "true"},
		{"mode", "demat"},
	}}
	want := ActionSet{
		tasks:      []string{"sendauthlinktoclient"},
		properties: []Property{{completionLabel, "sendauthlinktoclient"}},
	}
	*tests = append(*tests, doMatchTest{"ucc readyforauthlink", entity, ruleSets["ucccreation"], ActionSet{}, want})
}

func testUCCReadyForAuthLinkFail(tests *[]doMatchTest) {
	entity := Entity{uccCreationClass, []Attr{
		{lastStep, "readyforauthlink"},
		{lastStepSucc, "false"},
		{"mode", "demat"},
	}}
	want := ActionSet{
		properties: []Property{{completionLabel, endFailure}},
	}
	*tests = append(*tests, doMatchTest{"ucc readyforauthlink fail", entity, ruleSets["ucccreation"], ActionSet{}, want})
}

func testUCCEndSuccess(tests *[]doMatchTest) {
	entity := Entity{uccCreationClass, []Attr{
		{lastStep, "sendauthlinktoclient"},
		{lastStepSucc, "true"},
		{"mode", "demat"},
	}}
	want := ActionSet{
		properties: []Property{{completionLabel, endSuccess}},
	}
	*tests = append(*tests, doMatchTest{"ucc end-success", entity, ruleSets["ucccreation"], ActionSet{}, want})
}

func testUCCEndFailure(tests *[]doMatchTest) {
	entity := Entity{uccCreationClass, []Attr{
		{lastStep, "sendauthlinktoclient"},
		{lastStepSucc, "false"},
		{"mode", "demat"},
	}}
	want := ActionSet{
		properties: []Property{{completionLabel, endFailure}},
	}
	*tests = append(*tests, doMatchTest{"ucc end-failure", entity, ruleSets["ucccreation"], ActionSet{}, want})
}

func testComplexWF(tests *[]doMatchTest) {
	ruleSchemas = append(ruleSchemas, RuleSchema{
		complexWFClass,
		[]AttrSchema{
			{lastStep, "str"},
			{lastStepSucc, "bool"},
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
			{lastStep, "eq", start},
			{lastStepSucc, "eq", true},
		},
		RuleActions{
			tasks:      []string{"s1.1"},
			properties: []Property{{completionLabel, "s1.1"}},
		},
	}
	rule2 := Rule{
		[]RulePatternTerm{
			{lastStep, "eq", "s1.1"},
			{lastStepSucc, "eq", true},
			{"type", "eq", "saving"},
		},
		RuleActions{
			thenCall:   "rs2",
			elseCall:   "rs3",
			properties: []Property{{completionLabel, "l1.2"}},
		},
	}
	rule3 := Rule{
		[]RulePatternTerm{
			{lastStep, "eq", "s1.1"},
			{lastStepSucc, "eq", true},
		},
		RuleActions{
			tasks:      []string{"s1.2"},
			properties: []Property{{completionLabel, "l1.3"}},
		},
	}
	rule4 := Rule{
		[]RulePatternTerm{
			{lastStep, "eq", "l1.3"},
			{lastStepSucc, "eq", true},
		},
		RuleActions{
			properties: []Property{{completionLabel, endSuccess}},
		},
	}
	rule5 := Rule{
		[]RulePatternTerm{
			{lastStep, "eq", "l3.2"},
			{lastStepSucc, "eq", true},
		},
		RuleActions{
			properties: []Property{{completionLabel, endSuccess}},
		},
	}
	rule6 := Rule{
		[]RulePatternTerm{
			{lastStepSucc, "eq", false},
		},
		RuleActions{
			properties: []Property{{completionLabel, endFailure}},
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
			{lastStep, "eq", "s1.1"},
			{lastStepSucc, "eq", true},
		},
		RuleActions{
			tasks:      []string{"s2.1"},
			properties: []Property{{completionLabel, "l2.1"}},
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
			{lastStep, "eq", "s1.1"},
			{lastStepSucc, "eq", true},
			{"loc", "eq", "urban"},
		},
		RuleActions{
			tasks:      []string{"s3.1"},
			properties: []Property{{completionLabel, "l3.1"}},
			willReturn: true,
		},
	}
	rule2 := Rule{
		[]RulePatternTerm{
			{lastStep, "eq", "s1.1"},
			{lastStepSucc, "eq", true},
			{"loc", "eq", "rural"},
		},
		RuleActions{
			tasks:      []string{"s3.2"},
			properties: []Property{{completionLabel, "l3.2"}},
			willExit:   true,
		},
	}
	rule3 := Rule{
		[]RulePatternTerm{
			{lastStep, "eq", "s1.1"},
			{lastStepSucc, "eq", true},
		},
		RuleActions{
			tasks:      []string{"s3.3"},
			properties: []Property{{completionLabel, "l3.3"}},
		},
	}
	ruleSets["rs3"] = RuleSet{
		1, complexWFClass, "rs3",
		[]Rule{rule1, rule2, rule3},
	}
}

func testWFBasic(tests *[]doMatchTest) {
	entity := Entity{complexWFClass, []Attr{
		{lastStep, start},
		{lastStepSucc, "true"},
		{"type", "saving"},
		{"loc", "urban"},
	}}
	want := ActionSet{
		tasks:      []string{"s1.1"},
		properties: []Property{{completionLabel, "s1.1"}},
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
		{lastStep, "s1.1"},
		{lastStepSucc, "true"},
		{"type", "saving"},
		{"loc", "urban"},
	}}
	want := ActionSet{
		tasks:      []string{"s2.1", "s1.2"},
		properties: []Property{{completionLabel, "l1.3"}},
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
		{lastStep, "s1.1"},
		{lastStepSucc, "false"},
		{"type", "saving"},
		{"loc", "urban"},
	}}
	want := ActionSet{
		properties: []Property{{completionLabel, endFailure}},
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
		{lastStep, "s1.1"},
		{lastStepSucc, "true"},
		{"type", "checking"},
		{"loc", "semirural"},
	}}
	want := ActionSet{
		tasks:      []string{"s3.3", "s1.2"},
		properties: []Property{{completionLabel, "l1.3"}},
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
		{lastStep, "s1.1"},
		{lastStepSucc, "true"},
		{"type", "ppf"},
		{"loc", "semirural"},
	}}
	want := ActionSet{
		tasks:      []string{"s3.3", "s1.2"},
		properties: []Property{{completionLabel, "l1.3"}},
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
		{lastStep, "l1.3"},
		{lastStepSucc, "false"},
		{"type", "checking"},
		{"loc", "semirural"},
	}}
	want := ActionSet{
		properties: []Property{{completionLabel, endFailure}},
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
		{lastStep, "s1.1"},
		{lastStepSucc, "true"},
		{"type", "checking"},
		{"loc", "urban"},
	}}
	want := ActionSet{
		tasks:      []string{"s3.1", "s1.2"},
		properties: []Property{{completionLabel, "l1.3"}},
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
		{lastStep, "s1.1"},
		{lastStepSucc, "true"},
		{"type", "checking"},
		{"loc", "rural"},
	}}
	want := ActionSet{
		tasks:      []string{"s3.2"},
		properties: []Property{{completionLabel, "l3.2"}},
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
		{lastStep, "l3.2"},
		{lastStepSucc, "false"},
		{"type", "checking"},
		{"loc", "rural"},
	}}
	want := ActionSet{
		properties: []Property{{completionLabel, endFailure}},
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
		{lastStep, "l1.3"},
		{lastStepSucc, "true"},
		{"type", "saving"},
		{"loc", "urban"},
	}}
	want := ActionSet{
		properties: []Property{{completionLabel, endSuccess}},
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
		{lastStep, "l3.2"},
		{lastStepSucc, "true"},
		{"type", "checking"},
		{"loc", "rural"},
	}}
	want := ActionSet{
		properties: []Property{{completionLabel, endSuccess}},
	}
	*tests = append(*tests, doMatchTest{
		"wf succ 3.2",
		entity,
		ruleSets["main"],
		ActionSet{},
		want,
	})
}
