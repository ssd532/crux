package main

func collectActions(actionSet ActionSet, ruleActions RuleActions) ActionSet {
	newActionSet := ActionSet{}

	// Union-set of tasks
	newActionSet.tasks = append(newActionSet.tasks, actionSet.tasks...)
	for _, newTask := range ruleActions.tasks {
		found := false
		for _, task := range newActionSet.tasks {
			if newTask == task {
				found = true
				break
			}
		}
		if !found {
			newActionSet.tasks = append(newActionSet.tasks, newTask)
		}
	}

	// Union-set of workflow names
	newActionSet.workflows = append(newActionSet.workflows, actionSet.workflows...)
	for _, newWF := range ruleActions.workflows {
		found := false
		for _, wf := range newActionSet.workflows {
			if newWF == wf {
				found = true
				break
			}
		}
		if !found {
			newActionSet.workflows = append(newActionSet.workflows, newWF)
		}
	}

	// Perform "union-set" of properties, overwriting previous property values if needed
	newActionSet.properties = append(newActionSet.properties, actionSet.properties...)
	for _, newProperty := range ruleActions.properties {
		found := false
		for i, property := range newActionSet.properties {
			if property.name == newProperty.name {
				newActionSet.properties[i].val = newProperty.val
				found = true
				break
			}
		}
		if !found {
			newActionSet.properties = append(newActionSet.properties, newProperty)
		}
	}
	return newActionSet
}
