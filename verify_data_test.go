package main

import (
	"fmt"
	"testing"
)

const (
	incorrectOutputRSMsg = "incorrect output when verifying ruleset with "
	incorrectOutputWFMsg = "incorrect output when verifying workflow with "
	UCCCreation          = "ucccreation"
)

type verifySchemaTest struct {
	name    string
	rs      RuleSchema
	isWF    bool
	want    bool
	wantErr bool
}

func TestVerifySchema(t *testing.T) {
	tests := []verifySchemaTest{}

	// Business rules schema tests
	testCorrectBRSchema(&tests)
	testSchemaEmptyClass(&tests)
	testEmptyPatternSchema(&tests)
	testAttrNameIsNotCruxID(&tests)
	testInvalidValType(&tests)
	testNoValsForEnum(&tests)
	testEnumValIsNotCruxID(&tests)
	testBothTasksAndPropsEmpty(&tests)
	testTaskIsNotCruxID(&tests)
	testPropNameNotCruxID(&tests)

	// Workflow schema tests
	testCorrectWFSchema(&tests)
	testMissingStart(&tests)
	testMissingStep(&tests)
	testAdditionalProps(&tests)
	testMissingNextStep(&tests)
	testTasksAndStepDiscrepancy(&tests)

	fmt.Println("Running", len(tests), "verifyRuleSchema() tests")
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := verifyRuleSchema(tt.rs, tt.isWF)
			if (err != nil) != tt.wantErr {
				t.Errorf("verifyRuleSchema() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("verifyRuleSchema() = %v, want %v", got, tt.want)
			}
		})
	}
}

func testSchemaEmptyClass(tests *[]verifySchemaTest) {
	rs := RuleSchema{class: "",
		patternSchema: []AttrSchema{
			{name: "productname", valType: typeStr},
			{name: "price", valType: typeInt},
			{name: "inwintersale", valType: typeBool},
			{name: "paymenttype", valType: typeEnum, vals: map[string]bool{"cash": true, "card": true}},
			{name: "ismember", valType: typeBool},
		},
		actionSchema: ActionSchema{
			tasks:      []string{"freepen", "freemug", "freebag"},
			properties: []string{"discount", "pointsmult"},
		},
	}
	*tests = append(*tests, verifySchemaTest{
		name:    "schema with empty class",
		rs:      rs,
		isWF:    false,
		want:    false,
		wantErr: true,
	})
}

func testEmptyPatternSchema(tests *[]verifySchemaTest) {
	rs := RuleSchema{class: transactionClass,
		patternSchema: []AttrSchema{},
		actionSchema: ActionSchema{
			tasks:      []string{"freepen", "freemug", "freebag"},
			properties: []string{"discount", "pointsmult"},
		},
	}
	*tests = append(*tests, verifySchemaTest{
		name:    "empty pattern schema",
		rs:      rs,
		isWF:    false,
		want:    false,
		wantErr: true,
	})
}

func testAttrNameIsNotCruxID(tests *[]verifySchemaTest) {
	rs := RuleSchema{class: transactionClass,
		patternSchema: []AttrSchema{
			{name: "1productname", valType: typeStr},
			{name: "price", valType: typeInt},
			{name: "inwintersale", valType: typeBool},
			{name: "paymenttype", valType: typeEnum, vals: map[string]bool{"cash": true, "card": true}},
			{name: "ismember", valType: typeBool},
		},
		actionSchema: ActionSchema{
			tasks:      []string{"freepen", "freemug", "freebag"},
			properties: []string{"discount", "pointsmult"},
		},
	}
	*tests = append(*tests, verifySchemaTest{
		name:    "attr name is not CruxID",
		rs:      rs,
		isWF:    false,
		want:    false,
		wantErr: true,
	})
}

func testInvalidValType(tests *[]verifySchemaTest) {
	rs := RuleSchema{class: transactionClass,
		patternSchema: []AttrSchema{
			{name: "productname", valType: typeStr},
			{name: "price", valType: typeInt},
			{name: "inwintersale", valType: "abc"},
			{name: "paymenttype", valType: typeEnum, vals: map[string]bool{"cash": true, "card": true}},
			{name: "ismember", valType: typeBool},
		},
		actionSchema: ActionSchema{
			tasks:      []string{"freepen", "freemug", "freebag"},
			properties: []string{"discount", "pointsmult"},
		},
	}
	*tests = append(*tests, verifySchemaTest{
		name:    "invalid value type",
		rs:      rs,
		isWF:    false,
		want:    false,
		wantErr: true,
	})
}

