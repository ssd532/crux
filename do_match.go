/*
This file contains doMatch() and a helper function called by doMatch().
It also contains ruleSets, a map in which we are currently storing all
rulesets for the purpose of testing doMatch().
*/

package main

import (
	"errors"
	"fmt"
)

var ruleSets = make(map[string]RuleSet)

func doMatch(entity Entity, ruleSet RuleSet, actionSet ActionSet, seenRuleSets map[string]bool) (ActionSet, bool, error) {
	if seenRuleSets[ruleSet.SetName] {
		return ActionSet{}, false, errors.New("ruleset has already been traversed")
	}
	seenRuleSets[ruleSet.SetName] = true
	for _, rule := range ruleSet.Rules {
		willExit := false
		matched, err := matchPattern(entity, rule.RulePattern, actionSet)
		if err != nil {
			return ActionSet{}, false, err
		}
		if matched {
			actionSet = collectActions(actionSet, rule.RuleActions)
			if len(rule.RuleActions.ThenCall) > 0 {
				setToCall := ruleSets[rule.RuleActions.ThenCall]
				if setToCall.Class != entity.class {
					return inconsistentRuleSet(setToCall.SetName, ruleSet.SetName)
				}
				var err error
				actionSet, willExit, err = doMatch(entity, setToCall, actionSet, seenRuleSets)
				if err != nil {
					return ActionSet{}, false, err
				}
			}
			if willExit || rule.RuleActions.WillExit {
				return actionSet, true, nil
			}
			if rule.RuleActions.WillReturn {
				delete(seenRuleSets, ruleSet.SetName)
				return actionSet, false, nil
			}
		} else if len(rule.RuleActions.ElseCall) > 0 {
			setToCall := ruleSets[rule.RuleActions.ElseCall]
			if setToCall.Class != entity.class {
				return inconsistentRuleSet(setToCall.SetName, ruleSet.SetName)
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
	delete(seenRuleSets, ruleSet.SetName)
	return actionSet, false, nil
}

func inconsistentRuleSet(calledSetName string, currSetName string) (ActionSet, bool, error) {
	return ActionSet{}, false, fmt.Errorf("system inconsistency with BRE rule terms, attempting to call %v from %v",
		calledSetName, currSetName,
	)
}
