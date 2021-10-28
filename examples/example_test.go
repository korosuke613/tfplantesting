package tfplantesting

import (
	tfjson "github.com/hashicorp/terraform-json"
	tfplantesting "github.com/korosuke613/tfplantesting/src"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestExamples(t *testing.T) {
	t.Parallel()

	t.Run("The change should be 7", func(t *testing.T) {
		plan, err := tfplantesting.ParsePlanJson("./testdata/basic/plan.json")
		if err != nil {
			assert.Fail(t, err.Error())
		}

		resourceChanges := tfplantesting.GetResourceChanges(plan)
		numberOfChangeResources := len(resourceChanges)

		assert.Equal(t, 7, numberOfChangeResources)
	})

	t.Run("Only module bar should be changed", func(t *testing.T) {
		plan, err := tfplantesting.ParsePlanJson("./testdata/deep_module/plan.json")
		if err != nil {
			assert.Fail(t, err.Error())
		}

		resourceChanges := tfplantesting.GetResourceChanges(plan)
		changesInFoo, changesWithoutFoo := tfplantesting.SearchResourceModule(resourceChanges, "bar")

		assert.Equal(t, 0, len(changesWithoutFoo))
		assert.NotEqual(t, 0, len(changesInFoo))
	})

	t.Run("Only resource-type null_resource should be changed", func(t *testing.T) {
		plan, err := tfplantesting.ParsePlanJson("./testdata/config_resource_depends_on/plan.json")
		if err != nil {
			assert.Fail(t, err.Error())
		}

		resourceChanges := tfplantesting.GetResourceChanges(plan)
		changesInNullResource, changesWithoutNullResource := tfplantesting.SearchResourceType(resourceChanges, "null_resource")

		assert.Equal(t, 0, len(changesWithoutNullResource))
		assert.NotEqual(t, 0, len(changesInNullResource))
	})

	t.Run("Resource bar or foo in null_resource should be removed", func(t *testing.T) {
		plan, err := tfplantesting.ParsePlanJson("./testdata/config_resource_depends_on/plan.json")
		if err != nil {
			assert.Fail(t, err.Error())
		}

		resourceChanges := tfplantesting.GetResourceChanges(plan)
		deleteResources := tfplantesting.SearchAction(resourceChanges, tfjson.ActionDelete)

		assert.NotNil(t, deleteResources)

		for _, r := range deleteResources {
			assert.Equal(t, r.Type, "null_resource")
			assert.Contains(t, []string{"bar", "foo"}, r.Name)
		}
	})
}
