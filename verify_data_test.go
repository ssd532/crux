/*
This file contains TestVerifySchema() and TestVerifyRuleSet(). These two functions run tests for
verifyRuleSchema() and verifyRuleSet() respectively.
*/

package main

import (
	"testing"
)

const (
	incorrectOutputRSMsg = "incorrect output when verifying ruleset with "
	incorrectOutputWFMsg = "incorrect output when verifying workflow with "
	uccCreation          = "ucccreation"
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

	/* Business rules schema tests */
	// the only test that involves no error, because the schema is correct
	testCorrectBRSchema(&tests)
	// in the rest of these tests, verifyRuleSchema() should return an error
	testSchemaEmptyClass(&tests)
	testEmptyPatternSchema(&tests)
	testAttrNameIsNotCruxID(&tests)
	testInvalidValType(&tests)
	testNoValsForEnum(&tests)
	testEnumValIsNotCruxID(&tests)
	testBothTasksAndPropsEmpty(&tests)
	testTaskIsNotCruxID(&tests)
	testPropNameNotCruxID(&tests)

	/* Workflow schema tests */
	// the only test that involves no error, because the workflow schema is correct
	testCorrectWFSchema(&tests)
	// in the rest of these tests, verifyRuleSchema() should return an error
	testMissingStart(&tests)
	testMissingStep(&tests)
	testAdditionalProps(&tests)
	testMissingNextStep(&tests)
	testTasksAndStepDiscrepancy(&tests)

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
			// 1productname is not a CruxID
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
			// "abc" is not a valid valType
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
			// The "vals" "hash-set" below, which is the set of valid values for the
			// enum "paymenttype", shold not be empty
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
			// 1cash is not a CruxID
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
		// Both tasks and properties should not be empty
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
			// free*mug is not a CruxID
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
			tasks: []string{"freepen", "freemug", "freebag"},
			// Discount is not a CruxID
			properties: []string{"Discount", "pointsmult"},
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

func testMissingStart(tests *[]verifySchemaTest) {
	rs := RuleSchema{
		class: uccCreationClass,
		patternSchema: []AttrSchema{
			// vals below should also contain '"START": true'
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
			// there should be a "step" attribute-schema here
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
			tasks: []string{"getcustdetails", "aof", "sendauthlinktoclient"},
			// "abcd" should not be in properties
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
			tasks: []string{"getcustdetails", "aof", "sendauthlinktoclient"},
			// properties should contain "nextstep" (and should not contain "abcd")
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
				// "vals" should have exactly the same strings as "tasks" below, except "start" which is only in "vals"
				vals: map[string]bool{start: true, "getcustdetails": true, "aof": true, "sendauthlinktoclient": true},
			},
			{name: stepFailed, valType: typeBool},
			{name: "mode", valType: typeEnum, vals: map[string]bool{"physical": true, "demat": true}},
		},
		actionSchema: ActionSchema{
			// "tasks" should have exactly the same strings as "vals" above, except for "start"
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

func TestVerifyRuleSet(t *testing.T) {

	/* Business rules tests */
	setupPurchaseRuleSchema()
	setupRuleSetForPurchases()
	// the only two tests that involve no error, because the ruleset is correct
	testCorrectRS(t)
	testTaskAsAttrName(t)
	// in the rest of these tests, verifyRuleSet() should return an error
	testInvalidAttrName(t)
	testWrongAttrValType(t)
	testInvalidOp(t)
	testTaskNotInSchema(t)
	testPropNameNotInSchema(t)
	testBothReturnAndExit(t)

	/* Workflow tests */
	setupUCCCreationSchema()
	setupUCCCreationRuleSet()
	// the only test that involves no error, because the ruleset is correct
	testCorrectWF(t)
	// in the rest of these tests, verifyRuleSet() should return an error
	testWFRuleMissingStep(t)
	testWFRuleMissingBothNSAndDone(t)
	testWFNoTasksAndNotDone(t)
	testWFNextStepValNotInTasks(t)
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
	ruleSets[mainRS].Rules[1].RulePattern = []RulePatternTerm{
		{"product", opEQ, "jacket"},
		// priceabc is not in the schema
		{"priceabc", opGT, 50.0},
	}
	ok, err := verifyRuleSet(ruleSets[mainRS], false)
	if ok || err == nil {
		t.Errorf(incorrectOutputRSMsg + "invalid attr name")
	}
	ruleSets[mainRS].Rules[1].RulePattern = correctRP
}

func testTaskAsAttrName(t *testing.T) {
	ruleSets[mainRS].Rules[1].RulePattern = []RulePatternTerm{
		{"product", opEQ, "jacket"},
		// freejar is not in the pattern-schema, but it is a task in the action-schema
		{"freejar", opEQ, true},
	}
	ok, err := verifyRuleSet(ruleSets[mainRS], false)
	if !ok || err != nil {
		t.Errorf(incorrectOutputRSMsg + "a task 'tag' as an attribute name")
	}
	ruleSets[mainRS].Rules[1].RulePattern = correctRP
}

func testWrongAttrValType(t *testing.T) {
	ruleSets[mainRS].Rules[1].RulePattern = []RulePatternTerm{
		{"product", opEQ, "jacket"},
		// price should be a float, not a string
		{"price", opGT, "abc"},
	}
	ok, err := verifyRuleSet(ruleSets[mainRS], false)
	if ok || err == nil {
		t.Errorf(incorrectOutputRSMsg + "wrong attribute value type")
	}
	ruleSets[mainRS].Rules[1].RulePattern = correctRP
}

func testInvalidOp(t *testing.T) {
	ruleSets[mainRS].Rules[1].RulePattern = []RulePatternTerm{
		{"product", opEQ, "jacket"},
		// it should be "gt" (opGT), not "greater than"
		{"price", "greater than", 50.0},
	}
	ok, err := verifyRuleSet(ruleSets[mainRS], false)
	if ok || err == nil {
		t.Errorf(incorrectOutputRSMsg + "invalid operation")
	}
	ruleSets[mainRS].Rules[1].RulePattern = correctRP
}

// In each of the rule-action tests below, a rule-action is modified temporarily.
// After each test, we must reset the rule-action to the correct one below before
// moving on to the next test.
var correctRA RuleActions = RuleActions{
	Tasks:      []string{"freemug", "freejar", "freeplant"},
	Properties: []Property{{"discount", "20"}},
}

func testTaskNotInSchema(t *testing.T) {
	ruleSets[mainRS].Rules[3].RuleActions = RuleActions{
		// freeeraser is not in the schema
		Tasks:      []string{"freemug", "freeeraser"},
		Properties: []Property{{"discount", "20"}},
	}
	ok, err := verifyRuleSet(ruleSets[mainRS], false)
	if ok || err == nil {
		t.Errorf(incorrectOutputRSMsg + "task not in schema")
	}
	ruleSets[mainRS].Rules[3].RuleActions = correctRA
}

func testPropNameNotInSchema(t *testing.T) {
	ruleSets[mainRS].Rules[3].RuleActions = RuleActions{
		Tasks: []string{"freemug", "freejar", "freeplant"},
		// cashback is not a property in the action-schema
		Properties: []Property{{"cashback", "5"}},
	}
	ok, err := verifyRuleSet(ruleSets[mainRS], false)
	if ok || err == nil {
		t.Errorf(incorrectOutputRSMsg + "property name not in schema")
	}
	ruleSets[mainRS].Rules[3].RuleActions = correctRA
}

func testBothReturnAndExit(t *testing.T) {
	ruleSets[mainRS].Rules[3].RuleActions = RuleActions{
		Tasks:      []string{"freemug", "freejar", "freeplant"},
		Properties: []Property{{"discount", "20"}},
		// both WillReturn and WillExit below should not be true
		WillReturn: true,
		WillExit:   true,
	}
	ok, err := verifyRuleSet(ruleSets[mainRS], false)
	if ok || err == nil {
		t.Errorf(incorrectOutputRSMsg + "both RETURN and EXIT instructions")
	}
	ruleSets[mainRS].Rules[3].RuleActions = correctRA
}

func testCorrectWF(t *testing.T) {
	ok, err := verifyRuleSet(ruleSets[uccCreation], true)
	if !ok || err != nil {
		t.Errorf(incorrectOutputWFMsg + "no issues")
	}
}

func testWFRuleMissingStep(t *testing.T) {
	ruleSets[uccCreation].Rules[1].RulePattern = []RulePatternTerm{
		// there should be a "step" attribute here
		{stepFailed, opEQ, false},
		{"mode", opEQ, "physical"},
	}
	ok, err := verifyRuleSet(ruleSets[uccCreation], true)
	if ok || err == nil {
		t.Errorf(incorrectOutputWFMsg + "a rule missing 'step'")
	}
	// Reset to original correct rule-pattern
	ruleSets[uccCreation].Rules[1].RulePattern = []RulePatternTerm{
		{step, opEQ, "getcustdetails"},
		{stepFailed, opEQ, false},
		{"mode", opEQ, "physical"},
	}
}

// In each of the (workflow) rule-action tests below, a rule-action is modified temporarily.
// After each test, we must reset the rule-action to the correct one below before
// moving on to the next test.
var correctWorkflowRA RuleActions = RuleActions{
	Tasks:      []string{"aof", "kycvalid", "nomauth", "bankaccvalid"},
	Properties: []Property{{nextStep, "aof"}},
}

func testWFRuleMissingBothNSAndDone(t *testing.T) {
	ruleSets[uccCreation].Rules[1].RuleActions = RuleActions{
		Tasks: []string{"aof", "kycvalid", "nomauth", "bankaccvalid"},
		// Properties below should contain at least one of "nextstep" and "done"
		Properties: []Property{},
	}
	ok, err := verifyRuleSet(ruleSets[uccCreation], true)
	if ok || err == nil {
		t.Errorf(incorrectOutputWFMsg + "a rule missing both 'nextstep' and 'done'")
	}
	ruleSets[uccCreation].Rules[1].RuleActions = correctWorkflowRA
}

func testWFNoTasksAndNotDone(t *testing.T) {
	ruleSets[uccCreation].Rules[1].RuleActions = RuleActions{
		// Either Tasks below should not be empty, or Properties below should contain {"done", "true"}
		Tasks:      []string{},
		Properties: []Property{{nextStep, "abc"}},
	}
	ok, err := verifyRuleSet(ruleSets[uccCreation], true)
	if ok || err == nil {
		t.Errorf(incorrectOutputWFMsg + "a rule with no tasks and no 'done=true'")
	}
	ruleSets[uccCreation].Rules[1].RuleActions = correctWorkflowRA
}

func testWFNextStepValNotInTasks(t *testing.T) {
	ruleSets[uccCreation].Rules[1].RuleActions = RuleActions{
		Tasks: []string{"aof", "kycvalid", "nomauth", "bankaccvalid"},
		// "abcd" below is not in "Tasks" above
		Properties: []Property{{nextStep, "abcd"}},
	}
	ok, err := verifyRuleSet(ruleSets[uccCreation], true)
	if ok || err == nil {
		t.Errorf(incorrectOutputWFMsg + "a 'nextstep' value not in its rule's 'tasks'")
	}
	ruleSets[uccCreation].Rules[1].RuleActions = correctWorkflowRA
}
