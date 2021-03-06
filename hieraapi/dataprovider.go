package hieraapi

import "github.com/lyraproj/pcore/px"

type DataProvider interface {
	// Lookup performs an lookup using this data provider
	Lookup(key Key, invocation Invocation, merge MergeStrategy) px.Value

	// FullName returns a descriptive name of the data provider. Used by the explainer
	FullName() string
}
