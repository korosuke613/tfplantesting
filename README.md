# tfplantesting
Go library for testing `terraform plan` results.

See, [examples/example_test.go](examples/example_test.go).

## Getting Started

### Generate Plan.json

```shell
terraform plan -out ./plan.out
terraform show -json ./plan.out > plan.json
```

### Write the test code
Testing: Only module bar should be changed.

```go
package your_package

import (
	tfplantesting "github.com/korosuke613/tfplantesting/src"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestPlan(t *testing.T) {
	t.Parallel()

	t.Run("Only module bar should be changed", func(t *testing.T) {
		// Parse plan.json.
		plan, err := tfplantesting.ParsePlanJson("./plan.json")
		if err != nil {
			assert.Fail(t, err.Error())
		}

		// Get resources that will be changed.
		changeResources := tfplantesting.GetResourceChanges(plan)
		
		// Get the resource to be changed for the module whose name contains `foo`.
		changesInFoo, changesWithoutFoo := tfplantesting.SearchResourceModule(changeResources, "bar")

		// Verify that the resources of the module whose name contains `foo` are changed.
		assert.NotEqual(t, 0, len(changesInFoo))

		// Verify that no other resources have changed.
		assert.Equal(t, 0, len(changesWithoutFoo))
	})
}
```

## License
The basic license is the MIT license, but the code in the following directory is under the MPL2 license.

- [examples/testdata](examples/testdata)