func testNoValsForEnum(tests *[]verifySchemaTest) {
	rs := RuleSchema{class: transactionClass,
		patternSchema: []AttrSchema{
			{name: "productname", valType: typeStr},
			{name: "price", valType: typeInt},
			{name: "inwintersale", valType: typeBool},
			{name: "paymenttype", valType: typeEnum, vals: map[string]bool{}},
			{name: "ismember", valType: typeBool},
		},
		actionSchema: ActionSchema{
			tasks:      []string{"freepen", "freemug", "freebag"},
			properties: []string{"discount", "pointsmult"},
		},
	}
	*tests = append(*tests, verifySchemaTest{
		name:    "no vals for enum",
		rs:      rs,
		isWF:    false,
		want:    false,
		wantErr: true,
	})
}

func testEnumValIsNotCruxID(tests *[]verifySchemaTest) {
	rs := RuleSchema{class: transactionClass,
		patternSchema: []AttrSchema{
			{name: "productname", valType: typeStr},
			{name: "price", valType: typeInt},
			{name: "inwintersale", valType: typeBool},
			{name: "paymenttype", valType: typeEnum, vals: map[string]bool{"1cash": true, "card": true}},
			{name: "ismember", valType: typeBool},
		},
		actionSchema: ActionSchema{
			tasks:      []string{"freepen", "freemug", "freebag"},
			properties: []string{"discount", "pointsmult"},
		},
	}
	*tests = append(*tests, verifySchemaTest{
		name:    "enum val is not CruxID",
		rs:      rs,
		isWF:    false,
		want:    false,
		wantErr: true,
	})
}

func testBothTasksAndPropsEmpty(tests *[]verifySchemaTest) {
	rs := RuleSchema{class: transactionClass,
		patternSchema: []AttrSchema{
			{name: "productname", valType: typeStr},
			{name: "price", valType: typeInt},
			{name: "inwintersale", valType: typeBool},
			{name: "paymenttype", valType: typeEnum, vals: map[string]bool{"cash": true, "card": true}},
			{name: "ismember", valType: typeBool},
		},
		actionSchema: ActionSchema{
			tasks:      []string{},
			properties: []string{},
		},
	}
	*tests = append(*tests, verifySchemaTest{
		name:    "both tasks and properties empty",
		rs:      rs,
		isWF:    false,
		want:    false,
		wantErr: true,
	})
}

func testTaskIsNotCruxID(tests *[]verifySchemaTest) {
	rs := RuleSchema{class: transactionClass,
		patternSchema: []AttrSchema{
			{name: "productname", valType: typeStr},
			{name: "price", valType: typeInt},
			{name: "inwintersale", valType: typeBool},
			{name: "paymenttype", valType: typeEnum, vals: map[string]bool{"cash": true, "card": true}},
			{name: "ismember", valType: typeBool},
		},
		actionSchema: ActionSchema{
			tasks:      []string{"freepen", "free*mug", "freebag"},
			properties: []string{"discount", "pointsmult"},
		},
	}
	*tests = append(*tests, verifySchemaTest{
		name:    "task is not CruxID",
		rs:      rs,
		isWF:    false,
		want:    false,
		wantErr: true,
	})
}

func testPropNameNotCruxID(tests *[]verifySchemaTest) {
	rs := RuleSchema{class: transactionClass,
		patternSchema: []AttrSchema{
			{name: "productname", valType: typeStr},
			{name: "price", valType: typeInt},
			{name: "inwintersale", valType: typeBool},
			{name: "paymenttype", valType: typeEnum, vals: map[string]bool{"cash": true, "card": true}},
			{name: "ismember", valType: typeBool},
		},
		actionSchema: ActionSchema{
			tasks:      []string{"freepen", "freemug", "freebag"},
			properties: []string{"1discount", "pointsmult"},
		},
	}
	*tests = append(*tests, verifySchemaTest{
		name:    "property name is not CruxID",
		rs:      rs,
		isWF:    false,
		want:    false,
		wantErr: true,
	})
}

func testCorrectBRSchema(tests *[]verifySchemaTest) {
	rs := RuleSchema{class: transactionClass,
		patternSchema: []AttrSchema{
			{name: "productname", valType: typeStr},
			{name: "price", valType: typeInt},
			{name: "inwintersale", valType: typeBool},
			{name: "paymenttype", valType: typeEnum, vals: map[string]bool{"cash": true, "card": true}},
			{name: "ismember", valType: typeBool},
		},
		actionSchema: ActionSchema{
			tasks:      []string{"freepen", "freemug", "freebag"},
			properties: []string{"discount", "pointsmult"},
		},
	}
	*tests = append(*tests, verifySchemaTest{
		name:    "correct business-rules schema",
		rs:      rs,
		isWF:    false,
		want:    true,
		wantErr: false,
	})
}

