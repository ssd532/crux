package main

import (
	"testing"
	"time"
)

func TestMatchPattern(t *testing.T) {
	var testNames []string
	var entities []Entity
	var rulePatterns []([]RulePatternTerm)
	var resultsExpected []any

	actionSet := ActionSet{tasks: []string{"dodiscount", "yearendsale"}}

	// Test: many terms, everything matches
	testNames = append(testNames, "everything matches")
	entities = append(entities, sampleEntity)
	receivedTime, _ := time.Parse(timeLayout, "2018-05-15T12:00:00Z")
	rulePatterns = append(rulePatterns, []RulePatternTerm{
		{"cat", "eq", "textbook"},
		{"fullname", "eq", "Advanced Physics"},
		{"ageinstock", "le", 7},
		{"mrp", "lt", 51.2},
		{"received", "gt", receivedTime},
		{"bulkorder", "ne", false},
		{"dodiscount", "eq", true},
	})
	resultsExpected = append(resultsExpected, true)

	// Test: many terms, mrp doesn't match
	testNames = append(testNames, "mrp doesn't match")
	entities = append(entities, sampleEntity)
	receivedTime, _ = time.Parse(timeLayout, "2018-05-15T12:00:00Z")
	rulePatterns = append(rulePatterns, []RulePatternTerm{
		{"cat", "eq", "textbook"},
		{"fullname", "eq", "Advanced Physics"},
		{"ageinstock", "le", 7},
		{"mrp", "ge", 51.2},
		{"received", "gt", receivedTime},
		{"bulkorder", "ne", false},
		{"dodiscount", "eq", true},
	})
	resultsExpected = append(resultsExpected, false)

	// Test: bool "ne"
	testNames = append(testNames, "bool ne")
	entities = append(entities, sampleEntity)
	rulePatterns = append(rulePatterns, []RulePatternTerm{
		{"bulkorder", "ne", true},
	})
	resultsExpected = append(resultsExpected, false)

	// Test: enum "ne"
	testNames = append(testNames, "enum ne")
	entities = append(entities, sampleEntity)
	rulePatterns = append(rulePatterns, []RulePatternTerm{
		{"cat", "ne", "refbook"},
	})
	resultsExpected = append(resultsExpected, true)

	// Test: float "eq"
	testNames = append(testNames, "float eq")
	entities = append(entities, sampleEntity)
	rulePatterns = append(rulePatterns, []RulePatternTerm{
		{"mrp", "eq", 50.8},
	})
	resultsExpected = append(resultsExpected, true)

	// Test: float "ge"
	testNames = append(testNames, "float ge")
	entities = append(entities, sampleEntity)
	rulePatterns = append(rulePatterns, []RulePatternTerm{
		{"mrp", "ge", 50.8},
	})
	resultsExpected = append(resultsExpected, true)

	// Test: timestamp "lt"
	testNames = append(testNames, "timestamp lt")
	entities = append(entities, sampleEntity)
	receivedTime, _ = time.Parse(timeLayout, "2018-06-10T15:04:05Z")
	rulePatterns = append(rulePatterns, []RulePatternTerm{
		{"received", "lt", receivedTime},
	})
	resultsExpected = append(resultsExpected, true)

	// Test: timestamp "le"
	testNames = append(testNames, "timestamp le")
	entities = append(entities, sampleEntity)
	receivedTime, _ = time.Parse(timeLayout, "2018-05-01T15:04:05Z")
	rulePatterns = append(rulePatterns, []RulePatternTerm{
		{"received", "le", receivedTime},
	})
	resultsExpected = append(resultsExpected, false)

	// Test: string "lt"
	testNames = append(testNames, "string lt")
	entities = append(entities, sampleEntity)
	rulePatterns = append(rulePatterns, []RulePatternTerm{
		{"fullname", "lt", "Advanced Science"},
	})
	resultsExpected = append(resultsExpected, true)

	// Test: string "gt"
	testNames = append(testNames, "string gt")
	entities = append(entities, sampleEntity)
	rulePatterns = append(rulePatterns, []RulePatternTerm{
		{"fullname", "gt", "Accelerated Physics"},
	})
	resultsExpected = append(resultsExpected, true)

	// Test: task wanted in pattern, found in action set
	testNames = append(testNames, "tasks found in action set")
	entities = append(entities, sampleEntity)
	rulePatterns = append(rulePatterns, []RulePatternTerm{
		{"dodiscount", "eq", true},
		{"yearendsale", "ne", false},
	})
	resultsExpected = append(resultsExpected, true)

	// Test: task wanted in pattern, but not found in action set
	testNames = append(testNames, "task not in action set")
	entities = append(entities, sampleEntity)
	rulePatterns = append(rulePatterns, []RulePatternTerm{
		{"dodiscount", "eq", true},
		{"summersale", "eq", true},
	})
	resultsExpected = append(resultsExpected, false)

	// Test: task not wanted in pattern, and not found in action set
	testNames = append(testNames, "task 'eq false' in pattern, and not in action set")
	entities = append(entities, sampleEntity)
	rulePatterns = append(rulePatterns, []RulePatternTerm{
		{"summersale", "eq", false},
	})
	resultsExpected = append(resultsExpected, true)

	// Test: edge case - no rule pattern
	testNames = append(testNames, "no rule pattern")
	entities = append(entities, sampleEntity)
	rulePatterns = append(rulePatterns, []RulePatternTerm{})
	resultsExpected = append(resultsExpected, true)

	// Test: error converting value
	testNames = append(testNames, "deliberate error converting value")
	entities = append(entities, Entity{"inventoryitems", []Attr{
		{"ageinstock", "abc"},
	}})
	rulePatterns = append(rulePatterns, []RulePatternTerm{
		{"ageinstock", "gt", 5},
	})
	resultsExpected = append(resultsExpected, nil)

	// Test: error - not an ordered type
	testNames = append(testNames, "deliberate error: not an ordered type")
	entities = append(entities, sampleEntity)
	rulePatterns = append(rulePatterns, []RulePatternTerm{
		{"bulkorder", "gt", true},
	})
	resultsExpected = append(resultsExpected, nil)

	// Run the tests
	t.Log("==Running", len(rulePatterns), "matchPattern tests==")

	for i, rulePattern := range rulePatterns {
		t.Logf("Test: %s", testNames[i])
		res, err := matchPattern(entities[i], rulePattern, actionSet)
		if resultsExpected[i] == nil && err == nil {
			t.Errorf("Expected but did not get error")
			continue
		} else if resultsExpected[i] == nil && err != nil {
			continue
		} else if res != resultsExpected[i] {
			t.Errorf("FAIL output=%v, expected=%v", res, resultsExpected[i])
		}
	}
}