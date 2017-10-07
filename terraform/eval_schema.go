package terraform

import (
	"fmt"
	"strings"

	"github.com/hashicorp/terraform/config/configschema"
)

// EvalGetProviderSchema gets the schema for a provider's main configuration,
// as would appear inside a "provider" block.
type EvalGetProviderSchema struct {
	ProviderName string
	Provider     *ResourceProvider
	Output       **configschema.Block
}

func (n *EvalGetProviderSchema) Eval(ctx EvalContext) (interface{}, error) {
	if !providerSupportsSchema(*n.Provider) {
		return nil, fmt.Errorf(strings.TrimSpace(evalGetProviderSchemaNotCompatible), n.ProviderName)
	}

	schema, err := (*n.Provider).ProviderSchema()
	if err != nil {
		return nil, err
	}

	// a nil schema is used to indicate no schema is available, so if we
	// actually _got_ a nil schema then we'll promote it to an empty one
	// to remove this ambiguity.
	if schema == nil {
		schema = &configschema.Block{}
	}

	*n.Output = schema
	return nil, nil
}

// The following error uses user-oriented terminology, assuming that
// the only reason we'd require schema is if we're in an HCL2-format
// config file.
const evalGetProviderSchemaNotCompatible = `
provider %q is not compatible with the HCL2 experiment. A newer version may be compatible; if not, configuration for this provider must be placed in a configuration file that does not opt in to the experiment.
`

func providerSupportsSchema(provider ResourceProvider) bool {
	// Since the "ProviderSchema" function was added to ResourceProvider
	// without a change to the provider plugin protocol version, we must
	// sniff for support of this new feature, which we do by verifying
	// that at least one resource or data source has the SchemaAvailable
	// flag set. This weird sniffing protocol is designed to work within
	// the pre-existing set of methods so a breaking change could be avoided.
	if resources := provider.Resources(); len(resources) > 0 {
		return resources[0].SchemaAvailable
	} else if dataSources := provider.DataSources(); len(dataSources) > 0 {
		return dataSources[0].SchemaAvailable
	}

	// (a provider with no resources or data sources can't support schema
	// per this sniffing approach, but an empty provider would be useless anyway.)

	return false
}
