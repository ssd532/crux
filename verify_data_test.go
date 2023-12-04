package main

import "testing"

type verifySchemaTest struct {
	name    string
	rs      RuleSchema
	isWF    bool
	want    bool
	wantErr bool
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
