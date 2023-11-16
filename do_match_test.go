package main

import (
	"fmt"
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

	/**************
	   BRE tests
	**************/
	testBasic(&tests)
	testExit(&tests)
	testReturn(&tests)
	testTransactions(&tests)
	testPurchases(&tests)
	testOrders(&tests)

	/**************
	   WFE tests
	**************/
	testUCCCreation(&tests)
	testPrepareAOF(&tests)
	testValidateAOF(&tests)
	// An artificial workflow with random meaningless data, for testing purposes
	testComplexWF(&tests)

	fmt.Printf("Running %v doMatch() tests\n", len(tests))
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, _, _ := doMatch(tt.entity, tt.ruleSet, tt.actionSet, map[string]bool{})
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("\n\ndoMatch() = %v, \n\nwant        %v\n\n", got, tt.want)
			}
		})
	}

	// Test for cyclical rulesets
	testCycleError(t)
}
