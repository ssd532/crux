/*
This file contains WFE tests for doMatch()

Some of the definitions of rulesets below deliberately use a lot of whitespace to keep the code consistent
and to make it easier to understand, add to, and edit these tests
*/

package main

const (
	uccCreationClass = "ucccreation"
	prepareAOFClass  = "prepareaof"
	validateAOFClass = "validateaof"
	complexWFClass   = "complexwf"
)

func testUCCCreation(tests *[]doMatchTest) {
	ruleSchemas = append(ruleSchemas, RuleSchema{
		class: uccCreationClass,
		patternSchema: []AttrSchema{
			{name: step, valType: typeStr},
			{name: stepFailed, valType: typeBool},
			{name: "mode", valType: typeEnum},
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
			{step, opEQ, start},
			{stepFailed, opEQ, false},
		},
		RuleActions{
			tasks:      []string{"getcustdetails"},
			properties: []Property{{nextStep, "getcustdetails"}},
		},
	}
	rule2 := Rule{
		[]RulePatternTerm{
			{step, opEQ, "getcustdetails"},
			{stepFailed, opEQ, false},
		},
		RuleActions{
			tasks:      []string{"aof", "kycvalid", "nomauth"},
			properties: []Property{{nextStep, "readyforauthlink"}},
		},
	}
	rule3 := Rule{
		[]RulePatternTerm{
			{step, opEQ, "getcustdetails"},
			{stepFailed, opEQ, false},
			{"mode", opEQ, "physical"},
		},
		RuleActions{
			tasks:      []string{"bankaccvalid"},
			properties: []Property{{nextStep, "readyforauthlink"}},
		},
	}
	rule4 := Rule{
		[]RulePatternTerm{
			{step, opEQ, "getcustdetails"},
			{stepFailed, opEQ, false},
			{"mode", opEQ, "demat"},
		},
		RuleActions{
			tasks:      []string{"dpandbankaccvalid"},
			properties: []Property{{nextStep, "readyforauthlink"}},
		},
	}
	rule5 := Rule{
		[]RulePatternTerm{
			{step, opEQ, "readyforauthlink"},
			{stepFailed, opEQ, false},
		},
		RuleActions{
			tasks:      []string{"sendauthlinktoclient"},
			properties: []Property{{nextStep, "sendauthlinktoclient"}},
		},
	}
	rule6 := Rule{
		[]RulePatternTerm{
			{step, opEQ, "sendauthlinktoclient"},
			{stepFailed, opEQ, false},
		},
		RuleActions{
			properties: []Property{{done, trueStr}},
		},
	}
	rule7 := Rule{
		[]RulePatternTerm{
			{stepFailed, opEQ, true},
		},
		RuleActions{
			properties: []Property{{done, trueStr}},
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
		properties: []Property{{done, trueStr}},
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
		properties: []Property{{done, trueStr}},
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
		properties: []Property{{done, trueStr}},
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
		properties: []Property{{done, trueStr}},
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
		properties: []Property{{done, trueStr}},
	}
	*tests = append(*tests, doMatchTest{"ucc end-failure", entity, ruleSets["ucccreation"], ActionSet{}, want})
}

func testPrepareAOF(tests *[]doMatchTest) {
	ruleSchemas = append(ruleSchemas, RuleSchema{
		class: prepareAOFClass,
		patternSchema: []AttrSchema{
			{name: step, valType: typeStr},
			{name: stepFailed, valType: typeBool},
		},
	})

	setupRuleSetForPrepareAOF()

	testDownloadAOF(tests)
	testDownloadAOFFail(tests)
	testPrintAOF(tests)
	testSignAOF(tests)
	testSignAOFFail(tests)
	testReceiveSignedAOF(tests)
	testUploadAOF(tests)
	testPrepareAOFEnd(tests)
}

func testDownloadAOF(tests *[]doMatchTest) {
	entity := Entity{prepareAOFClass, []Attr{
		{step, start},
		{stepFailed, falseStr},
	}}
	want := ActionSet{
		tasks:      []string{"downloadform"},
		properties: []Property{{nextStep, "downloadform"}},
	}
	*tests = append(*tests, doMatchTest{"download aof", entity, ruleSets["prepareaof"], ActionSet{}, want})
}

func testDownloadAOFFail(tests *[]doMatchTest) {
	entity := Entity{prepareAOFClass, []Attr{
		{step, "downloadform"},
		{stepFailed, trueStr},
	}}
	want := ActionSet{
		properties: []Property{{done, trueStr}},
	}
	*tests = append(*tests, doMatchTest{"download aof fail", entity, ruleSets["prepareaof"], ActionSet{}, want})
}

func testPrintAOF(tests *[]doMatchTest) {
	entity := Entity{prepareAOFClass, []Attr{
		{step, "downloadform"},
		{stepFailed, falseStr},
	}}
	want := ActionSet{
		tasks:      []string{"printprefilledform"},
		properties: []Property{{nextStep, "printprefilledform"}},
	}
	*tests = append(*tests, doMatchTest{"print prefilled aof", entity, ruleSets["prepareaof"], ActionSet{}, want})
}

func testSignAOF(tests *[]doMatchTest) {
	entity := Entity{prepareAOFClass, []Attr{
		{step, "printprefilledform"},
		{stepFailed, falseStr},
	}}
	want := ActionSet{
		tasks:      []string{"signform"},
		properties: []Property{{nextStep, "signform"}},
	}
	*tests = append(*tests, doMatchTest{"sign aof", entity, ruleSets["prepareaof"], ActionSet{}, want})
}

func testSignAOFFail(tests *[]doMatchTest) {
	entity := Entity{prepareAOFClass, []Attr{
		{step, "signform"},
		{stepFailed, trueStr},
	}}
	want := ActionSet{
		properties: []Property{{done, trueStr}},
	}
	*tests = append(*tests, doMatchTest{"sign aof fail", entity, ruleSets["prepareaof"], ActionSet{}, want})
}

func testReceiveSignedAOF(tests *[]doMatchTest) {
	entity := Entity{prepareAOFClass, []Attr{
		{step, "signform"},
		{stepFailed, falseStr},
	}}
	want := ActionSet{
		tasks:      []string{"receivesignedform"},
		properties: []Property{{nextStep, "receivesignedform"}},
	}
	*tests = append(*tests, doMatchTest{"receive signed aof", entity, ruleSets["prepareaof"], ActionSet{}, want})
}

func testUploadAOF(tests *[]doMatchTest) {
	entity := Entity{prepareAOFClass, []Attr{
		{step, "receivesignedform"},
		{stepFailed, falseStr},
	}}
	want := ActionSet{
		tasks:      []string{"uploadsignedform"},
		properties: []Property{{nextStep, "uploadsignedform"}},
	}
	*tests = append(*tests, doMatchTest{"upload signed aof", entity, ruleSets["prepareaof"], ActionSet{}, want})
}

func testPrepareAOFEnd(tests *[]doMatchTest) {
	entity := Entity{prepareAOFClass, []Attr{
		{step, "uploadsignedform"},
		{stepFailed, falseStr},
	}}
	want := ActionSet{
		properties: []Property{{done, trueStr}},
	}
	*tests = append(*tests, doMatchTest{"prepare aof end", entity, ruleSets["prepareaof"], ActionSet{}, want})
}

func setupRuleSetForPrepareAOF() {
	rule1 := Rule{
		[]RulePatternTerm{
			{step, opEQ, start},
			{stepFailed, opEQ, false},
		},
		RuleActions{
			tasks:      []string{"downloadform"},
			properties: []Property{{nextStep, "downloadform"}},
		},
	}
	rule2 := Rule{
		[]RulePatternTerm{
			{step, opEQ, "downloadform"},
			{stepFailed, opEQ, false},
		},
		RuleActions{
			tasks:      []string{"printprefilledform"},
			properties: []Property{{nextStep, "printprefilledform"}},
		},
	}
	rule3 := Rule{
		[]RulePatternTerm{
			{step, opEQ, "printprefilledform"},
			{stepFailed, opEQ, false},
		},
		RuleActions{
			tasks:      []string{"signform"},
			properties: []Property{{nextStep, "signform"}},
		},
	}
	rule4 := Rule{
		[]RulePatternTerm{
			{step, opEQ, "signform"},
			{stepFailed, opEQ, false},
		},
		RuleActions{
			tasks:      []string{"receivesignedform"},
			properties: []Property{{nextStep, "receivesignedform"}},
		},
	}
	rule5 := Rule{
		[]RulePatternTerm{
			{step, opEQ, "receivesignedform"},
			{stepFailed, opEQ, false},
		},
		RuleActions{
			tasks:      []string{"uploadsignedform"},
			properties: []Property{{nextStep, "uploadsignedform"}},
		},
	}
	rule6 := Rule{
		[]RulePatternTerm{
			{step, opEQ, "uploadsignedform"},
			{stepFailed, opEQ, false},
		},
		RuleActions{
			tasks:      []string{},
			properties: []Property{{done, trueStr}},
		},
	}
	rule7 := Rule{
		[]RulePatternTerm{
			{stepFailed, opEQ, true},
		},
		RuleActions{
			properties: []Property{{done, trueStr}},
		},
	}
	ruleSets["prepareaof"] = RuleSet{1, prepareAOFClass, "prepareaof",
		[]Rule{rule1, rule2, rule3, rule4, rule5, rule6, rule7},
	}
}

func testValidateAOF(tests *[]doMatchTest) {
	ruleSchemas = append(ruleSchemas, RuleSchema{
		class: validateAOFClass,
		patternSchema: []AttrSchema{
			{name: step, valType: typeStr},
			{name: stepFailed, valType: typeBool},
			{name: "aofexists", valType: typeBool},
		},
	})

	setupRuleSetForValidateAOF()

	testValidateExistingAOF(tests)
	testValidateAOFStart(tests)
	testAOFGetResponseFromRTA(tests)
	testValidateAOFEnd(tests)
}

func testValidateExistingAOF(tests *[]doMatchTest) {
	entity := Entity{validateAOFClass, []Attr{
		{step, start},
		{stepFailed, falseStr},
		{"aofexists", trueStr},
	}}
	want := ActionSet{
		properties: []Property{{done, trueStr}},
	}
	*tests = append(*tests, doMatchTest{"validate existing aof", entity, ruleSets["validateaof"], ActionSet{}, want})
}

func testValidateAOFStart(tests *[]doMatchTest) {
	entity := Entity{validateAOFClass, []Attr{
		{step, start},
		{stepFailed, falseStr},
		{"aofexists", falseStr},
	}}
	want := ActionSet{
		tasks:      []string{"sendaoftorta"},
		properties: []Property{{nextStep, "sendaoftorta"}},
	}
	*tests = append(*tests, doMatchTest{"send aof to rta", entity, ruleSets["validateaof"], ActionSet{}, want})
}

func testAOFGetResponseFromRTA(tests *[]doMatchTest) {
	entity := Entity{validateAOFClass, []Attr{
		{step, "sendaoftorta"},
		{stepFailed, falseStr},
		{"aofexists", falseStr},
	}}
	want := ActionSet{
		tasks:      []string{"getresponsefromrta"},
		properties: []Property{{nextStep, "getresponsefromrta"}},
	}
	*tests = append(*tests, doMatchTest{"aof - get response from rta", entity, ruleSets["validateaof"], ActionSet{}, want})
}

func testValidateAOFEnd(tests *[]doMatchTest) {
	entity := Entity{validateAOFClass, []Attr{
		{step, "getresponsefromrta"},
		{stepFailed, falseStr},
		{"aofexists", falseStr},
	}}
	want := ActionSet{
		properties: []Property{{done, trueStr}},
	}
	*tests = append(*tests, doMatchTest{"validate aof end", entity, ruleSets["validateaof"], ActionSet{}, want})
}

func setupRuleSetForValidateAOF() {
	rule1 := Rule{
		[]RulePatternTerm{
			{step, opEQ, start},
			{stepFailed, opEQ, false},
			{"aofexists", opEQ, true},
		},
		RuleActions{
			properties: []Property{{done, trueStr}},
		},
	}
	rule2 := Rule{
		[]RulePatternTerm{
			{step, opEQ, start},
			{stepFailed, opEQ, false},
			{"aofexists", opEQ, false},
		},
		RuleActions{
			tasks:      []string{"sendaoftorta"},
			properties: []Property{{nextStep, "sendaoftorta"}},
		},
	}
	rule3 := Rule{
		[]RulePatternTerm{
			{step, opEQ, "sendaoftorta"},
			{stepFailed, opEQ, false},
			{"aofexists", opEQ, false},
		},
		RuleActions{
			tasks:      []string{"getresponsefromrta"},
			properties: []Property{{nextStep, "getresponsefromrta"}},
		},
	}
	rule4 := Rule{
		[]RulePatternTerm{
			{step, opEQ, "getresponsefromrta"},
			{stepFailed, opEQ, false},
			{"aofexists", opEQ, false},
		},
		RuleActions{
			properties: []Property{{done, trueStr}},
		},
	}
	rule5 := Rule{
		[]RulePatternTerm{
			{stepFailed, opEQ, true},
		},
		RuleActions{
			properties: []Property{{done, trueStr}},
		},
	}
	ruleSets["validateaof"] = RuleSet{1, validateAOFClass, "validateaof",
		[]Rule{rule1, rule2, rule3, rule4, rule5},
	}
}

func testComplexWF(tests *[]doMatchTest) {
	ruleSchemas = append(ruleSchemas, RuleSchema{
		class: complexWFClass,
		patternSchema: []AttrSchema{
			{name: step, valType: typeStr},
			{name: stepFailed, valType: typeBool},
			{name: "type", valType: typeEnum},
			{name: "loc", valType: typeEnum},
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
			{step, opEQ, start},
			{stepFailed, opEQ, false},
		},
		RuleActions{
			tasks:      []string{"s1.1"},
			properties: []Property{{nextStep, "s1.1"}},
		},
	}
	rule2 := Rule{
		[]RulePatternTerm{
			{step, opEQ, "s1.1"},
			{stepFailed, opEQ, false},
			{"type", opEQ, "saving"},
		},
		RuleActions{
			thenCall:   "rs2",
			elseCall:   "rs3",
			properties: []Property{{nextStep, "l1.2"}},
		},
	}
	rule3 := Rule{
		[]RulePatternTerm{
			{step, opEQ, "s1.1"},
			{stepFailed, opEQ, false},
		},
		RuleActions{
			tasks:      []string{"s1.2"},
			properties: []Property{{nextStep, "l1.3"}},
		},
	}
	rule4 := Rule{
		[]RulePatternTerm{
			{step, opEQ, "l1.3"},
			{stepFailed, opEQ, false},
		},
		RuleActions{
			properties: []Property{{done, trueStr}},
		},
	}
	rule5 := Rule{
		[]RulePatternTerm{
			{step, opEQ, "l3.2"},
			{stepFailed, opEQ, false},
		},
		RuleActions{
			properties: []Property{{done, trueStr}},
		},
	}
	rule6 := Rule{
		[]RulePatternTerm{
			{stepFailed, opEQ, true},
		},
		RuleActions{
			properties: []Property{{done, trueStr}},
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
			{step, opEQ, "s1.1"},
			{stepFailed, opEQ, false},
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
			{step, opEQ, "s1.1"},
			{stepFailed, opEQ, false},
			{"loc", opEQ, "urban"},
		},
		RuleActions{
			tasks:      []string{"s3.1"},
			properties: []Property{{nextStep, "l3.1"}},
			willReturn: true,
		},
	}
	rule2 := Rule{
		[]RulePatternTerm{
			{step, opEQ, "s1.1"},
			{stepFailed, opEQ, false},
			{"loc", opEQ, "rural"},
		},
		RuleActions{
			tasks:      []string{"s3.2"},
			properties: []Property{{nextStep, "l3.2"}},
			willExit:   true,
		},
	}
	rule3 := Rule{
		[]RulePatternTerm{
			{step, opEQ, "s1.1"},
			{stepFailed, opEQ, false},
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
		properties: []Property{{done, trueStr}},
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
		properties: []Property{{done, trueStr}},
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
		properties: []Property{{done, trueStr}},
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
		properties: []Property{{done, trueStr}},
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
		properties: []Property{{done, trueStr}},
	}
	*tests = append(*tests, doMatchTest{
		"wf succ 3.2",
		entity,
		ruleSets["main"],
		ActionSet{},
		want,
	})
}
