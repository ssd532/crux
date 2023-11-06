package main

import (
	"errors"
	"fmt"
	"strconv"
	"time"
)

func matchPattern(entity Entity, rulePattern []RulePatternTerm, actionSet ActionSet) (bool, error) {
	for _, term := range rulePattern {
		valType := ""
		entityAttrVal := ""
		for _, entityAttr := range entity.attrs {
			if entityAttr.name == term.attrName {
				entityAttrVal = entityAttr.val
				valType = getTypeFromSchema(entity.class, entityAttr.name)
			}
		}
		if entityAttrVal == "" {
			// Check whether the attribute name in the pattern term matches any of the tasks in
			// the action-set
			for _, task := range actionSet.tasks {
				if task == term.attrName {
					entityAttrVal = "true"
					valType = "bool"
				}
			}
		}
		if entityAttrVal == "" {
			entityAttrVal = "false"
			valType = "bool"
		}
		matched, err := makeComparison(entityAttrVal, term.attrVal, valType, term.op)
		if err != nil {
			return false, fmt.Errorf("error making comparison %w", err)
		}
		if !matched {
			return false, nil
		}
	}
	return true, nil
}

func getTypeFromSchema(class string, attrName string) string {
	for _, ruleSchema := range ruleSchemas {
		if ruleSchema.class == class {
			for _, attrSchema := range ruleSchema.patternSchema {
				if attrSchema.name == attrName {
					return attrSchema.valType
				}
			}
		}
	}
	return ""
}

func makeComparison(entityAttrVal string, termAttrVal any, valType string, op string) (bool, error) {
	entityAttrValConv, err := convertEntityAttrVal(entityAttrVal, valType)
	if err != nil {
		return false, fmt.Errorf("error converting value: %w", err)
	}
	switch op {
	case "eq":
		return entityAttrValConv == termAttrVal, nil
	case "ne":
		return entityAttrValConv != termAttrVal, nil
	}
	orderedTypes := map[string]bool{"int": true, "float": true, "ts": true, "str": true}
	if !orderedTypes[valType] {
		return false, errors.New("not an ordered type")
	}
	var result int8
	var match bool
	switch op {
	case "lt":
		result, err = compare(entityAttrValConv, termAttrVal)
		match = (result == -1)
	case "le":
		result, err = compare(entityAttrValConv, termAttrVal)
		match = (result == -1) || (result == 0)
	case "gt":
		result, err = compare(entityAttrValConv, termAttrVal)
		match = (result == 1)
	case "ge":
		result, err = compare(entityAttrValConv, termAttrVal)
		match = (result == 1) || (result == 0)
	}
	if err != nil {
		return false, fmt.Errorf("error making comparison %w", err)
	}
	return match, nil
}

func convertEntityAttrVal(entityAttrVal string, valType string) (any, error) {
	var entityAttrValConv any
	var err error
	switch valType {
	case "bool":
		entityAttrValConv, err = strconv.ParseBool(entityAttrVal)
	case "int":
		entityAttrValConv, err = strconv.Atoi(entityAttrVal)
	case "float":
		entityAttrValConv, err = strconv.ParseFloat(entityAttrVal, 64)
	case "str", "enum":
		entityAttrValConv = entityAttrVal
	case "ts":
		entityAttrValConv, err = time.Parse(timeLayout, entityAttrVal)
	}
	if err != nil {
		return nil, err
	}
	return entityAttrValConv, nil
}

// The compare function returns:
// 0 if a == b,
// -1 if a < b, or
// 1 if a > b
func compare(a any, b any) (int8, error) {
	if a == b {
		return 0, nil
	}
	var lessThan bool
	switch a.(type) {
	case int:
		if a.(int) < b.(int) {
			lessThan = true
		}
	case float64:
		if a.(float64) < b.(float64) {
			lessThan = true
		}
	case string:
		if a.(string) < b.(string) {
			lessThan = true
		}
	case time.Time:
		if a.(time.Time).Before(b.(time.Time)) {
			lessThan = true
		}
	default:
		return -2, errors.New("invalid type")
	}
	if lessThan {
		return -1, nil
	} else {
		return 1, nil
	}
}
