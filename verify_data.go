package main

import (
	"fmt"
	"reflect"
	"regexp"
)

const (
	cruxIDRegExp = `^[a-z][a-z0-9_]*$`

	step       = "step"
	stepFailed = "stepfailed"
	start      = "START"
	nextStep   = "nextstep"
	done       = "done"
)

var validTypes = map[string]bool{
	typeBool: true, typeInt: true, typeFloat: true, typeStr: true, typeEnum: true, typeTS: true,
}

// Parameters
// rs RuleSchema: the RuleSchema to be verified
// isWF bool: true if the RuleSchema applies to a workflow, otherwise false
func verifyRuleSchema(rs RuleSchema, isWF bool) (bool, error) {
	if len(rs.class) == 0 {
		return false, fmt.Errorf("schema class is empty string")
	}
	if _, err := verifyPatternSchema(rs, isWF); err != nil {
		return false, err
	}
	if _, err := verifyActionSchema(rs, isWF); err != nil {
		return false, err
	}
	return true, nil
}

func verifyPatternSchema(rs RuleSchema, isWF bool) (bool, error) {
	if len(rs.patternSchema) == 0 {
		return false, fmt.Errorf("pattern-schema for %v is empty", rs.class)
	}
	re := regexp.MustCompile(cruxIDRegExp)
	// Bools needed for workflows only
	stepFound, stepFailedFound := false, false

	for _, attrSchema := range rs.patternSchema {
		if !re.MatchString(attrSchema.name) {
			return false, fmt.Errorf("attribute name %v is not a valid CruxID", attrSchema.name)
		} else if !validTypes[attrSchema.valType] {
			return false, fmt.Errorf("%v is not a valid value-type", attrSchema.valType)
		} else if attrSchema.valType == typeEnum && len(attrSchema.vals) == 0 {
			return false, fmt.Errorf("no valid values for enum %v", attrSchema.name)
		}
		for val := range attrSchema.vals {
			if !re.MatchString(val) && val != start {
				return false, fmt.Errorf("enum value %v is not a valid CruxID", val)
			}
		}

		// Workflows only
		if attrSchema.name == step && attrSchema.valType == typeEnum {
			stepFound = true
		}
		if isWF && attrSchema.name == step && !attrSchema.vals[start] {
			return false, fmt.Errorf("workflow schema for %v doesn't allow step=START", rs.class)
		}
		if attrSchema.name == stepFailed && attrSchema.valType == typeBool {
			stepFailedFound = true
		}
	}

	// Workflows only
	if isWF && (!stepFound || !stepFailedFound) {
		return false, fmt.Errorf("necessary workflow attributes absent in schema for class %v", rs.class)
	}

	return true, nil
}

func verifyActionSchema(rs RuleSchema, isWF bool) (bool, error) {
	re := regexp.MustCompile(cruxIDRegExp)
	if len(rs.actionSchema.tasks) == 0 && len(rs.actionSchema.properties) == 0 {
		return false, fmt.Errorf("both tasks and properties are empty in schema for class %v", rs.class)
	}
	for _, task := range rs.actionSchema.tasks {
		if !re.MatchString(task) {
			return false, fmt.Errorf("task %v is not a valid CruxID", task)
		}
	}

	// Workflows only
	if isWF && len(rs.actionSchema.properties) != 2 {
		return false, fmt.Errorf("action-schema for %v does not contain exactly two properties", rs.class)
	}
	nextStepFound, doneFound := false, false

	for _, propName := range rs.actionSchema.properties {
		if !re.MatchString(propName) {
			return false, fmt.Errorf("property name %v is not a valid CruxID", propName)
		} else if propName == nextStep {
			nextStepFound = true
		} else if propName == done {
			doneFound = true
		}
	}

	// Workflows only
	if isWF && (!nextStepFound || !doneFound) {
		return false, fmt.Errorf("action-schema for %v does not contain both the properties 'nextstep' and 'done'", rs.class)
	}
	if isWF && !reflect.DeepEqual(getTasksMapForWF(rs.actionSchema.tasks), getStepAttrVals(rs)) {
		return false, fmt.Errorf("action-schema tasks for %v are not the same as valid values for 'step' in pattern-schema", rs.class)
	}
	return true, nil
}

func getTasksMapForWF(tasks []string) map[string]bool {
	tm := map[string]bool{}
	for _, t := range tasks {
		tm[t] = true
	}
	// To allow comparison with the set of valid values for the 'step' attribute, which includes "START"
	tm[start] = true
	return tm
}

func getStepAttrVals(rs RuleSchema) map[string]bool {
	for _, ps := range rs.patternSchema {
		if ps.name == step {
			return ps.vals
		}
	}
	return nil
}
