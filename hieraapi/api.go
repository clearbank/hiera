package hieraapi

import (
	"github.com/lyraproj/issue/issue"
	"github.com/lyraproj/pcore/px"
)

type Kind string

// HieraRoot is an option key that can be used to change the default root which is the current working directory
const HieraRoot = `Hiera::Root`

// HieraConfigFileName is an option that can be used to change the default file name 'hiera.yaml'
const HieraConfigFileName = `Hiera::ConfigFileName`

// HieraConfig is an option that can be used to change absolute path of the hiera config. When specified, the
// HieraRoot and HieraConfigFileName will not have any effect.
const HieraConfig = `Hiera::Config`

// HieraScope is an option that can be used to pass a variable scope to Hiera. This scope is used
// by the 'scope' lookup_key provider function and when doing variable interpolations
const HieraScope = `Hiera::Scope`

const KindDataDig = Kind(`data_dig`)
const KindDataHash = Kind(`data_hash`)
const KindLookupKey = Kind(`lookup_key`)

var FunctionKeys = []string{string(KindDataDig), string(KindDataHash), string(KindLookupKey)}

var LocationKeys = []string{string(LcPath), `paths`, string(LcGlob), `globs`, string(LcUri), `uris`, string(LcMappedPaths)}

var ReservedOptionKeys = []string{string(LcPath), string(LcUri)}

type Function interface {
	Kind() Kind
	Name() string
	Resolve(ic Invocation) (Function, bool)
}

type Entry interface {
	Copy(Config) Entry
	Options() px.OrderedMap
	OptionsMap() map[string]px.Value
	DataDir() string
	Function() Function
	Name() string
	Resolve(ic Invocation, defaults Entry) Entry
	CreateProvider() DataProvider
	Locations() []Location
}

type Config interface {
	// Root returns the directory holding this Config
	Root() string

	// Path is the full path to this Config
	Path() string

	// Defaults returns the Defaults entry
	Defaults() Entry

	// Hierarchy returns the configuration hierarchy slice
	Hierarchy() []Entry

	// DefaultHierarchy returns the default hierarchy slice
	DefaultHierarchy() []Entry

	// Resolve resolves this instance into a ResolveHierarchy. Resolving means creating the proper
	// DataProviders for all Hierarchy entries
	Resolve(ic Invocation) ResolvedConfig
}

type ResolvedConfig interface {
	// Config returns the original Config that the receiver was created from
	Config() Config

	// Hierarchy returns the DataProvider slice
	Hierarchy() []DataProvider

	// DefaultHierarchy returns the DataProvider slice for the configured default_hierarchy.
	// The slice will be empty if no such hierarchy has been defined.
	DefaultHierarchy() []DataProvider

	// LookupOptions returns the resolved lookup_options value for the given key or nil
	// if no such options exists
	LookupOptions(key Key) map[string]px.Value
}

// An Invocation keeps track of one specific lookup invocation implements a guard against
// endless recursion
type Invocation interface {
	px.Context

	Config() ResolvedConfig

	DoWithScope(scope px.Keyed, doer px.Doer)

	// Call doer and while it is executing, don't reveal any found values in logs
	DoRedacted(doer px.Doer)

	// ReportText will add the message returned by the given function to the
	// lookup explainer. The method will only get called when the explanation
	// support is enabled
	ReportText(messageProducer func() string)
	ReportLocationNotFound()
	ReportFound(key interface{}, value px.Value)
	ReportMergeResult(value px.Value)
	ReportMergeSource(source string)
	ReportNotFound(key interface{})

	WithDataProvider(pvd DataProvider, f px.Producer) px.Value

	WithInterpolation(expr string, f px.Producer) px.Value

	WithInvalidKey(key interface{}, f px.Producer) px.Value

	WithLocation(loc Location, f px.Producer) px.Value

	WithLookup(key Key, f px.Producer) px.Value

	WithMerge(ms MergeStrategy, f px.Producer) px.Value

	WithSegment(seg interface{}, f px.Producer) px.Value

	WithSubLookup(key Key, f px.Producer) px.Value

	// ExplainMode returns true if explain support is active
	ExplainMode() bool

	// ForConfig returns an Invocation that without explainer support
	ForConfig() Invocation

	// ForData returns an Invocation that has adjusted its explainer according to
	// how it should report lookup of data as opposed to the "lookup_options" key.
	ForData() Invocation

	// ForLookupOptions returns an Invocation that has adjusted its explainer according to
	// how it should report lookup of the "lookup_options" key.
	ForLookupOptions() Invocation
}

// NotFound is the error that Hiera will panic with when a value cannot be found and no default
// value has been defined
var NotFound issue.Reported

type DataDig func(ic ProviderContext, key Key, options map[string]px.Value) px.Value

type DataHash func(ic ProviderContext, options map[string]px.Value) px.OrderedMap

type LookupKey func(ic ProviderContext, key string, options map[string]px.Value) px.Value