func testMissingStart(tests *[]verifySchemaTest) {
	rs := RuleSchema{
		class: uccCreationClass,
		patternSchema: []AttrSchema{
			{name: step, valType: typeEnum,
				vals: map[string]bool{"getcustdetails": true, "aof": true, "sendauthlinktoclient": true},
			},
			{name: stepFailed, valType: typeBool},
			{name: "mode", valType: typeEnum, vals: map[string]bool{"physical": true, "demat": true}},
		},
		actionSchema: ActionSchema{
			tasks:      []string{"getcustdetails", "aof", "sendauthlinktoclient"},
			properties: []string{nextStep, done},
		},
	}
	*tests = append(*tests, verifySchemaTest{
		name:    "missing START",
		rs:      rs,
		isWF:    true,
		want:    false,
		wantErr: true,
	})
}

func testMissingStep(tests *[]verifySchemaTest) {
	rs := RuleSchema{
		class: uccCreationClass,
		patternSchema: []AttrSchema{
			{name: stepFailed, valType: typeBool},
			{name: "mode", valType: typeEnum, vals: map[string]bool{"physical": true, "demat": true}},
		},
		actionSchema: ActionSchema{
			tasks:      []string{"getcustdetails", "aof", "sendauthlinktoclient"},
			properties: []string{nextStep, done},
		},
	}
	*tests = append(*tests, verifySchemaTest{
		name:    "missing step",
		rs:      rs,
		isWF:    true,
		want:    false,
		wantErr: true,
	})
}

func testAdditionalProps(tests *[]verifySchemaTest) {
	rs := RuleSchema{
		class: uccCreationClass,
		patternSchema: []AttrSchema{
			{name: step, valType: typeEnum,
				vals: map[string]bool{start: true, "getcustdetails": true, "aof": true, "sendauthlinktoclient": true},
			},
			{name: stepFailed, valType: typeBool},
			{name: "mode", valType: typeEnum, vals: map[string]bool{"physical": true, "demat": true}},
		},
		actionSchema: ActionSchema{
			tasks:      []string{"getcustdetails", "aof", "sendauthlinktoclient"},
			properties: []string{nextStep, done, "abcd"},
		},
	}
	*tests = append(*tests, verifySchemaTest{
		name:    "additional property other than nextstep and done",
		rs:      rs,
		isWF:    true,
		want:    false,
		wantErr: true,
	})
}

func testMissingNextStep(tests *[]verifySchemaTest) {
	rs := RuleSchema{
		class: uccCreationClass,
		patternSchema: []AttrSchema{
			{name: step, valType: typeEnum,
				vals: map[string]bool{start: true, "getcustdetails": true, "aof": true, "sendauthlinktoclient": true},
			},
			{name: stepFailed, valType: typeBool},
			{name: "mode", valType: typeEnum, vals: map[string]bool{"physical": true, "demat": true}},
		},
		actionSchema: ActionSchema{
			tasks:      []string{"getcustdetails", "aof", "sendauthlinktoclient"},
			properties: []string{done, "abcd"},
		},
	}
	*tests = append(*tests, verifySchemaTest{
		name:    "missing nextstep",
		rs:      rs,
		isWF:    true,
		want:    false,
		wantErr: true,
	})
}

func testTasksAndStepDiscrepancy(tests *[]verifySchemaTest) {
	rs := RuleSchema{
		class: uccCreationClass,
		patternSchema: []AttrSchema{
			{name: step, valType: typeEnum,
				vals: map[string]bool{start: true, "getcustdetails": true, "aof": true, "sendauthlinktoclient": true},
			},
			{name: stepFailed, valType: typeBool},
			{name: "mode", valType: typeEnum, vals: map[string]bool{"physical": true, "demat": true}},
		},
		actionSchema: ActionSchema{
			tasks:      []string{"getcustinfo", "aof", "sendauthlinktoclient"},
			properties: []string{nextStep, done},
		},
	}
	*tests = append(*tests, verifySchemaTest{
		name:    "tasks and steps discrepancy",
		rs:      rs,
		isWF:    true,
		want:    false,
		wantErr: true,
	})
}

