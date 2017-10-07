package terraform

import (
	"testing"

	"github.com/hashicorp/terraform/config/configschema"
)

func TestEvalGetProviderSchema(t *testing.T) {
	var provider ResourceProvider
	schema := &configschema.Block{}
	provider = &MockResourceProvider{
		ProviderSchemaReturn: schema,

		// Need to have at least one resource with SchemaAvailable so we
		// can sniff to see that this provider supports schema methods.
		ResourcesReturn: []ResourceType{
			{
				Name:            "baz_bar",
				SchemaAvailable: true,
			},
		},
	}
	var got *configschema.Block

	n := &EvalGetProviderSchema{
		ProviderName: "baz",
		Provider:     &provider,
		Output:       &got,
	}

	_, err := n.Eval(nil)
	if err != nil {
		t.Errorf("unexpected error: %s", err)
	}

	if got != schema {
		t.Errorf("result is not the pointer we provided in the mock")
	}
}
