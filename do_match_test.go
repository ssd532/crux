package main

import (
	"reflect"
	"testing"
)

type doMatchTest struct {
	name      string
	entity    Entity
	ruleSet   RuleSet
	actionSet ActionSet
	want      ActionSet
}

func TestDoMatch(t *testing.T) {
	tests := []doMatchTest{}

	// BRE tests
	testBasic(&tests)
	testExit(&tests)
	testReturn(&tests)
	testsWithTransactions(&tests)
	testsWithOrders(&tests)

	// WFE tests
	testUCCCreation(&tests)
	testComplexWF(&tests)

	t.Logf("Running %v doMatch() tests\n", len(tests))
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, _, _ := doMatch(tt.entity, tt.ruleSet, tt.actionSet, map[string]bool{})
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("\n\ndoMatch() = %v, \n\nwant %v\n\n", got, tt.want)
			}
		})
	}

	// Test for cyclical rulesets
	testCycleError(t)
}