func testCorrectWFSchema(tests *[]verifySchemaTest) {
	rs := RuleSchema{
		class: uccCreationClass,
		patternSchema: []AttrSchema{
			{name: step, valType: typeEnum,
				vals: map[string]bool{start: true, "getcustdetails": true, "aof": true, "sendauthlinktoclient": true},
			},
			{name: stepFailed, valType: typeBool},
			{name: "mode", valType: typeEnum, vals: map[string]bool{"physical": true, "demat": true}},
		},
		actionSchema: ActionSchema{
			tasks:      []string{"getcustdetails", "aof", "sendauthlinktoclient"},
			properties: []string{nextStep, done},
		},
	}
	*tests = append(*tests, verifySchemaTest{
		name:    "correct workflow schema",
		rs:      rs,
		isWF:    true,
		want:    true,
		wantErr: false,
	})
}

func testCorrectRS(t *testing.T) {
	ok, err := verifyRuleSet(ruleSets[mainRS], false)
	if !ok || err != nil {
		t.Errorf(incorrectOutputRSMsg + "no issues")
	}
}

// In each of the rule-pattern tests below, a rule-pattern is modified temporarily.
// After each test, we must reset the rule-pattern to the correct one below before
// moving on to the next test.
var correctRP = []RulePatternTerm{
	{"product", opEQ, "jacket"},
	{"price", opGT, 50.0},
}

func testInvalidAttrName(t *testing.T) {
	ruleSets[mainRS].rules[1].rulePattern = []RulePatternTerm{
		{"product", opEQ, "jacket"},
		{"priceabc", opGT, 50.0},
	}
	ok, err := verifyRuleSet(ruleSets[mainRS], false)
	if ok || err == nil {
		t.Errorf(incorrectOutputRSMsg + "invalid attr name")
	}
	ruleSets[mainRS].rules[1].rulePattern = correctRP
}

func testTaskAsAttrName(t *testing.T) {
	ruleSets[mainRS].rules[1].rulePattern = []RulePatternTerm{
		{"product", opEQ, "jacket"},
		{"freejar", opEQ, true},
	}
	ok, err := verifyRuleSet(ruleSets[mainRS], false)
	if !ok || err != nil {
		t.Errorf(incorrectOutputRSMsg + "a task 'tag' as an attribute name")
	}
	ruleSets[mainRS].rules[1].rulePattern = correctRP
}

func testWrongAttrValType(t *testing.T) {
	ruleSets[mainRS].rules[1].rulePattern = []RulePatternTerm{
		{"product", opEQ, "jacket"},
		{"price", opGT, "abc"},
	}
	ok, err := verifyRuleSet(ruleSets[mainRS], false)
	if ok || err == nil {
		t.Errorf(incorrectOutputRSMsg + "wrong attribute value type")
	}
	ruleSets[mainRS].rules[1].rulePattern = correctRP
}

func testInvalidOp(t *testing.T) {
	ruleSets[mainRS].rules[1].rulePattern = []RulePatternTerm{
		{"product", opEQ, "jacket"},
		{"price", "greater than", 50.0},
	}
	ok, err := verifyRuleSet(ruleSets[mainRS], false)
	if ok || err == nil {
		t.Errorf(incorrectOutputRSMsg + "invalid operation")
	}
	ruleSets[mainRS].rules[1].rulePattern = correctRP
}

// In each of the rule-action tests below, a rule-action is modified temporarily.
// After each test, we must reset the rule-action to the correct one below before
// moving on to the next test.
var correctRA RuleActions = RuleActions{
	tasks:      []string{"freemug", "freejar", "freeplant"},
	properties: []Property{{"discount", "20"}},
}

func testTaskNotInSchema(t *testing.T) {
	ruleSets[mainRS].rules[3].ruleActions = RuleActions{
		tasks:      []string{"freemug", "freeeraser"},
		properties: []Property{{"discount", "20"}},
	}
	ok, err := verifyRuleSet(ruleSets[mainRS], false)
	if ok || err == nil {
		t.Errorf(incorrectOutputRSMsg + "task not in schema")
	}
	ruleSets[mainRS].rules[3].ruleActions = correctRA
}

func testPropNameNotInSchema(t *testing.T) {
	ruleSets[mainRS].rules[3].ruleActions = RuleActions{
		tasks:      []string{"freemug", "freejar", "freeplant"},
		properties: []Property{{"cashback", "5"}},
	}
	ok, err := verifyRuleSet(ruleSets[mainRS], false)
	if ok || err == nil {
		t.Errorf(incorrectOutputRSMsg + "property name not in schema")
	}
	ruleSets[mainRS].rules[3].ruleActions = correctRA
}

