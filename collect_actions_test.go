package main

import (
	"reflect"
	"testing"
)

func TestCollectActionsBasic(t *testing.T) {
	actionSet := ActionSet{
		tasks:      []string{"dodiscount", "yearendsale"},
		properties: []Property{{"discount", "7"}, {"shipby", "fedex"}},
	}

	ruleActions := RuleActions{
		Tasks:      []string{"yearendsale", "summersale"},
		Properties: []Property{{"cashback", "10"}, {"discount", "9"}},
		ThenCall:   "domesticpo",
		WillReturn: false,
		WillExit:   true,
	}

	want := ActionSet{
		tasks:      []string{"dodiscount", "yearendsale", "summersale"},
		properties: []Property{{"discount", "9"}, {"shipby", "fedex"}, {"cashback", "10"}},
	}

	res := collectActions(actionSet, ruleActions)
	if !reflect.DeepEqual(want, res) {
		t.Errorf("\n\ncollectActions() = %v, \n\nwant %v\n\n", res, want)
	}
}

func TestCollectActionsWithEmptyRuleActions(t *testing.T) {
	actionSet := ActionSet{
		tasks:      []string{"dodiscount", "yearendsale"},
		properties: []Property{{"discount", "7"}, {"shipby", "fedex"}},
	}

	ruleActions := RuleActions{}

	want := ActionSet{
		tasks:      []string{"dodiscount", "yearendsale"},
		properties: []Property{{"discount", "7"}, {"shipby", "fedex"}},
	}

	res := collectActions(actionSet, ruleActions)
	if !reflect.DeepEqual(want, res) {
		t.Errorf("\n\ncollectActions() = %v, \n\nwant %v\n\n", res, want)
	}
}

func TestCollectActionsWithEmptyActionSet(t *testing.T) {
	actionSet := ActionSet{}

	ruleActions := RuleActions{
		Tasks:      []string{"dodiscount", "yearendsale"},
		Properties: []Property{{"discount", "7"}, {"shipby", "fedex"}},
		ThenCall:   "overseaspo",
		WillReturn: true,
		WillExit:   false,
	}

	want := ActionSet{
		tasks:      []string{"dodiscount", "yearendsale"},
		properties: []Property{{"discount", "7"}, {"shipby", "fedex"}},
	}

	res := collectActions(actionSet, ruleActions)
	if !reflect.DeepEqual(want, res) {
		t.Errorf("\n\ncollectActions() = %v, \n\nwant %v\n\n", res, want)
	}
}
