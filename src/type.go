package tfplantesting

import tfjson "github.com/hashicorp/terraform-json"

type ChangeAction string

const (
	NoOp                             = ChangeAction(tfjson.ActionNoop)
	Create                           = ChangeAction(tfjson.ActionCreate)
	Read                             = ChangeAction(tfjson.ActionRead)
	Update                           = ChangeAction(tfjson.ActionUpdate)
	Delete                           = ChangeAction(tfjson.ActionDelete)
	CreateBeforeDestroy ChangeAction = "create-before-destroy"
	DestroyBeforeCreate ChangeAction = "destroy-before-create"
	Replace             ChangeAction = "replace"
)

type ChangeResources map[string]tfjson.ResourceChange