func testBothReturnAndExit(t *testing.T) {
	ruleSets[mainRS].rules[3].ruleActions = RuleActions{
		tasks:      []string{"freemug", "freejar", "freeplant"},
		properties: []Property{{"discount", "20"}},
		willReturn: true,
		willExit:   true,
	}
	ok, err := verifyRuleSet(ruleSets[mainRS], false)
	if ok || err == nil {
		t.Errorf(incorrectOutputRSMsg + "both RETURN and EXIT instructions")
	}
	ruleSets[mainRS].rules[3].ruleActions = correctRA
}

func testCorrectWF(t *testing.T) {
	ok, err := verifyRuleSet(ruleSets[UCCCreation], true)
	if !ok || err != nil {
		t.Errorf(incorrectOutputWFMsg + "no issues")
	}
}

func testWFRuleMissingStep(t *testing.T) {
	ruleSets[UCCCreation].rules[1].rulePattern = []RulePatternTerm{
		{stepFailed, opEQ, false},
		{"mode", opEQ, "physical"},
	}
	ok, err := verifyRuleSet(ruleSets[UCCCreation], true)
	if ok || err == nil {
		t.Errorf(incorrectOutputWFMsg + "a rule missing 'step'")
	}
	// Reset to original correct rule-pattern
	ruleSets[UCCCreation].rules[1].rulePattern = []RulePatternTerm{
		{step, opEQ, "getcustdetails"},
		{stepFailed, opEQ, false},
		{"mode", opEQ, "physical"},
	}
}

// In each of the (workflow) rule-action tests below, a rule-action is modified temporarily.
// After each test, we must reset the rule-action to the correct one below before
// moving on to the next test.
var correctWorkflowRA RuleActions = RuleActions{
	tasks:      []string{"aof", "kycvalid", "nomauth", "bankaccvalid"},
	properties: []Property{{nextStep, "aof"}},
}

func testWFRuleMissingBothNSAndDone(t *testing.T) {
	ruleSets[UCCCreation].rules[1].ruleActions = RuleActions{
		tasks:      []string{"aof", "kycvalid", "nomauth", "bankaccvalid"},
		properties: []Property{},
	}
	ok, err := verifyRuleSet(ruleSets[UCCCreation], true)
	if ok || err == nil {
		t.Errorf(incorrectOutputWFMsg + "a rule missing both 'nextstep' and 'done'")
	}
	ruleSets[UCCCreation].rules[1].ruleActions = correctWorkflowRA
}

func testWFNoTasksAndNotDone(t *testing.T) {
	ruleSets[UCCCreation].rules[1].ruleActions = RuleActions{
		tasks:      []string{},
		properties: []Property{{nextStep, "abc"}},
	}
	ok, err := verifyRuleSet(ruleSets[UCCCreation], true)
	if ok || err == nil {
		t.Errorf(incorrectOutputWFMsg + "a rule with no tasks and no 'done=true'")
	}
	ruleSets[UCCCreation].rules[1].ruleActions = correctWorkflowRA
}

func testWFNextStepValNotInTasks(t *testing.T) {
	ruleSets[UCCCreation].rules[1].ruleActions = RuleActions{
		tasks:      []string{"aof", "kycvalid", "nomauth", "bankaccvalid"},
		properties: []Property{{nextStep, "abcd"}},
	}
	ok, err := verifyRuleSet(ruleSets[UCCCreation], true)
	if ok || err == nil {
		t.Errorf(incorrectOutputWFMsg + "a 'nextstep' value not in its rule's 'tasks'")
	}
	ruleSets[UCCCreation].rules[1].ruleActions = correctWorkflowRA
}

func TestVerifyRuleSet(t *testing.T) {

	// Business rules tests
	setupPurchaseRuleSchema()
	setupRuleSetForPurchases()
	testCorrectRS(t)
	testInvalidAttrName(t)
	testTaskAsAttrName(t)
	testWrongAttrValType(t)
	testInvalidOp(t)
	testTaskNotInSchema(t)
	testPropNameNotInSchema(t)
	testBothReturnAndExit(t)

	// Workflow tests
	setupUCCCreationSchema()
	setupUCCCreationRuleSet()
	testCorrectWF(t)
	testWFRuleMissingStep(t)
	testWFRuleMissingBothNSAndDone(t)
	testWFNoTasksAndNotDone(t)
	testWFNextStepValNotInTasks(t)
}
