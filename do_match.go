package main

import (
	"errors"
	"fmt"
)

const (
	timeLayout = "2006-01-02T15:04:05Z"
)

type Entity struct {
	class string
	attrs []Attr
}

type Attr struct {
	name string
	val  string
}

type ActionSet struct {
	tasks      []string
	properties []Property
}

type Property struct {
	name string
	val  string
}

var ruleSets = make(map[string]RuleSet)

type RuleSet struct {
	ver     int
	class   string
	setName string
	rules   []Rule
}

type Rule struct {
	rulePattern []RulePatternTerm
	ruleActions RuleActions
}

type RulePatternTerm struct {
	attrName string
	op       string
	attrVal  any
}

type RuleActions struct {
	tasks      []string
	properties []Property
	thenCall   string
	elseCall   string
	willReturn bool
	willExit   bool
}

func doMatch(entity Entity, ruleSet RuleSet, actionSet ActionSet, seenRuleSets map[string]bool) (ActionSet, bool, error) {
	if seenRuleSets[ruleSet.setName] {
		return ActionSet{}, false, errors.New("ruleset has already been traversed")
	}
	seenRuleSets[ruleSet.setName] = true
	for _, rule := range ruleSet.rules {
		willExit := false
		matched, err := matchPattern(entity, rule.rulePattern, actionSet)
		if err != nil {
			return ActionSet{}, false, err
		}
		if matched {
			actionSet = collectActions(actionSet, rule.ruleActions)
			if len(rule.ruleActions.thenCall) > 0 {
				setToCall := ruleSets[rule.ruleActions.thenCall]
				if setToCall.class != entity.class {
					return inconsistentRuleSet(setToCall.setName, ruleSet.setName)
				}
				var err error
				actionSet, willExit, err = doMatch(entity, setToCall, actionSet, seenRuleSets)
				if err != nil {
					return ActionSet{}, false, err
				}
			}
			if willExit || rule.ruleActions.willExit {
				return actionSet, true, nil
			}
			if rule.ruleActions.willReturn {
				delete(seenRuleSets, ruleSet.setName)
				return actionSet, false, nil
			}
		} else if len(rule.ruleActions.elseCall) > 0 {
			setToCall := ruleSets[rule.ruleActions.elseCall]
			if setToCall.class != entity.class {
				return inconsistentRuleSet(setToCall.setName, ruleSet.setName)
			}
			var err error
			actionSet, willExit, err = doMatch(entity, setToCall, actionSet, seenRuleSets)
			if err != nil {
				return ActionSet{}, false, err
			} else if willExit {
				return actionSet, true, nil
			}
		}
	}
	delete(seenRuleSets, ruleSet.setName)
	return actionSet, false, nil
}

func inconsistentRuleSet(calledSetName string, currSetName string) (ActionSet, bool, error) {
	return ActionSet{}, false, fmt.Errorf("system inconsistency with BRE rule terms, attempting to call %v from %v",
		calledSetName, currSetName,
	)
}
