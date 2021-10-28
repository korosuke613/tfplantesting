package tfplantesting

import (
	"encoding/json"
	"fmt"
	tfjson "github.com/hashicorp/terraform-json"
	"io/ioutil"
	"regexp"
)

func GetCompositeActionType(actions *tfjson.Actions) ChangeAction {
	switch {
	case actions.NoOp():
		return NoOp
	case actions.Create():
		return Create
	case actions.Read():
		return Read
	case actions.Update():
		return Update
	case actions.Delete():
		return Delete
	case actions.CreateBeforeDestroy():
		return CreateBeforeDestroy
	case actions.DestroyBeforeCreate():
		return DestroyBeforeCreate
	case actions.Replace():
		return Replace
	default:
		panic("detected a non-existent Action")
	}
}

func ParsePlanJson(planFilePath string) (*tfjson.Plan, error) {
	planJson, err := ioutil.ReadFile(planFilePath)
	if err != nil {
		return nil, fmt.Errorf("fail open %v: %v", planFilePath, err)
	}

	plan := tfjson.Plan{}

	if err := json.Unmarshal(planJson, &plan); err != nil {
		return nil, fmt.Errorf("fail parse %v: %v", planFilePath, err)
	}

	return &plan, nil
}

func GetResourceChanges(plan *tfjson.Plan) ChangeResources {
	changeResources := ChangeResources{}

	for _, change := range plan.ResourceChanges {
		if change.Change.Actions.Read() {
			continue
		}
		changeResources[change.Address] = *change
	}

	return changeResources
}

func GetResourceChangesWithData(plan *tfjson.Plan) ChangeResources {
	changeResources := ChangeResources{}

	for _, change := range plan.ResourceChanges {
		changeResources[change.Address] = *change
	}

	return changeResources
}

func GetActions(changeResources ChangeResources) map[tfjson.Action][]tfjson.ResourceChange {
	actions := map[tfjson.Action][]tfjson.ResourceChange{}
	for _, resource := range changeResources {
		for _, action := range resource.Change.Actions {
			actions[action] = append(actions[action], resource)
		}
	}

	return actions
}

func GetCompositeActions(changeResources ChangeResources) map[ChangeAction][]tfjson.ResourceChange {
	actions := map[ChangeAction][]tfjson.ResourceChange{}
	for _, resource := range changeResources {
		action := GetCompositeActionType(&resource.Change.Actions)
		actions[action] = append(actions[action], resource)
	}

	return actions
}

func SearchAction(changeResources ChangeResources, actionType tfjson.Action) []tfjson.ResourceChange {
	// 特定のActionがあるかどうか
	actions := GetActions(changeResources)

	if _, ok := actions[actionType]; ok {
		return actions[actionType]
	}
	return nil
}

func SearchCompositeAction(changeResources ChangeResources, actionType ChangeAction) []tfjson.ResourceChange {
	// 特定の複合Actionがあるかどうか
	actions := GetCompositeActions(changeResources)

	if _, ok := actions[actionType]; ok {
		return actions[actionType]
	}
	return nil
}

func SearchResourceType(changeResources ChangeResources, resourceType string) (map[string][]tfjson.ResourceChange, map[string][]tfjson.ResourceChange) {
	foundTypes := map[string][]tfjson.ResourceChange{}
	withoutTypes := map[string][]tfjson.ResourceChange{}

	for _, resource := range changeResources {
		if resource.Change.Actions.NoOp() {
			continue
		}
		if regexp.MustCompile(resourceType).Match([]byte(resource.Type)) {
			foundTypes[resource.Type] = append(foundTypes[resourceType], resource)
		} else {
			withoutTypes[resource.Type] = append(withoutTypes[resourceType], resource)
		}
	}
	return foundTypes, withoutTypes
}

func SearchResourceModule(changeResources ChangeResources, resourceModule string) (map[string][]tfjson.ResourceChange, map[string][]tfjson.ResourceChange) {
	foundModule := map[string][]tfjson.ResourceChange{}
	withoutModule := map[string][]tfjson.ResourceChange{}

	for _, resource := range changeResources {
		if resource.Change.Actions.NoOp() {
			continue
		}
		if regexp.MustCompile(resourceModule).Match([]byte(resource.ModuleAddress)) {
			foundModule[resource.Type] = append(foundModule[resourceModule], resource)
		} else {
			withoutModule[resource.Type] = append(withoutModule[resourceModule], resource)
		}
	}
	return foundModule, withoutModule
}
